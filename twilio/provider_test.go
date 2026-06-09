package twilio

import (
	"os"
	"testing"

	"github.com/twilio/terraform-provider-twilio/twilio/resources"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const ProviderName = "twilio"

var testAccProviderFactories map[string]func() (*schema.Provider, error)

func init() {
	testAccProviderFactories = map[string]func() (*schema.Provider, error){
		ProviderName: func() (*schema.Provider, error) { return Provider(), nil },
	}
}

func TestMain(m *testing.M) {
	resource.TestMain(m)
}

func TestTwilioResourcesMap(t *testing.T) {
	twilioResources := resources.NewTwilioResources()
	if twilioResources.Map["twilio_api_accounts_messages"] == nil {
		t.Fatal("expected twilio_api_accounts_messages to be registered")
	}
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv(AccountSid); v == "" {
		t.Fatal("TWILIO_ACCOUNT_SID must be set for acceptance tests")
	}
	if v := os.Getenv(ApiKey); v == "" {
		t.Fatal("TWILIO_API_KEY must be set for acceptance tests")
	}
	if v := os.Getenv(ApiSecret); v == "" {
		t.Fatal("TWILIO_API_SECRET must be set for acceptance tests")
	}
}
