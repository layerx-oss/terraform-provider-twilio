package core

import (
	"encoding/json"
	"testing"
)

func TestNullableStringGet(t *testing.T) {

	T1 := NullableString{Valid: true, Value: "str"}

	value, ok := T1.Get()
	if !ok {
		t.Error("got value is not valid")
	}
	if value != "[str]" {
		t.Errorf("Get() value = %q, want %q", value, "[str]")
	}
}

func TestNullableStringGetEmpty(t *testing.T) {

	T1 := NullableString{Valid: true, Value: ""}

	value, ok := T1.Get()
	if !ok {
		t.Error("got value is not valid")
	}
	if value != "[]" {
		t.Errorf("Get() value = %q, want %q", value, "[]")
	}
}

func TestNullableStringGetNil(t *testing.T) {

	T1 := NullableString{Valid: false, Value: ""}

	_, ok := T1.Get()
	if ok {
		t.Error("got value is valid")
	}
}

func TestNullableStringGetNativePresentation(t *testing.T) {

	T1 := NullableString{Valid: true, Value: "str"}

	value, ok := T1.GetNativePresentation()
	if !ok {
		t.Error("got value is not valid")
	}
	if value != "str" {
		t.Errorf("GetNativePresentation() value = %q, want %q", value, "str")
	}
}

func TestNullableStringGetNativePresentationEmpty(t *testing.T) {

	T1 := NullableString{Valid: true, Value: ""}

	value, ok := T1.GetNativePresentation()
	if !ok {
		t.Error("got value is not valid")
	}
	if value != "" {
		t.Errorf("GetNativePresentation() value = %q, want %q", value, "")
	}
}

func TestNullableStringGetNativePresentationNil(t *testing.T) {

	T1 := NullableString{Valid: false, Value: ""}

	_, ok := T1.GetNativePresentation()
	if ok {
		t.Error("got value is valid")
	}
}

func TestNullableStringSet(t *testing.T) {
	var test NullableString

	err := test.Set("[str]")
	if err != nil {
		t.Fatalf("test did not set: %v", err)
	}
	if !test.Valid {
		t.Error("marshaled value is not valid")
	}
	if test.Value != "str" {
		t.Errorf("Value = %q, want %q", test.Value, "str")
	}
}

func TestNullableStringSetEmpty(t *testing.T) {
	var test NullableString

	err := test.Set("[]")
	if err != nil {
		t.Fatalf("test did not set: %v", err)
	}
	if !test.Valid {
		t.Error("marshaled value is not valid")
	}
	if test.Value != "" {
		t.Errorf("Value = %q, want %q", test.Value, "")
	}
}

func TestNullableStringSetNil(t *testing.T) {
	var test NullableString

	err := test.Set("")
	if err != nil {
		t.Fatalf("test did not set: %v", err)
	}
	if test.Valid {
		t.Error("marshaled value is valid")
	}
}

func TestNullableStringMarshal(t *testing.T) {

	type test struct {
		T1 NullableString
	}

	jsonValue := test{T1: NullableString{Valid: true, Value: "str"}}

	jsonString, err := json.Marshal(jsonValue)
	if err != nil {
		t.Fatalf("test did not marshal: %v", err)
	}
	if got := string(jsonString); got != "{\"T1\":\"str\"}" {
		t.Errorf("marshaled value = %q, want %q", got, "{\"T1\":\"str\"}")
	}
}

func TestNullableStringMarshalEmpty(t *testing.T) {

	type test struct {
		T1 NullableString
	}

	jsonValue := test{T1: NullableString{Valid: true, Value: ""}}

	jsonString, err := json.Marshal(jsonValue)
	if err != nil {
		t.Fatalf("test did not marshal: %v", err)
	}
	if got := string(jsonString); got != "{\"T1\":\"\"}" {
		t.Errorf("marshaled value = %q, want %q", got, "{\"T1\":\"\"}")
	}
}

func TestNullableStringMarshalNil(t *testing.T) {

	type test struct {
		T1 NullableString
	}

	jsonValue := test{T1: NullableString{Valid: false, Value: ""}}

	jsonString, err := json.Marshal(jsonValue)
	if err != nil {
		t.Fatalf("test did not marshal: %v", err)
	}
	if got := string(jsonString); got != "{\"T1\":null}" {
		t.Errorf("marshaled value = %q, want %q", got, "{\"T1\":null}")
	}
}

func TestNullableStringUmarshal(t *testing.T) {

	type test struct {
		T1 NullableString
	}

	var jsonValue test

	err := json.Unmarshal([]byte("{\"T1\":\"str\"}"), &jsonValue)
	if err != nil {
		t.Fatalf("test did not unmarshal: %v", err)
	}
	if !jsonValue.T1.Valid {
		t.Error("marshaled value is not valid")
	}
	if jsonValue.T1.Value != "str" {
		t.Errorf("Value = %q, want %q", jsonValue.T1.Value, "str")
	}
}

func TestNullableStringUmarshalEmpty(t *testing.T) {

	type test struct {
		T1 NullableString
	}

	var jsonValue test

	err := json.Unmarshal([]byte("{\"T1\":\"\"}"), &jsonValue)
	if err != nil {
		t.Fatalf("test did not unmarshal: %v", err)
	}
	if !jsonValue.T1.Valid {
		t.Error("marshaled value is not valid")
	}
	if jsonValue.T1.Value != "" {
		t.Errorf("Value = %q, want %q", jsonValue.T1.Value, "")
	}
}

func TestNullableStringUmarshalNil(t *testing.T) {

	type test struct {
		T1 NullableString
	}

	var jsonValue test

	err := json.Unmarshal([]byte("{\"T1\":null}"), &jsonValue)
	if err != nil {
		t.Fatalf("test did not unmarshal: %v", err)
	}
	if jsonValue.T1.Valid {
		t.Error("marshaled value is valid")
	}
}

// ---

func TestNullableBoolGet(t *testing.T) {

	T1 := NullableBool{Valid: true, Value: true}

	value, ok := T1.Get()
	if !ok {
		t.Error("got value is not valid")
	}
	if value != "true" {
		t.Errorf("Get() value = %q, want %q", value, "true")
	}
}

func TestNullableBoolGetFalse(t *testing.T) {

	T1 := NullableBool{Valid: true, Value: false}

	value, ok := T1.Get()
	if !ok {
		t.Error("got value is not valid")
	}
	if value != "false" {
		t.Errorf("Get() value = %q, want %q", value, "false")
	}
}

func TestNullableBoolGetNil(t *testing.T) {

	T1 := NullableBool{Valid: false, Value: false}

	_, ok := T1.Get()
	if ok {
		t.Error("got value is valid")
	}
}

func TestNullableBoolGetNativePresentation(t *testing.T) {

	T1 := NullableBool{Valid: true, Value: true}

	value, ok := T1.GetNativePresentation()
	if !ok {
		t.Error("got value is not valid")
	}
	if value != true {
		t.Errorf("GetNativePresentation() value = %v, want %v", value, true)
	}
}

func TestNullableBoolGetNativePresentationFalse(t *testing.T) {

	T1 := NullableBool{Valid: true, Value: false}

	value, ok := T1.GetNativePresentation()
	if !ok {
		t.Error("got value is not valid")
	}
	if value != false {
		t.Errorf("GetNativePresentation() value = %v, want %v", value, false)
	}
}

func TestNullableBoolGetNativePresentationNil(t *testing.T) {

	T1 := NullableBool{Valid: false, Value: false}

	_, ok := T1.GetNativePresentation()
	if ok {
		t.Error("got value is valid")
	}
}

func TestNullableBoolSet(t *testing.T) {
	var test NullableBool

	err := test.Set("true")
	if err != nil {
		t.Fatalf("test did not set: %v", err)
	}
	if !test.Valid {
		t.Error("marshaled value is not valid")
	}
	if test.Value != true {
		t.Errorf("Value = %v, want %v", test.Value, true)
	}
}

func TestNullableBoolSetFalse(t *testing.T) {
	var test NullableBool

	err := test.Set("false")
	if err != nil {
		t.Fatalf("test did not set: %v", err)
	}
	if !test.Valid {
		t.Error("marshaled value is not valid")
	}
	if test.Value != false {
		t.Errorf("Value = %v, want %v", test.Value, false)
	}
}

func TestNullableBoolSetNil(t *testing.T) {
	var test NullableBool

	err := test.Set("")
	if err != nil {
		t.Fatalf("test did not set: %v", err)
	}
	if test.Valid {
		t.Error("marshaled value is valid")
	}
}

func TestNullableBoolMarshal(t *testing.T) {

	type test struct {
		T1 NullableBool
	}

	jsonValue := test{T1: NullableBool{Valid: true, Value: true}}

	jsonString, err := json.Marshal(jsonValue)
	if err != nil {
		t.Fatalf("test did not marshal: %v", err)
	}
	if got := string(jsonString); got != "{\"T1\":true}" {
		t.Errorf("marshaled value = %q, want %q", got, "{\"T1\":true}")
	}
}

func TestNullableBoolMarshalFalse(t *testing.T) {

	type test struct {
		T1 NullableBool
	}

	jsonValue := test{T1: NullableBool{Valid: true, Value: false}}

	jsonString, err := json.Marshal(jsonValue)
	if err != nil {
		t.Fatalf("test did not marshal: %v", err)
	}
	if got := string(jsonString); got != "{\"T1\":false}" {
		t.Errorf("marshaled value = %q, want %q", got, "{\"T1\":false}")
	}
}

func TestNullableBoolMarshalNil(t *testing.T) {

	type test struct {
		T1 NullableBool
	}

	jsonValue := test{T1: NullableBool{Valid: false, Value: false}}

	jsonString, err := json.Marshal(jsonValue)
	if err != nil {
		t.Fatalf("test did not marshal: %v", err)
	}
	if got := string(jsonString); got != "{\"T1\":null}" {
		t.Errorf("marshaled value = %q, want %q", got, "{\"T1\":null}")
	}
}

func TestNullableBoolUmarshal(t *testing.T) {

	type test struct {
		T1 NullableBool
	}

	var jsonValue test

	err := json.Unmarshal([]byte("{\"T1\":true}"), &jsonValue)
	if err != nil {
		t.Fatalf("test did not unmarshal: %v", err)
	}
	if !jsonValue.T1.Valid {
		t.Error("marshaled value is not valid")
	}
	if jsonValue.T1.Value != true {
		t.Errorf("Value = %v, want %v", jsonValue.T1.Value, true)
	}
}

func TestNullableBoolUmarshalFalse(t *testing.T) {

	type test struct {
		T1 NullableBool
	}

	var jsonValue test

	err := json.Unmarshal([]byte("{\"T1\":false}"), &jsonValue)
	if err != nil {
		t.Fatalf("test did not unmarshal: %v", err)
	}
	if !jsonValue.T1.Valid {
		t.Error("marshaled value is not valid")
	}
	if jsonValue.T1.Value != false {
		t.Errorf("Value = %v, want %v", jsonValue.T1.Value, false)
	}
}

func TestNullableBoolUmarshalNil(t *testing.T) {

	type test struct {
		T1 NullableBool
	}

	var jsonValue test

	err := json.Unmarshal([]byte("{\"T1\":null}"), &jsonValue)
	if err != nil {
		t.Fatalf("test did not unmarshal: %v", err)
	}
	if jsonValue.T1.Valid {
		t.Error("marshaled value is valid")
	}
}

// ---

func TestNullableIntGet(t *testing.T) {

	T1 := NullableInt{Valid: true, Value: 1}

	value, ok := T1.Get()
	if !ok {
		t.Error("got value is not valid")
	}
	if value != "1" {
		t.Errorf("Get() value = %q, want %q", value, "1")
	}
}

func TestNullableIntGet0(t *testing.T) {

	T1 := NullableInt{Valid: true, Value: 0}

	value, ok := T1.Get()
	if !ok {
		t.Error("got value is not valid")
	}
	if value != "0" {
		t.Errorf("Get() value = %q, want %q", value, "0")
	}
}

func TestNullableIntGetNil(t *testing.T) {

	T1 := NullableInt{Valid: false, Value: 0}

	_, ok := T1.Get()
	if ok {
		t.Error("got value is valid")
	}
}

func TestNullableIntGetNativePresentation(t *testing.T) {

	T1 := NullableInt{Valid: true, Value: 1}

	value, ok := T1.GetNativePresentation()
	if !ok {
		t.Error("got value is not valid")
	}
	if value != int64(1) {
		t.Errorf("GetNativePresentation() value = %v, want %v", value, int64(1))
	}
}

func TestNullableIntGetNativePresentation0(t *testing.T) {

	T1 := NullableInt{Valid: true, Value: 0}

	value, ok := T1.GetNativePresentation()
	if !ok {
		t.Error("got value is not valid")
	}
	if value != int64(0) {
		t.Errorf("GetNativePresentation() value = %v, want %v", value, int64(0))
	}
}

func TestNullableIntGetNativePresentationNil(t *testing.T) {

	T1 := NullableInt{Valid: false, Value: 0}

	_, ok := T1.GetNativePresentation()
	if ok {
		t.Error("got value is valid")
	}
}

func TestNullableIntSet(t *testing.T) {
	var test NullableInt

	err := test.Set("1")
	if err != nil {
		t.Fatalf("test did not set: %v", err)
	}
	if !test.Valid {
		t.Error("marshaled value is not valid")
	}
	if test.Value != int64(1) {
		t.Errorf("Value = %v, want %v", test.Value, int64(1))
	}
}

func TestNullableIntSet0(t *testing.T) {
	var test NullableInt

	err := test.Set("0")
	if err != nil {
		t.Fatalf("test did not set: %v", err)
	}
	if !test.Valid {
		t.Error("marshaled value is not valid")
	}
	if test.Value != int64(0) {
		t.Errorf("Value = %v, want %v", test.Value, int64(0))
	}
}

func TestNullableIntSetNil(t *testing.T) {
	var test NullableInt

	err := test.Set("")
	if err != nil {
		t.Fatalf("test did not set: %v", err)
	}
	if test.Valid {
		t.Error("marshaled value is valid")
	}
}

func TestNullableIntMarshal(t *testing.T) {

	type test struct {
		T1 NullableInt
	}

	jsonValue := test{T1: NullableInt{Valid: true, Value: 1}}

	jsonString, err := json.Marshal(jsonValue)
	if err != nil {
		t.Fatalf("test did not marshal: %v", err)
	}
	if got := string(jsonString); got != "{\"T1\":1}" {
		t.Errorf("marshaled value = %q, want %q", got, "{\"T1\":1}")
	}
}

func TestNullableIntMarshal0(t *testing.T) {

	type test struct {
		T1 NullableInt
	}

	jsonValue := test{T1: NullableInt{Valid: true, Value: 0}}

	jsonString, err := json.Marshal(jsonValue)
	if err != nil {
		t.Fatalf("test did not marshal: %v", err)
	}
	if got := string(jsonString); got != "{\"T1\":0}" {
		t.Errorf("marshaled value = %q, want %q", got, "{\"T1\":0}")
	}
}

func TestNullableIntMarshalNil(t *testing.T) {

	type test struct {
		T1 NullableInt
	}

	jsonValue := test{T1: NullableInt{Valid: false, Value: 0}}

	jsonString, err := json.Marshal(jsonValue)
	if err != nil {
		t.Fatalf("test did not marshal: %v", err)
	}
	if got := string(jsonString); got != "{\"T1\":null}" {
		t.Errorf("marshaled value = %q, want %q", got, "{\"T1\":null}")
	}
}

func TestNullableIntUmarshal(t *testing.T) {

	type test struct {
		T1 NullableInt
	}

	var jsonValue test

	err := json.Unmarshal([]byte("{\"T1\":1}"), &jsonValue)
	if err != nil {
		t.Fatalf("test did not unmarshal: %v", err)
	}
	if !jsonValue.T1.Valid {
		t.Error("marshaled value is not valid")
	}
	if jsonValue.T1.Value != int64(1) {
		t.Errorf("Value = %v, want %v", jsonValue.T1.Value, int64(1))
	}
}

func TestNullableIntUmarshal0(t *testing.T) {

	type test struct {
		T1 NullableInt
	}

	var jsonValue test

	err := json.Unmarshal([]byte("{\"T1\":0}"), &jsonValue)
	if err != nil {
		t.Fatalf("test did not unmarshal: %v", err)
	}
	if !jsonValue.T1.Valid {
		t.Error("marshaled value is not valid")
	}
	if jsonValue.T1.Value != int64(0) {
		t.Errorf("Value = %v, want %v", jsonValue.T1.Value, int64(0))
	}
}

func TestNullableIntUmarshalNil(t *testing.T) {

	type test struct {
		T1 NullableInt
	}

	var jsonValue test

	err := json.Unmarshal([]byte("{\"T1\":null}"), &jsonValue)
	if err != nil {
		t.Fatalf("test did not unmarshal: %v", err)
	}
	if jsonValue.T1.Valid {
		t.Error("marshaled value is valid")
	}
}

// ---

func TestNullableFloatGet(t *testing.T) {

	T1 := NullableFloat{Valid: true, Value: 1}

	value, ok := T1.Get()
	if !ok {
		t.Error("got value is not valid")
	}
	if value != "1" {
		t.Errorf("Get() value = %q, want %q", value, "1")
	}
}

func TestNullableFloatGet0(t *testing.T) {

	T1 := NullableFloat{Valid: true, Value: 0}

	value, ok := T1.Get()
	if !ok {
		t.Error("got value is not valid")
	}
	if value != "0" {
		t.Errorf("Get() value = %q, want %q", value, "0")
	}
}

func TestNullableFloatGetNil(t *testing.T) {

	T1 := NullableFloat{Valid: false, Value: 0}

	_, ok := T1.Get()
	if ok {
		t.Error("got value is valid")
	}
}

func TestNullableFloatGetNativePresentation(t *testing.T) {

	T1 := NullableFloat{Valid: true, Value: 1}

	value, ok := T1.GetNativePresentation()
	if !ok {
		t.Error("got value is not valid")
	}
	if value != 1. {
		t.Errorf("GetNativePresentation() value = %v, want %v", value, 1.)
	}
}

func TestNullableFloatGetNativePresentation0(t *testing.T) {

	T1 := NullableFloat{Valid: true, Value: 0}

	value, ok := T1.GetNativePresentation()
	if !ok {
		t.Error("got value is not valid")
	}
	if value != 0. {
		t.Errorf("GetNativePresentation() value = %v, want %v", value, 0.)
	}
}

func TestNullableFloatGetNativePresentationNil(t *testing.T) {

	T1 := NullableFloat{Valid: false, Value: 0}

	_, ok := T1.GetNativePresentation()
	if ok {
		t.Error("got value is valid")
	}
}

func TestNullableFloatSet(t *testing.T) {
	var test NullableFloat

	err := test.Set("1")
	if err != nil {
		t.Fatalf("test did not set: %v", err)
	}
	if !test.Valid {
		t.Error("marshaled value is not valid")
	}
	if test.Value != 1. {
		t.Errorf("Value = %v, want %v", test.Value, 1.)
	}
}

func TestNullableFloatSet0(t *testing.T) {
	var test NullableFloat

	err := test.Set("0")
	if err != nil {
		t.Fatalf("test did not set: %v", err)
	}
	if !test.Valid {
		t.Error("marshaled value is not valid")
	}
	if test.Value != 0. {
		t.Errorf("Value = %v, want %v", test.Value, 0.)
	}
}

func TestNullableFloatSetNil(t *testing.T) {
	var test NullableFloat

	err := test.Set("")
	if err != nil {
		t.Fatalf("test did not set: %v", err)
	}
	if test.Valid {
		t.Error("marshaled value is valid")
	}
}

func TestNullableFloatMarshal(t *testing.T) {

	type test struct {
		T1 NullableFloat
	}

	jsonValue := test{T1: NullableFloat{Valid: true, Value: 1}}

	jsonString, err := json.Marshal(jsonValue)
	if err != nil {
		t.Fatalf("test did not marshal: %v", err)
	}
	if got := string(jsonString); got != "{\"T1\":1}" {
		t.Errorf("marshaled value = %q, want %q", got, "{\"T1\":1}")
	}
}

func TestNullableFloatMarshal0(t *testing.T) {

	type test struct {
		T1 NullableFloat
	}

	jsonValue := test{T1: NullableFloat{Valid: true, Value: 0}}

	jsonString, err := json.Marshal(jsonValue)
	if err != nil {
		t.Fatalf("test did not marshal: %v", err)
	}
	if got := string(jsonString); got != "{\"T1\":0}" {
		t.Errorf("marshaled value = %q, want %q", got, "{\"T1\":0}")
	}
}

func TestNullableFloatMarshalNil(t *testing.T) {

	type test struct {
		T1 NullableFloat
	}

	jsonValue := test{T1: NullableFloat{Valid: false, Value: 0}}

	jsonString, err := json.Marshal(jsonValue)
	if err != nil {
		t.Fatalf("test did not marshal: %v", err)
	}
	if got := string(jsonString); got != "{\"T1\":null}" {
		t.Errorf("marshaled value = %q, want %q", got, "{\"T1\":null}")
	}
}

func TestNullableFloatUmarshal(t *testing.T) {

	type test struct {
		T1 NullableFloat
	}

	var jsonValue test

	err := json.Unmarshal([]byte("{\"T1\":1}"), &jsonValue)
	if err != nil {
		t.Fatalf("test did not unmarshal: %v", err)
	}
	if !jsonValue.T1.Valid {
		t.Error("marshaled value is not valid")
	}
	if jsonValue.T1.Value != 1. {
		t.Errorf("Value = %v, want %v", jsonValue.T1.Value, 1.)
	}
}

func TestNullableFloatUmarshal0(t *testing.T) {

	type test struct {
		T1 NullableFloat
	}

	var jsonValue test

	err := json.Unmarshal([]byte("{\"T1\":0}"), &jsonValue)
	if err != nil {
		t.Fatalf("test did not unmarshal: %v", err)
	}
	if !jsonValue.T1.Valid {
		t.Error("marshaled value is not valid")
	}
	if jsonValue.T1.Value != 0. {
		t.Errorf("Value = %v, want %v", jsonValue.T1.Value, 0.)
	}
}

func TestNullableFloatUmarshalNil(t *testing.T) {

	type test struct {
		T1 NullableFloat
	}

	var jsonValue test

	err := json.Unmarshal([]byte("{\"T1\":null}"), &jsonValue)
	if err != nil {
		t.Fatalf("test did not unmarshal: %v", err)
	}
	if jsonValue.T1.Valid {
		t.Error("marshaled value is valid")
	}
}
