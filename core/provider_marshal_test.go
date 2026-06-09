package core

import (
	"reflect"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestTagNaming(t *testing.T) {
	type innerStruct struct {
		T1 string `json:"t1"`
		T2 string
		T3 string `json:"t3a" provider:"t3b"`
		T4 string `provider:"t4"`
		T5 string `json:"t5" provider:",id"`
		T6 string `provider:",id"`
		T7 string `provider:",ignore"`
		T8 string `provider:",flatten"`
	}
	var field reflect.StructField
	var tag tagInfo
	field, _ = reflect.TypeOf(innerStruct{}).FieldByName("T1")
	tag = getNameFromTag(field)
	if tag.name != "t1" {
		t.Errorf("wrong name: got %q, want %q", tag.name, "t1")
	}
	if tag.isId {
		t.Error("Id set")
	}
	if tag.flatten {
		t.Error("Flatten set")
	}
	if tag.ignore {
		t.Error("ignore set")
	}
	field, _ = reflect.TypeOf(innerStruct{}).FieldByName("T2")
	tag = getNameFromTag(field)
	if tag.name != "T2" {
		t.Errorf("wrong name: got %q, want %q", tag.name, "T2")
	}
	if tag.isId {
		t.Error("Id set")
	}
	if tag.flatten {
		t.Error("Flatten set")
	}
	if tag.ignore {
		t.Error("ignore set")
	}
	field, _ = reflect.TypeOf(innerStruct{}).FieldByName("T3")
	tag = getNameFromTag(field)
	if tag.name != "t3b" {
		t.Errorf("wrong name: got %q, want %q", tag.name, "t3b")
	}
	if tag.isId {
		t.Error("Id set")
	}
	if tag.flatten {
		t.Error("Flatten set")
	}
	if tag.ignore {
		t.Error("ignore set")
	}
	field, _ = reflect.TypeOf(innerStruct{}).FieldByName("T4")
	tag = getNameFromTag(field)
	if tag.name != "t4" {
		t.Errorf("wrong name: got %q, want %q", tag.name, "t4")
	}
	if tag.isId {
		t.Error("Id set")
	}
	if tag.flatten {
		t.Error("Flatten set")
	}
	if tag.ignore {
		t.Error("ignore set")
	}
	field, _ = reflect.TypeOf(innerStruct{}).FieldByName("T5")
	tag = getNameFromTag(field)
	if tag.name != "t5" {
		t.Errorf("wrong name: got %q, want %q", tag.name, "t5")
	}
	if !tag.isId {
		t.Error("Id not set")
	}
	if tag.flatten {
		t.Error("Flatten set")
	}
	if tag.ignore {
		t.Error("ignore set")
	}
	field, _ = reflect.TypeOf(innerStruct{}).FieldByName("T6")
	tag = getNameFromTag(field)
	if tag.name != "T6" {
		t.Errorf("wrong name: got %q, want %q", tag.name, "T6")
	}
	if !tag.isId {
		t.Error("Id not set")
	}
	if tag.flatten {
		t.Error("Flatten set")
	}
	if tag.ignore {
		t.Error("ignore set")
	}
	field, _ = reflect.TypeOf(innerStruct{}).FieldByName("T7")
	tag = getNameFromTag(field)
	if tag.name != "T7" {
		t.Errorf("wrong name: got %q, want %q", tag.name, "T7")
	}
	if tag.isId {
		t.Error("Id set")
	}
	if tag.flatten {
		t.Error("Flatten set")
	}
	if !tag.ignore {
		t.Error("ignore not set")
	}
	field, _ = reflect.TypeOf(innerStruct{}).FieldByName("T8")
	tag = getNameFromTag(field)
	if tag.name != "T8" {
		t.Errorf("wrong name: got %q, want %q", tag.name, "T8")
	}
	if tag.isId {
		t.Error("Id set")
	}
	if !tag.flatten {
		t.Error("Flatten not set")
	}
	if tag.ignore {
		t.Error("ignore set")
	}
}

func TestSimpleUnmarshal(t *testing.T) {
	terraformSchema := map[string]*schema.Schema{
		"T1": {
			Type:     schema.TypeBool,
			Required: true,
		},
		"T2": {
			Type:     schema.TypeInt,
			Required: true,
		},
		"T3": {
			Type:     schema.TypeString,
			Required: true,
		},
		"T4": {
			Type:     schema.TypeInt,
			Optional: true,
		},
	}
	data := map[string]interface{}{
		"T1": true,
		"T2": 1,
		"T3": "t3",
	}
	resourceData := schema.TestResourceDataRaw(t, terraformSchema, data)

	type innerStruct struct {
		T1 bool
		T2 int
		T3 string
		T4 *int
	}

	testStruct := innerStruct{}
	if err := UnmarshalSchema(&testStruct, resourceData); err != nil {
		t.Errorf("Unmarshall failed: result '%v'", err)
	}
	if testStruct.T1 != true {
		t.Error("T1 did not unmarshal")
	}
	if testStruct.T2 != 1 {
		t.Error("T2 did not unmarshal")
	}
	if testStruct.T3 != "t3" {
		t.Error("T3 did not unmarshal")
	}
	if testStruct.T4 != nil {
		t.Error("T4 should be nil")
	}
}

func TestComplexUnmarshal(t *testing.T) {
	terraformSchema := map[string]*schema.Schema{
		"T1": {
			Type:     schema.TypeString,
			Required: true,
		},
		"T2": {
			Type: schema.TypeList,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
			Required: true,
		},
		"T3": {
			Type:     schema.TypeString,
			Required: true,
		},
		"T4": {
			Type:     schema.TypeInt,
			Required: true,
		},
		"T5": {
			Type:     schema.TypeFloat,
			Required: true,
		},
		"T6": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"T7": {
			Type:     schema.TypeInt,
			Optional: true,
		},
	}
	data := map[string]interface{}{
		"T1": "AC00112233445566778899aabbccddeeff",
		"T2": []interface{}{"t2a", "t2b"},
		"T3": "t3",
		"T4": 1,
		"T5": 1.0,
		"T6": "AC00112233445566778899aabbccddeeff",
	}
	resourceData := schema.TestResourceDataRaw(t, terraformSchema, data)
	resourceData.SetId("t0")

	type innerStruct struct {
		T0 *string `provider:",id"`
		T1 AccountSid
		T2 []string
		T3 *string
		T4 *int
		T5 *float64
		T6 *AccountSid
		T7 *int
	}

	testStruct := innerStruct{}
	if err := UnmarshalSchema(&testStruct, resourceData); err != nil {
		t.Errorf("Unmarshall failed: result '%v'", err)
	}
	if *testStruct.T0 != "t0" {
		t.Error("T0 did not unmarshal")
	}
	if testStruct.T1.String() != "AC00112233445566778899aabbccddeeff" {
		t.Error("T1 did not unmarshal")
	}
	if len(testStruct.T2) != 2 {
		t.Error("T2 did not unmarshal")
	}
	if testStruct.T2[0] != "t2a" {
		t.Error("T2 did not unmarshal")
	}
	if testStruct.T2[1] != "t2b" {
		t.Error("T2 did not unmarshal")
	}
	if *testStruct.T3 != "t3" {
		t.Error("T3 did not unmarshal")
	}
	if *testStruct.T4 != 1 {
		t.Error("T4 did not unmarshal")
	}
	if *testStruct.T5 != 1.0 {
		t.Error("T5 did not unmarshal")
	}
	if testStruct.T6.String() != "AC00112233445566778899aabbccddeeff" {
		t.Error("T6 did not unmarshal")
	}
	if testStruct.T7 != nil {
		t.Error("T7 did not unmarshal")
	}
}

func TestBoolUnmarshal(t *testing.T) {
	terraformSchema := map[string]*schema.Schema{
		"custom_code_enabled": {
			Type:     schema.TypeBool,
			Optional: true,
			Computed: true,
		},
		"lookup_enabled": {
			Type:     schema.TypeBool,
			Optional: true,
		},
	}
	data := map[string]interface{}{
		"custom_code_enabled": false,
		"lookup_enabled":      false,
	}

	resourceData := schema.TestResourceDataRaw(t, terraformSchema, data)
	resourceData.SetId("t0")

	type innerStruct struct {
		CustomCodeEnabled *bool `json:"CustomCodeEnabled,omitempty"`
		LookupEnabled     *bool `json:"LookupEnabled,omitempty"`
	}

	testStruct := innerStruct{}
	if err := UnmarshalSchema(&testStruct, resourceData); err != nil {
		t.Errorf("Unmarshall failed: result '%v'", err)
	}

	if *testStruct.LookupEnabled != false {
		t.Error("LookupEnabled did not unmarshal")
	}
	if *testStruct.CustomCodeEnabled != false {
		t.Error("CustomCodeEnabled did not unmarshal")
	}
}

func TestTimeUnMarshal(t *testing.T) {
	terraformSchema := map[string]*schema.Schema{
		"T0": {
			Type:     schema.TypeString,
			Required: true,
		},
	}
	data := map[string]interface{}{
		"T0": "2021-05-17T01:35:33Z",
	}
	resourceData := schema.TestResourceDataRaw(t, terraformSchema, data)
	err := resourceData.Set("T0", "2021-05-17T01:35:33Z")

	type innerStruct struct {
		T0 *time.Time
	}

	testStruct := innerStruct{}
	if err := UnmarshalSchema(&testStruct, resourceData); err != nil {
		t.Errorf("Unmarshall failed: result '%v'", err)
	}

	if err != nil {
		t.Errorf("Date T0 did not unmarshal: %v", err)
	}
	if resourceData.Get("T0") != "2021-05-17T01:35:33Z" {
		t.Error("Date T0 did not unmarshal")
	}
}

func TestUnmarshalNilValueToPointer(t *testing.T) {
	terraformSchema := map[string]*schema.Schema{
		"T1": {
			Type:     schema.TypeString,
			Required: true,
		},
		"T2": {
			Type:     schema.TypeList,
			Required: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
	}
	data := map[string]interface{}{
		"T1": nil,
		"T2": nil,
	}
	resourceData := schema.TestResourceDataRaw(t, terraformSchema, data)

	type innerStruct struct {
		T1 *string
		T2 *[]string
	}

	testStruct := innerStruct{}
	if err := UnmarshalSchema(&testStruct, resourceData); err != nil {
		t.Errorf("Unmarshall failed: result '%v'", err)
	}
	if testStruct.T1 != nil {
		t.Error("T1 did not unmarshal")
	}
	if testStruct.T2 != nil {
		t.Error("T2 did not unmarshal")
	}
}

func TestFlattenUnmarshal(t *testing.T) {
	terraformSchema := map[string]*schema.Schema{
		"T1": {
			Type:     schema.TypeString,
			Required: true,
		},
		"T2": {
			Type:     schema.TypeString,
			Required: true,
		},
	}
	data := map[string]interface{}{
		"T1": "t1",
		"T2": "t2",
	}
	resourceData := schema.TestResourceDataRaw(t, terraformSchema, data)
	resourceData.SetId("t0")

	type test2 struct {
		T1 string
		T2 string
	}

	type innerStruct struct {
		T0     *string `provider:",id"`
		Nested test2   `provider:",flatten"`
	}

	testStruct := innerStruct{}
	if err := UnmarshalSchema(&testStruct, resourceData); err != nil {
		t.Errorf("Unmarshall failed: result '%v'", err)
	}
	if *testStruct.T0 != "t0" {
		t.Error("T0 did not unmarshal")
	}
	if testStruct.Nested.T1 != "t1" {
		t.Error("T1 did not unmarshal")
	}
	if testStruct.Nested.T2 != "t2" {
		t.Error("T2 did not unmarshal")
	}
}

func TestListUnmarshal(t *testing.T) {
	terraformSchema := map[string]*schema.Schema{
		"T1": {
			Type: schema.TypeList,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"T1a": {
						Type:     schema.TypeString,
						Required: true,
					},
					"T1b": {
						Type:     schema.TypeString,
						Required: true,
					},
				},
			},
			Required: true,
		},
	}
	data := map[string]interface{}{
		"T1": []interface{}{
			map[string]interface{}{"T1a": "r1", "T1b": "r2"},
			map[string]interface{}{"T1a": "r3", "T1b": "r4"},
		},
	}
	resourceData := schema.TestResourceDataRaw(t, terraformSchema, data)
	resourceData.SetId("t0")

	type test2 struct {
		T1a string
		T1b string
	}

	type innerStruct struct {
		T1 []test2
	}

	testStruct := innerStruct{}
	if err := UnmarshalSchema(&testStruct, resourceData); err != nil {
		t.Errorf("Unmarshall failed: result '%v'", err)
	}
	if len(testStruct.T1) != 2 {
		t.Error("T1 did not unmarshal")
	}
	if testStruct.T1[0].T1a != "r1" {
		t.Error("T1 did not unmarshal")
	}
	if testStruct.T1[0].T1b != "r2" {
		t.Error("T1 did not unmarshal")
	}
	if testStruct.T1[1].T1a != "r3" {
		t.Error("T1 did not unmarshal")
	}
	if testStruct.T1[1].T1b != "r4" {
		t.Error("T1 did not unmarshal")
	}
}

func TestSidListUnmarshal(t *testing.T) {
	terraformSchema := map[string]*schema.Schema{
		"T1": {
			Type:     schema.TypeList,
			Elem:     &schema.Schema{Type: schema.TypeString},
			Required: true,
		},
	}
	data := map[string]interface{}{
		"T1": []interface{}{"AC00112233445566778899aabbccddeefe", "AC00112233445566778899aabbccddeeff"},
	}
	resourceData := schema.TestResourceDataRaw(t, terraformSchema, data)
	resourceData.SetId("t0")

	type innerStruct struct {
		T1 []AccountSid
	}

	testStruct := innerStruct{}
	if err := UnmarshalSchema(&testStruct, resourceData); err != nil {
		t.Errorf("Unmarshall failed: result '%v'", err)
	}
	if len(testStruct.T1) != 2 {
		t.Error("T1 did not unmarshal")
	}
	if testStruct.T1[0].String() != "AC00112233445566778899aabbccddeefe" {
		t.Error("T1 did not unmarshal")
	}
	if testStruct.T1[1].String() != "AC00112233445566778899aabbccddeeff" {
		t.Error("T1 did not unmarshal")
	}
}

func TestMapUnmarshal(t *testing.T) {
	terraformSchema := map[string]*schema.Schema{
		"T1": {
			Type:     schema.TypeMap,
			Elem:     &schema.Schema{Type: schema.TypeString},
			Required: true,
		},
	}
	data := map[string]interface{}{
		"T1": map[string]interface{}{"T1a": "r1", "T1b": "r2"},
	}
	resourceData := schema.TestResourceDataRaw(t, terraformSchema, data)
	resourceData.SetId("t0")

	type innerStruct struct {
		T1 map[string]string
	}

	testStruct := innerStruct{}
	if err := UnmarshalSchema(&testStruct, resourceData); err != nil {
		t.Errorf("Unmarshall failed: result '%v'", err)
	}
	if len(testStruct.T1) != 2 {
		t.Error("T1 did not unmarshal")
	}
	if testStruct.T1["T1a"] != "r1" {
		t.Error("T1 did not unmarshal")
	}
	if testStruct.T1["T1b"] != "r2" {
		t.Error("T1 did not unmarshal")
	}
}

func TestSidMapUnmarshal(t *testing.T) {
	terraformSchema := map[string]*schema.Schema{
		"T1": {
			Type:     schema.TypeMap,
			Elem:     &schema.Schema{Type: schema.TypeString},
			Required: true,
		},
	}
	data := map[string]interface{}{
		"T1": map[string]interface{}{"T1a": "AC00112233445566778899aabbccddeefe", "T1b": "AC00112233445566778899aabbccddeeff"},
	}
	resourceData := schema.TestResourceDataRaw(t, terraformSchema, data)
	resourceData.SetId("t0")

	type innerStruct struct {
		T1 map[string]AccountSid
	}

	testStruct := innerStruct{}
	if err := UnmarshalSchema(&testStruct, resourceData); err != nil {
		t.Errorf("Unmarshall failed: result '%v'", err)
	}
	if len(testStruct.T1) != 2 {
		t.Error("T1 did not unmarshal")
	}
	if testStruct.T1["T1a"].String() != "AC00112233445566778899aabbccddeefe" {
		t.Error("T1 did not unmarshal")
	}
	if testStruct.T1["T1b"].String() != "AC00112233445566778899aabbccddeeff" {
		t.Error("T1 did not unmarshal")
	}
}

func TestComplexMapUnmarshal(t *testing.T) {
	terraformSchema := map[string]*schema.Schema{
		"T1": {
			Type:     schema.TypeString,
			Required: true,
		},
	}
	data := map[string]interface{}{
		"T1": "{\"M0\":{\"T1a\":\"r1\",\"T1b\":true},\"M1\":{\"T1a\":\"r2\",\"T1b\":false}}",
	}
	resourceData := schema.TestResourceDataRaw(t, terraformSchema, data)
	resourceData.SetId("t0")

	type test2 struct {
		T1a string
		T1b bool
	}

	type innerStruct struct {
		T1 map[string]test2
	}

	testStruct := innerStruct{}
	if err := UnmarshalSchema(&testStruct, resourceData); err != nil {
		t.Errorf("Unmarshall failed: result '%v'", err)
	}
	if len(testStruct.T1) != 2 {
		t.Error("T1 did not unmarshal")
	}
	if testStruct.T1["M0"].T1a != "r1" {
		t.Error("T1 did not unmarshal")
	}
	if testStruct.T1["M0"].T1b != true {
		t.Error("T1 did not unmarshal")
	}
	if testStruct.T1["M1"].T1a != "r2" {
		t.Error("T1 did not unmarshal")
	}
	if testStruct.T1["M1"].T1b != false {
		t.Error("T1 did not unmarshal")
	}
}

func TestPureJsonUnmarshal(t *testing.T) {
	terraformSchema := map[string]*schema.Schema{
		"T1": {
			Type:     schema.TypeString,
			Required: true,
		},
	}
	data := map[string]interface{}{
		"T1": "{\"test\":\"test_value\"}",
	}
	resourceData := schema.TestResourceDataRaw(t, terraformSchema, data)

	type innerStruct struct {
		T1 interface{}
	}

	testStruct := innerStruct{}
	if err := UnmarshalSchema(&testStruct, resourceData); err != nil {
		t.Errorf("Unmarshall failed: result '%v'", err)
	}
	if testStruct.T1.(map[string]interface{})["test"] != "test_value" {
		t.Error("T1 did not unmarshal")
	}
}

func TestJsonEncodedNilUnMarshal(t *testing.T) {
	terraformSchema := map[string]*schema.Schema{
		"T1": {
			Type:     schema.TypeString,
			Required: true,
		},
	}
	data := map[string]interface{}{}
	resourceData := schema.TestResourceDataRaw(t, terraformSchema, data)

	type test struct {
		T1 *interface{}
	}

	testStruct := test{}
	if err := UnmarshalSchema(&testStruct, resourceData); err != nil {
		t.Errorf("Unmarshall failed: result '%v'", err)
	}
	if testStruct.T1 != nil {
		t.Error("T2 did not unmarshal")
	}
}

func TestOptionalFlattenUnmarshal(t *testing.T) {
	terraformSchema := map[string]*schema.Schema{
		"T1": {
			Type:     schema.TypeString,
			Required: true,
		},
		"T2": {
			Type:     schema.TypeString,
			Required: true,
		},
	}
	data := map[string]interface{}{}
	resourceData := schema.TestResourceDataRaw(t, terraformSchema, data)
	resourceData.SetId("t0")

	type test2 struct {
		T1 *string
		T2 *string
	}

	type innerStruct struct {
		T0     *string `provider:",id"`
		Nested test2   `provider:",flatten"`
	}

	testStruct := innerStruct{}
	if err := UnmarshalSchema(&testStruct, resourceData); err != nil {
		t.Errorf("Unmarshall failed: result '%v'", err)
	}
	if *testStruct.T0 != "t0" {
		t.Error("T0 did not unmarshal")
	}
	if testStruct.Nested.T1 != nil {
		t.Error("T1 did not unmarshal")
	}
	if testStruct.Nested.T2 != nil {
		t.Error("T2 did not unmarshal")
	}
}

func TestSimpleMarshal(t *testing.T) {
	terraformSchema := map[string]*schema.Schema{
		"T1": {
			Type:     schema.TypeBool,
			Required: true,
		},
		"T2": {
			Type:     schema.TypeInt,
			Required: true,
		},
		"T3": {
			Type:     schema.TypeString,
			Required: true,
		},
		"T4": {
			Type:     schema.TypeInt,
			Optional: true,
		},
		"T5": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"T6": {
			Type:     schema.TypeString,
			Optional: true,
		},
	}
	data := map[string]interface{}{}
	resourceData := schema.TestResourceDataRaw(t, terraformSchema, data)

	type innerStruct struct {
		T0 string `provider:",id"`
		T1 bool
		T2 int
		T3 string
		T4 *int
		T5 *interface{}
		T6 *interface{}
	}

	testStruct := innerStruct{T0: "t0", T1: true, T2: 1, T3: "t3", T5: nil}
	if err := MarshalSchema(resourceData, &testStruct); err != nil {
		t.Errorf("Marshall failed: result '%v'", err)
	}
	if resourceData.Id() != "t0" {
		t.Error("Id did not marshal")
	}
	if resourceData.Get("T1") != true {
		t.Error("T1 did not marshal")
	}
	if resourceData.Get("T2") != 1 {
		t.Error("T2 did not marshal")
	}
	if resourceData.Get("T3") != "t3" {
		t.Error("T3 did not marshal")
	}
	if resourceData.Get("T4") != 0 {
		t.Error("T4 did not marshal")
	}
	if resourceData.Get("T5") != "" {
		t.Error("T5 did not marshal")
	}
	if resourceData.Get("T6") != "" {
		t.Error("T6 did not marshal")
	}
}

func TestComplexMarshal(t *testing.T) {
	terraformSchema := map[string]*schema.Schema{
		"T1": {
			Type:     schema.TypeString,
			Required: true,
		},
		"T2": {
			Type: schema.TypeList,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
			Required: true,
		},
		"T3": {
			Type:     schema.TypeString,
			Required: true,
		},
		"T4": {
			Type:     schema.TypeInt,
			Required: true,
		},
	}
	data := map[string]interface{}{}
	resourceData := schema.TestResourceDataRaw(t, terraformSchema, data)

	type innerStruct struct {
		T0 *string `provider:",id"`
		T1 AccountSid
		T2 []string
		T3 *string
		T4 *int
	}

	testStr := "t0"
	testSid := AccountSid{}
	_ = testSid.Set("AC00112233445566778899aabbccddeeff")
	testDateString := "2010-04-01"
	testInt := 1

	testStruct := innerStruct{T0: &testStr, T1: testSid, T2: []string{"t2a", "t2b"}, T3: &testDateString, T4: &testInt}
	if err := MarshalSchema(resourceData, &testStruct); err != nil {
		t.Errorf("Marshall failed: result '%v'", err)
	}
	if resourceData.Id() != "t0" {
		t.Error("Id did not marshal")
	}
	if resourceData.Get("T1") != "AC00112233445566778899aabbccddeeff" {
		t.Error("T1 did not marshal")
	}
	if len(resourceData.Get("T2").([]interface{})) != 2 {
		t.Error("T2 did not unmarshal")
	}
	if resourceData.Get("T2.0") != "t2a" {
		t.Error("T2 did not unmarshal")
	}
	if resourceData.Get("T2.1") != "t2b" {
		t.Error("T2 did not unmarshal")
	}
	if resourceData.Get("T3") != "2010-04-01" {
		t.Error("T3 did not marshal")
	}
	if resourceData.Get("T4") != 1 {
		t.Error("T4 did not marshal")
	}
}

func TestObjectMarshal(t *testing.T) {
	terraformSchema := map[string]*schema.Schema{
		"Errors": {
			Type: schema.TypeList,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"Links": {
			Type: schema.TypeString,
		},
	}
	data := map[string]interface{}{}
	resourceData := schema.TestResourceDataRaw(t, terraformSchema, data)

	type innerStruct struct {
		Links  map[string]interface{}
		Errors []map[string]interface{}
	}

	testStruct := innerStruct{
		Links: map[string]interface{}{
			"test_users": "https://studio.twilio.com/v2/Flows/FWXX/TestUsers",
			"revisions":  "https://studio.twilio.com/v2/Flows/FWXX/Revisions",
			"executions": "https://studio.twilio.com/v2/Flows/FWXX/Executions",
		},
		Errors: []map[string]interface{}{
			{
				"message":       "some message",
				"property_path": "some property path",
			},
			{
				"message":       "some message 2",
				"property_path": "some property path 2",
			},
		},
	}

	if err := MarshalSchema(resourceData, &testStruct); err != nil {
		t.Errorf("Marshall failed: result '%v'", err)
	}

	if resourceData.Get("Links") != "{\"executions\":\"https://studio.twilio.com/v2/Flows/FWXX/Executions\",\"revisions\":\"https://studio.twilio.com/v2/Flows/FWXX/Revisions\",\"test_users\":\"https://studio.twilio.com/v2/Flows/FWXX/TestUsers\"}" {
		t.Error("Links did not marshal")
	}
	if resourceData.Get("Errors.0") != "{\"message\":\"some message\",\"property_path\":\"some property path\"}" {
		t.Error("Errors.0 did not marshal")
	}
	if resourceData.Get("Errors.1") != "{\"message\":\"some message 2\",\"property_path\":\"some property path 2\"}" {
		t.Error("Errors.1 did not marshal")
	}
}

func TestTimeMarshal(t *testing.T) {

	terraformSchema := map[string]*schema.Schema{
		"T0": {
			Type:     schema.TypeString,
			Required: true,
		},
	}
	data := map[string]interface{}{}
	resourceData := schema.TestResourceDataRaw(t, terraformSchema, data)

	type innerStruct struct {
		T0 *time.Time
	}

	testDate, _ := time.Parse(time.RFC3339, "2021-05-17T01:35:33Z")

	testStruct := innerStruct{T0: &testDate}
	if err := MarshalSchema(resourceData, &testStruct); err != nil {
		t.Errorf("Marshall failed: result '%v'", err)
	}
	if resourceData.Get("T0") != "2021-05-17T01:35:33Z" {
		t.Error("Date T0 did not marshal")
	}
}

func TestFlattenMarshal(t *testing.T) {
	terraformSchema := map[string]*schema.Schema{
		"T1": {
			Type:     schema.TypeString,
			Required: true,
		},
		"T2": {
			Type:     schema.TypeString,
			Required: true,
		},
	}

	data := map[string]interface{}{}
	resourceData := schema.TestResourceDataRaw(t, terraformSchema, data)

	type test2 struct {
		T1 string
		T2 string
	}

	type innerStruct struct {
		T0     *string `provider:",id"`
		Nested test2   `provider:",flatten"`
	}

	testString := "t0"

	testStruct := innerStruct{T0: &testString, Nested: test2{T1: "t1", T2: "t2"}}
	if err := MarshalSchema(resourceData, &testStruct); err != nil {
		t.Errorf("Marshall failed: result '%v'", err)
	}
	if resourceData.Id() != "t0" {
		t.Error("Id did not marshal")
	}
	if resourceData.Get("T1") != "t1" {
		t.Error("T1 did not marshal")
	}
	if resourceData.Get("T2") != "t2" {
		t.Error("T2 did not marshal")
	}
}

func TestListMarshal(t *testing.T) {
	terraformSchema := map[string]*schema.Schema{
		"T1": {
			Type: schema.TypeList,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"T1a": {
						Type:     schema.TypeString,
						Required: true,
					},
					"T1b": {
						Type:     schema.TypeString,
						Required: true,
					},
				},
			},
			Required: true,
		},
	}
	data := map[string]interface{}{}
	resourceData := schema.TestResourceDataRaw(t, terraformSchema, data)

	type test2 struct {
		T1a string
		T1b string
	}

	type innerStruct struct {
		T1 []test2
	}

	testStruct := innerStruct{T1: []test2{{T1a: "r1", T1b: "r2"}, {T1a: "r3", T1b: "r4"}}}
	if err := MarshalSchema(resourceData, &testStruct); err != nil {
		t.Errorf("Marshall failed: result '%v'", err)
	}
	if len(resourceData.Get("T1").([]interface{})) != 2 {
		t.Error("T1 did not unmarshal")
	}
	if resourceData.Get("T1.0.T1a") != "r1" {
		t.Error("T1 did not unmarshal")
	}
	if resourceData.Get("T1.0.T1b") != "r2" {
		t.Error("T1 did not unmarshal")
	}
	if resourceData.Get("T1.1.T1a") != "r3" {
		t.Error("T1 did not unmarshal")
	}
	if resourceData.Get("T1.1.T1b") != "r4" {
		t.Error("T1 did not unmarshal")
	}
}

func TestSidListMarshal(t *testing.T) {
	terraformSchema := map[string]*schema.Schema{
		"T1": {
			Type:     schema.TypeList,
			Elem:     &schema.Schema{Type: schema.TypeString},
			Required: true,
		},
	}
	data := map[string]interface{}{}
	resourceData := schema.TestResourceDataRaw(t, terraformSchema, data)

	type innerStruct struct {
		T1 []AccountSid
	}

	sid1, _ := CreateAccountSid("AC00112233445566778899aabbccddeefe")
	sid2, _ := CreateAccountSid("AC00112233445566778899aabbccddeeff")

	testStruct := innerStruct{T1: []AccountSid{sid1, sid2}}
	if err := MarshalSchema(resourceData, &testStruct); err != nil {
		t.Errorf("Marshall failed: result '%v'", err)
	}
	if len(resourceData.Get("T1").([]interface{})) != 2 {
		t.Error("T1 did not unmarshal")
	}
	if resourceData.Get("T1.0") != "AC00112233445566778899aabbccddeefe" {
		t.Error("T1 did not unmarshal")
	}
	if resourceData.Get("T1.1") != "AC00112233445566778899aabbccddeeff" {
		t.Error("T1 did not unmarshal")
	}
}

func TestMapMarshal(t *testing.T) {
	terraformSchema := map[string]*schema.Schema{
		"T1": {
			Type:     schema.TypeMap,
			Elem:     &schema.Schema{Type: schema.TypeString},
			Required: true,
		},
	}
	data := map[string]interface{}{}
	resourceData := schema.TestResourceDataRaw(t, terraformSchema, data)

	type innerStruct struct {
		T1 map[string]string
	}

	testStruct := innerStruct{T1: map[string]string{"T1a": "r1", "T1b": "r2"}}
	if err := MarshalSchema(resourceData, &testStruct); err != nil {
		t.Errorf("Marshall failed: result '%v'", err)
	}
	if len(resourceData.Get("T1").(map[string]interface{})) != 2 {
		t.Error("T1 did not unmarshal")
	}
	if resourceData.Get("T1.T1a") != "r1" {
		t.Error("T1 did not unmarshal")
	}
	if resourceData.Get("T1.T1b") != "r2" {
		t.Error("T1 did not unmarshal")
	}
}

func TestSidMapMarshal(t *testing.T) {
	terraformSchema := map[string]*schema.Schema{
		"T1": {
			Type:     schema.TypeMap,
			Elem:     &schema.Schema{Type: schema.TypeString},
			Required: true,
		},
	}
	data := map[string]interface{}{}
	resourceData := schema.TestResourceDataRaw(t, terraformSchema, data)

	type innerStruct struct {
		T1 map[string]AccountSid
	}

	sid1, _ := CreateAccountSid("AC00112233445566778899aabbccddeefe")
	sid2, _ := CreateAccountSid("AC00112233445566778899aabbccddeeff")

	testStruct := innerStruct{T1: map[string]AccountSid{"T1a": sid1, "T1b": sid2}}
	if err := MarshalSchema(resourceData, &testStruct); err != nil {
		t.Errorf("Marshall failed: result '%v'", err)
	}
	if len(resourceData.Get("T1").(map[string]interface{})) != 2 {
		t.Error("T1 did not unmarshal")
	}
	if resourceData.Get("T1.T1a") != "AC00112233445566778899aabbccddeefe" {
		t.Error("T1 did not unmarshal")
	}
	if resourceData.Get("T1.T1b") != "AC00112233445566778899aabbccddeeff" {
		t.Error("T1 did not unmarshal")
	}
}

func TestComplexMapMarshal(t *testing.T) {
	terraformSchema := map[string]*schema.Schema{
		"T1": {
			Type:     schema.TypeString,
			Required: true,
		},
	}
	data := map[string]interface{}{}
	resourceData := schema.TestResourceDataRaw(t, terraformSchema, data)

	type nestedStruct struct {
		T1a string
		T1b bool
	}

	type innerStruct struct {
		T1 map[string]nestedStruct
	}

	testStruct := innerStruct{T1: map[string]nestedStruct{"M0": {T1a: "r1", T1b: true}, "M1": {T1a: "r2", T1b: false}}}
	if err := MarshalSchema(resourceData, &testStruct); err != nil {
		t.Errorf("Marshall failed: result '%v'", err)
	}
	if resourceData.Get("T1") != "{\"M0\":{\"T1a\":\"r1\",\"T1b\":true},\"M1\":{\"T1a\":\"r2\",\"T1b\":false}}" {
		t.Error("T1 did not unmarshal")
	}
}

func TestPureJsonMarshal(t *testing.T) {
	terraformSchema := map[string]*schema.Schema{
		"T1": {
			Type:     schema.TypeString,
			Required: true,
		},
	}
	data := map[string]interface{}{}
	resourceData := schema.TestResourceDataRaw(t, terraformSchema, data)

	type innerStruct struct {
		T1 interface{}
	}

	testStruct := innerStruct{T1: map[string]interface{}{"test": "test_value"}}
	if err := MarshalSchema(resourceData, &testStruct); err != nil {
		t.Errorf("Marshall failed: result '%v'", err)
	}

	if resourceData.Get("T1") != "{\"test\":\"test_value\"}" {
		t.Error("T1 did not unmarshal")
	}
}

func TestJsonEncodedListOfObjectsMarshal(t *testing.T) {
	terraformSchema := map[string]*schema.Schema{
		"T1": {
			Type:     schema.TypeList,
			Required: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
	}

	type test2 struct {
		T1 []interface{}
	}

	data := map[string]interface{}{}
	resourceData := schema.TestResourceDataRaw(t, terraformSchema, data)

	m1 := map[string]interface{}{
		"foo": "bar1",
	}
	m2 := map[string]interface{}{
		"foo": "bar2",
	}
	testStruct := test2{
		T1: []interface{}{m1, m2},
	}
	if err := MarshalSchema(resourceData, &testStruct); err != nil {
		t.Errorf("Marshall failed: result '%v'", err)
	}
	if resourceData.Get("T1.0") != "{\"foo\":\"bar1\"}" {
		t.Error("T1 did not unmarshal")
	}
	if resourceData.Get("T1.1") != "{\"foo\":\"bar2\"}" {
		t.Error("T1 did not unmarshal")
	}
}

func TestImplicitNestedMarshal(t *testing.T) {
	terraformSchema := map[string]*schema.Schema{
		"limits_channel_members": {
			Type:     schema.TypeInt,
			Optional: true,
		},
		"notifications_log_enabled": {
			Type:     schema.TypeBool,
			Optional: true,
		},
		"notifications_new_message_sound": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"notifications_recipient_last_name": {
			Type:     schema.TypeString,
			Optional: true,
		},
	}
	data := map[string]interface{}{}
	resourceData := schema.TestResourceDataRaw(t, terraformSchema, data)

	type innerStruct struct {
		Limits        *map[string]interface{} `json:"limits"`
		Notifications *map[string]interface{} `json:"notifications"`
	}

	var recipient interface{} = map[string]interface{}{
		"first_name": "Turk",
		"last_name":  "Andjaydee",
		"email":      nil,
	}

	testStruct := innerStruct{
		Limits: &map[string]interface{}{
			"channel_members": 10,
		},
		Notifications: &map[string]interface{}{
			"log_enabled":       true,
			"new_message_sound": "LOUD NOISES!",
			"recipient":         &recipient,
		},
	}
	if err := MarshalSchema(resourceData, &testStruct); err != nil {
		t.Errorf("Marshall failed: result '%v'", err)
	}
	if resourceData.Get("limits_channel_members") != 10 {
		t.Error("limits_channel_members did not marshal")
	}
	if resourceData.Get("notifications_log_enabled") != true {
		t.Error("notifications_log_enabled did not marshal")
	}
	if resourceData.Get("notifications_new_message_sound") != "LOUD NOISES!" {
		t.Error("notifications_new_message_sound did not marshal")
	}
	if resourceData.Get("notifications_recipient_last_name") != "Andjaydee" {
		t.Error("notifications_recipient_last_name did not marshal")
	}
	if resourceData.Get("notifications_recipient_email") != nil {
		t.Error("notifications_recipient_email did not marshal")
	}
}

func TestSnakeCaseConversion(t *testing.T) {
	testStr := "Integration.FlowSid"
	result := ToSnakeCase(testStr)
	if result != "integration_flow_sid" {
		t.Errorf("ToSnakeCase(%q) = %q, want %q", testStr, result, "integration_flow_sid")
	}
	testStr = "Integration.Flow.Sid"
	result = ToSnakeCase(testStr)
	if result != "integration_flow_sid" {
		t.Errorf("ToSnakeCase(%q) = %q, want %q", testStr, result, "integration_flow_sid")
	}
	testStr = "integration.flow.sid"
	result = ToSnakeCase(testStr)
	if result != "integration_flow_sid" {
		t.Errorf("ToSnakeCase(%q) = %q, want %q", testStr, result, "integration_flow_sid")
	}
	testStr = "integration_flow_sid"
	result = ToSnakeCase(testStr)
	if result != "integration_flow_sid" {
		t.Errorf("ToSnakeCase(%q) = %q, want %q", testStr, result, "integration_flow_sid")
	}
	testStr = "IntegrationChannel123"
	result = ToSnakeCase(testStr)
	if result != "integration_channel123" {
		t.Errorf("ToSnakeCase(%q) = %q, want %q", testStr, result, "integration_channel123")
	}
	testStr = "IntegrationChannelSid"
	result = ToSnakeCase(testStr)
	if result != "integration_channel_sid" {
		t.Errorf("ToSnakeCase(%q) = %q, want %q", testStr, result, "integration_channel_sid")
	}
}
