package core

import "testing"

func TestIntToString(t *testing.T) {
	if got := IntToString(123); got != "123" {
		t.Errorf("IntToString(123) = %q, want %q", got, "123")
	}
	if got := IntToString(0); got != "0" {
		t.Errorf("IntToString(0) = %q, want %q", got, "0")
	}
	if got := IntToString(-123); got != "-123" {
		t.Errorf("IntToString(-123) = %q, want %q", got, "-123")
	}
}

func TestStringToInt(t *testing.T) {
	value, err := StringToInt("123")
	if err != nil {
		t.Fatalf("StringToInt(%q) unexpected error: %v", "123", err)
	}
	if value != 123 {
		t.Errorf("StringToInt(%q) = %d, want %d", "123", value, 123)
	}

	value, err = StringToInt("0")
	if err != nil {
		t.Fatalf("StringToInt(%q) unexpected error: %v", "0", err)
	}
	if value != 0 {
		t.Errorf("StringToInt(%q) = %d, want %d", "0", value, 0)
	}

	value, err = StringToInt("-123")
	if err != nil {
		t.Fatalf("StringToInt(%q) unexpected error: %v", "-123", err)
	}
	if value != -123 {
		t.Errorf("StringToInt(%q) = %d, want %d", "-123", value, -123)
	}

	if _, err = StringToInt("blurg"); err == nil {
		t.Errorf("StringToInt(%q) expected an error, got nil", "blurg")
	}
}
