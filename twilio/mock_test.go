package twilio

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"testing"

	"github.com/twilio/terraform-provider-twilio/core"
)

// mockTwilioAPI is an in-memory, stateful emulator of the Twilio REST API used
// by the acceptance tests so they can run in CI without contacting Twilio (no
// credentials, no billing). It stores the JSON object produced when a resource
// is created and returns it verbatim on subsequent reads, which keeps Terraform
// refresh and ImportStateVerify consistent.
type mockTwilioAPI struct {
	mu    sync.Mutex
	items map[string]map[string]any // canonical item path -> stored object
	seq   int
}

func newMockTwilioAPI() *mockTwilioAPI {
	return &mockTwilioAPI{items: map[string]map[string]any{}}
}

// sidPrefixes maps "<hostKeyword>/<Collection>" to the two-letter SID prefix
// Twilio would use. Collections that are not listed fall back to a generic
// prefix; that is fine for resources whose SID schema fields are plain strings.
var sidPrefixes = map[string]string{
	"serverless/Services":     "ZS",
	"serverless/Functions":    "ZH",
	"serverless/Environments": "ZE",
	"chat/Services":           "IS",
	"studio/Flows":            "FW",
	"flex/FlexFlows":          "FO",
}

var serviceSidRe = regexp.MustCompile(`/Services/([A-Za-z0-9]+)/`)

func (m *mockTwilioAPI) RoundTrip(req *http.Request) (*http.Response, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	path := req.URL.Path

	switch req.Method {
	case http.MethodDelete:
		delete(m.items, path)
		return jsonResponse(req, http.StatusNoContent, nil), nil

	case http.MethodGet:
		if obj, ok := m.items[path]; ok {
			return jsonResponse(req, http.StatusOK, obj), nil
		}
		return jsonResponse(req, http.StatusNotFound, map[string]any{
			"code": 20404, "message": "The requested resource was not found", "status": 404,
		}), nil

	case http.MethodPost:
		form := parseForm(req)
		if obj, ok := m.items[path]; ok {
			// Update an existing item: merge the posted params.
			mergeForm(obj, form)
			return jsonResponse(req, http.StatusOK, obj), nil
		}
		// Create a new item under the collection at `path`.
		m.seq++
		sid := m.newSid(req.URL.Host, path)
		obj := map[string]any{"sid": sid}
		mergeForm(obj, form)
		if match := serviceSidRe.FindStringSubmatch(path + "/"); match != nil {
			obj["service_sid"] = match[1]
		}
		m.items[path+"/"+sid] = obj
		return jsonResponse(req, http.StatusCreated, obj), nil
	}

	return jsonResponse(req, http.StatusMethodNotAllowed, nil), nil
}

func (m *mockTwilioAPI) newSid(host, path string) string {
	collection := lastSegment(path)
	prefix := "XX"
	for key, p := range sidPrefixes {
		parts := strings.SplitN(key, "/", 2)
		if strings.Contains(host, parts[0]) && parts[1] == collection {
			prefix = p
			break
		}
	}
	return fmt.Sprintf("%s%032d", prefix, m.seq)
}

func lastSegment(path string) string {
	trimmed := strings.TrimSuffix(strings.TrimSuffix(path, "/"), ".json")
	idx := strings.LastIndex(trimmed, "/")
	return strings.TrimSuffix(trimmed[idx+1:], ".json")
}

// mergeForm copies posted form params into obj, converting Twilio's PascalCase
// parameter names to the snake_case JSON field names the models expect and
// keeping only the first value for each key.
// stringValuedKeys are snake_case fields that the Twilio models type as string
// even though their values could look boolean/numeric (e.g.
// FlexFlow.contact_identity, or a unique_name that happens to be all digits).
// They must not be coerced away from string.
var stringValuedKeys = map[string]bool{
	"contact_identity": true,
	"unique_name":      true,
	"friendly_name":    true,
	"channel_type":     true,
	"commit_message":   true,
	"status":           true,
	"domain_suffix":    true,
}

func mergeForm(obj map[string]any, form url.Values) {
	for k, v := range form {
		if len(v) == 0 {
			continue
		}
		setNested(obj, k, v[0])
	}
}

// setNested stores a posted form value into obj. Twilio uses dotted parameter
// names for nested objects (e.g. "Integration.FlowSid"), which must be rebuilt
// as nested JSON so the typed models unmarshal correctly.
func setNested(obj map[string]any, formKey, raw string) {
	segments := strings.Split(formKey, ".")
	cur := obj
	for i, seg := range segments {
		key := core.ToSnakeCase(seg)
		if i == len(segments)-1 {
			if stringValuedKeys[key] {
				cur[key] = raw
			} else {
				cur[key] = coerce(raw)
			}
			return
		}
		next, ok := cur[key].(map[string]any)
		if !ok {
			next = map[string]any{}
			cur[key] = next
		}
		cur = next
	}
}

// coerce converts a form-encoded string into the JSON type the Twilio models
// expect (bool or number where applicable), so the response unmarshals into the
// typed twilio-go structs.
func coerce(s string) any {
	switch s {
	case "true":
		return true
	case "false":
		return false
	}
	if i, err := strconv.ParseInt(s, 10, 64); err == nil {
		return i
	}
	return s
}

func parseForm(req *http.Request) url.Values {
	if req.Body == nil {
		return url.Values{}
	}
	body, _ := io.ReadAll(req.Body)
	values, _ := url.ParseQuery(string(body))
	return values
}

func jsonResponse(req *http.Request, status int, body map[string]any) *http.Response {
	resp := &http.Response{
		StatusCode: status,
		Status:     fmt.Sprintf("%d %s", status, http.StatusText(status)),
		Header:     http.Header{"Content-Type": []string{"application/json"}},
		Request:    req,
	}
	if body == nil {
		resp.Body = io.NopCloser(bytes.NewReader(nil))
		return resp
	}
	encoded, _ := json.Marshal(body)
	resp.Body = io.NopCloser(bytes.NewReader(encoded))
	resp.ContentLength = int64(len(encoded))
	return resp
}

// setupMockProvider installs the in-memory Twilio API for the duration of the
// test and provides dummy credentials so the provider configures successfully.
func setupMockProvider(t *testing.T) {
	t.Helper()
	// twilio-go validates that credentials are alphanumeric before sending, so
	// the dummy values must contain no punctuation.
	t.Setenv(AccountSid, "AC00000000000000000000000000000000")
	t.Setenv(ApiKey, "SK00000000000000000000000000000000")
	t.Setenv(ApiSecret, "dummysecret00000000000000000000000")

	testHTTPTransport = newMockTwilioAPI()
	t.Cleanup(func() { testHTTPTransport = nil })
}
