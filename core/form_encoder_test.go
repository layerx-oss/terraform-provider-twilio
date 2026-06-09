package core

import (
	"net/url"
	"reflect"
	"testing"
)

func TestProcessFormValueBool(t *testing.T) {
	form := url.Values{}

	type testStruct struct {
		a bool
	}

	testValue := testStruct{a: true}

	err := processParamValue("form", &form, reflect.TypeOf(testValue).Field(0), reflect.ValueOf(testValue).Field(0))
	if err != nil {
		t.Fatalf("ProcessFormValue failed: %v", err)
	}
	if got, want := form["a"], []string{"true"}; !reflect.DeepEqual(got, want) {
		t.Errorf("form[\"a\"] = %#v, want %#v", got, want)
	}
}

func TestProcessFormValueInt(t *testing.T) {
	form := url.Values{}

	type testStruct struct {
		a int
	}

	testValue := testStruct{a: 1}

	err := processParamValue("form", &form, reflect.TypeOf(testValue).Field(0), reflect.ValueOf(testValue).Field(0))
	if err != nil {
		t.Fatalf("ProcessFormValue failed: %v", err)
	}
	if got, want := form["a"], []string{"1"}; !reflect.DeepEqual(got, want) {
		t.Errorf("form[\"a\"] = %#v, want %#v", got, want)
	}
}

func TestProcessFormValueFloat(t *testing.T) {
	form := url.Values{}

	type testStruct struct {
		a float64
	}

	testValue := testStruct{a: 1}

	err := processParamValue("form", &form, reflect.TypeOf(testValue).Field(0), reflect.ValueOf(testValue).Field(0))
	if err != nil {
		t.Fatalf("ProcessFormValue failed: %v", err)
	}
	if got, want := form["a"], []string{"1"}; !reflect.DeepEqual(got, want) {
		t.Errorf("form[\"a\"] = %#v, want %#v", got, want)
	}
}

func TestProcessFormValueString(t *testing.T) {
	form := url.Values{}

	type testStruct struct {
		a string
	}

	testValue := testStruct{a: "str"}

	err := processParamValue("form", &form, reflect.TypeOf(testValue).Field(0), reflect.ValueOf(testValue).Field(0))
	if err != nil {
		t.Fatalf("ProcessFormValue failed: %v", err)
	}
	if got, want := form["a"], []string{"str"}; !reflect.DeepEqual(got, want) {
		t.Errorf("form[\"a\"] = %#v, want %#v", got, want)
	}
}

func TestProcessFormValueIntList(t *testing.T) {
	form := url.Values{}

	type testStruct struct {
		a []int
	}

	testValue := testStruct{a: []int{1, 2, 3}}

	err := processParamValue("form", &form, reflect.TypeOf(testValue).Field(0), reflect.ValueOf(testValue).Field(0))
	if err != nil {
		t.Fatalf("ProcessFormValue failed: %v", err)
	}
	if got, want := form["a"], []string{"1", "2", "3"}; !reflect.DeepEqual(got, want) {
		t.Errorf("form[\"a\"] = %#v, want %#v", got, want)
	}
}

func TestProcessFormValueArray(t *testing.T) {
	form := url.Values{}

	type testStruct struct {
		a [1]int
	}

	testValue := testStruct{a: [1]int{1}}

	err := processParamValue("form", &form, reflect.TypeOf(testValue).Field(0), reflect.ValueOf(testValue).Field(0))
	if err == nil {
		t.Error("ProcessFormValue succeeded, want an error")
	}
}

func TestProcessFormValueBasic(t *testing.T) {
	form := url.Values{}

	type testStruct struct {
		A Sid
	}

	sid, _ := CreateSid("XX00112233445566778899aabbccddeeff")
	testValue := testStruct{A: sid}

	err := processParamValue("form", &form, reflect.TypeOf(testValue).Field(0), reflect.ValueOf(testValue).Field(0))
	if err != nil {
		t.Fatalf("ProcessFormValue failed: %v", err)
	}
	if got, want := form["A"], []string{"XX00112233445566778899aabbccddeeff"}; !reflect.DeepEqual(got, want) {
		t.Errorf("form[\"A\"] = %#v, want %#v", got, want)
	}
}

func TestProcessFormValuePtrNil(t *testing.T) {
	form := url.Values{}

	type testStruct struct {
		A *Sid
	}

	testValue := testStruct{A: nil}

	err := processParamValue("form", &form, reflect.TypeOf(testValue).Field(0), reflect.ValueOf(testValue).Field(0))
	if err == nil {
		t.Error("ProcessFormValue succeeded, want an error")
	}
}

func TestProcessFormValuePtrBasic(t *testing.T) {
	form := url.Values{}

	type testStruct struct {
		A *Sid
	}

	sid, _ := CreateSid("XX00112233445566778899aabbccddeeff")
	testValue := testStruct{A: &sid}

	err := processParamValue("form", &form, reflect.TypeOf(testValue).Field(0), reflect.ValueOf(testValue).Field(0))
	if err != nil {
		t.Fatalf("ProcessFormValue failed: %v", err)
	}
	if got, want := form["A"], []string{"XX00112233445566778899aabbccddeeff"}; !reflect.DeepEqual(got, want) {
		t.Errorf("form[\"A\"] = %#v, want %#v", got, want)
	}
}

func TestProcessFormValueJson(t *testing.T) {
	form := url.Values{}

	type testStruct struct {
		A interface{}
	}

	testValue := testStruct{A: map[string]interface{}{"a": 1, "b": "str"}}

	err := processParamValue("form", &form, reflect.TypeOf(testValue).Field(0), reflect.ValueOf(testValue).Field(0))
	if err != nil {
		t.Fatalf("ProcessFormValue failed: %v", err)
	}
	if got, want := form["A"], []string{"{\"a\":1,\"b\":\"str\"}"}; !reflect.DeepEqual(got, want) {
		t.Errorf("form[\"A\"] = %#v, want %#v", got, want)
	}
}

func TestProcessFormValueName(t *testing.T) {
	form := url.Values{}

	type testStruct struct {
		a int `form:"b"`
	}

	testValue := testStruct{a: 1}

	err := processParamValue("form", &form, reflect.TypeOf(testValue).Field(0), reflect.ValueOf(testValue).Field(0))
	if err != nil {
		t.Fatalf("ProcessFormValue failed: %v", err)
	}
	if got, want := form["b"], []string{"1"}; !reflect.DeepEqual(got, want) {
		t.Errorf("form[\"b\"] = %#v, want %#v", got, want)
	}
}

func TestProcessFormValueIgnore(t *testing.T) {
	form := url.Values{}

	type testStruct struct {
		a int `form:",ignore"`
	}

	testValue := testStruct{a: 1}

	err := processParamValue("form", &form, reflect.TypeOf(testValue).Field(0), reflect.ValueOf(testValue).Field(0))
	if err != nil {
		t.Fatalf("ProcessFormValue failed: %v", err)
	}
	if _, ok := form["a"]; ok {
		t.Error("form is not empty, want empty")
	}
}

func TestProcessFormValueOmitEmpty(t *testing.T) {
	form := url.Values{}

	type testStruct struct {
		A *Sid `form:",omitempty"`
	}

	testValue := testStruct{A: nil}

	err := processParamValue("form", &form, reflect.TypeOf(testValue).Field(0), reflect.ValueOf(testValue).Field(0))
	if err != nil {
		t.Fatalf("ProcessFormValue failed: %v", err)
	}
	if _, ok := form["A"]; ok {
		t.Error("form is not empty, want empty")
	}
}

func TestProcessFormValueFlatten(t *testing.T) {
	form := url.Values{}

	type testStruct2 struct {
		B int
	}

	type testStruct struct {
		A testStruct2 `form:",flatten"`
	}

	testValue := testStruct{A: testStruct2{B: 3}}

	err := processParamValue("form", &form, reflect.TypeOf(testValue).Field(0), reflect.ValueOf(testValue).Field(0))
	if err != nil {
		t.Fatalf("ProcessFormValue failed: %v", err)
	}
	if got, want := form["B"], []string{"3"}; !reflect.DeepEqual(got, want) {
		t.Errorf("form[\"B\"] = %#v, want %#v", got, want)
	}
}

func TestProcessFormValueNoFlatten(t *testing.T) {
	form := url.Values{}

	type testStruct2 struct {
		B int
	}

	type testStruct struct {
		A testStruct2
	}

	testValue := testStruct{A: testStruct2{B: 3}}

	err := processParamValue("form", &form, reflect.TypeOf(testValue).Field(0), reflect.ValueOf(testValue).Field(0))
	if err == nil {
		t.Error("ProcessFormValue succeeded, want an error")
	}
}

func TestFormEncoderSimple(t *testing.T) {

	type testStruct struct {
		a int
		b int
	}

	testValue := testStruct{a: 2, b: 3}

	form, err := FormEncoder(testValue)
	if err != nil {
		t.Fatalf("FormEncoder failed: %v", err)
	}
	if got, want := form["a"], []string{"2"}; !reflect.DeepEqual(got, want) {
		t.Errorf("form[\"a\"] = %#v, want %#v", got, want)
	}
	if got, want := form["b"], []string{"3"}; !reflect.DeepEqual(got, want) {
		t.Errorf("form[\"b\"] = %#v, want %#v", got, want)
	}
}

func TestFormEncoderNil(t *testing.T) {
	if _, err := FormEncoder(nil); err == nil {
		t.Error("FormEncoder(nil) succeeded, want an error")
	}
}

func TestFormEncoderArray(t *testing.T) {
	testValue := []int{1, 2}
	if _, err := FormEncoder(testValue); err == nil {
		t.Error("FormEncoder([]int) succeeded, want an error")
	}
}
