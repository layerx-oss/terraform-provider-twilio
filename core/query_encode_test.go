package core

import (
	"slices"
	"testing"
)

func TestQueryEncoderSimple(t *testing.T) {

	type testStruct struct {
		a int
		b int
	}

	testValue := testStruct{a: 2, b: 3}

	queryString, err := QueryEncoder(testValue)
	if err != nil {
		t.Fatalf("QueryEncoder failed: %v", err)
	}
	want := []string{"?a=2&b=3", "?b=3&a=2"}
	if !slices.Contains(want, queryString) {
		t.Errorf("QueryEncoder() = %q, want one of %v", queryString, want)
	}
}

func TestQueryEncoderNil(t *testing.T) {
	if _, err := FormEncoder(nil); err == nil {
		t.Error("FormEncoder(nil) succeeded, want an error")
	}
}

func TestQueryEncoderArray(t *testing.T) {
	testValue := []int{1, 2}
	if _, err := FormEncoder(testValue); err == nil {
		t.Error("FormEncoder([]int) succeeded, want an error")
	}
}
