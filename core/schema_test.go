package core

import (
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestRequiredSchema(t *testing.T) {
	testSchema := AsInt(SchemaRequired)
	if testSchema.Optional {
		t.Error("Optional = true, want false")
	}
	if !testSchema.Required {
		t.Error("Required = false, want true")
	}
	if testSchema.Computed {
		t.Error("Computed = true, want false")
	}
	if testSchema.ForceNew {
		t.Error("ForceNew = true, want false")
	}
}

func TestOptionalSchema(t *testing.T) {
	testSchema := AsInt(SchemaOptional)
	if !testSchema.Optional {
		t.Error("Optional = false, want true")
	}
	if testSchema.Required {
		t.Error("Required = true, want false")
	}
	if testSchema.Computed {
		t.Error("Computed = true, want false")
	}
	if testSchema.ForceNew {
		t.Error("ForceNew = true, want false")
	}
	if testSchema.Sensitive {
		t.Error("Sensitive = true, want false")
	}
}

func TestComputedSchema(t *testing.T) {
	testSchema := AsInt(SchemaComputed)
	if testSchema.Optional {
		t.Error("Optional = true, want false")
	}
	if testSchema.Required {
		t.Error("Required = true, want false")
	}
	if !testSchema.Computed {
		t.Error("Computed = false, want true")
	}
	if testSchema.ForceNew {
		t.Error("ForceNew = true, want false")
	}
	if testSchema.Sensitive {
		t.Error("Sensitive = true, want false")
	}
}

func TestComputedOptionalSchema(t *testing.T) {
	testSchema := AsInt(SchemaComputedOptional)
	if !testSchema.Optional {
		t.Error("Optional = false, want true")
	}
	if testSchema.Required {
		t.Error("Required = true, want false")
	}
	if !testSchema.Computed {
		t.Error("Computed = false, want true")
	}
	if testSchema.ForceNew {
		t.Error("ForceNew = true, want false")
	}
	if testSchema.Sensitive {
		t.Error("Sensitive = true, want false")
	}
}

func TestComputedSensitiveSchema(t *testing.T) {
	testSchema := AsInt(SchemaComputedSensitive)
	if testSchema.Optional {
		t.Error("Optional = true, want false")
	}
	if testSchema.Required {
		t.Error("Required = true, want false")
	}
	if !testSchema.Computed {
		t.Error("Computed = false, want true")
	}
	if testSchema.ForceNew {
		t.Error("ForceNew = true, want false")
	}
	if !testSchema.Sensitive {
		t.Error("Sensitive = false, want true")
	}
}

func TestForceNewRequiredSchemaSchema(t *testing.T) {
	testSchema := AsInt(SchemaForceNewRequired)
	if testSchema.Optional {
		t.Error("Optional = true, want false")
	}
	if !testSchema.Required {
		t.Error("Required = false, want true")
	}
	if testSchema.Computed {
		t.Error("Computed = true, want false")
	}
	if !testSchema.ForceNew {
		t.Error("ForceNew = false, want true")
	}
	if testSchema.Sensitive {
		t.Error("Sensitive = true, want false")
	}
}

func TestForceNewRequiredOptionalSchema(t *testing.T) {
	testSchema := AsInt(SchemaForceNewOptional)
	if !testSchema.Optional {
		t.Error("Optional = false, want true")
	}
	if testSchema.Required {
		t.Error("Required = true, want false")
	}
	if !testSchema.ForceNew {
		t.Error("ForceNew = false, want true")
	}
	if testSchema.Sensitive {
		t.Error("Sensitive = true, want false")
	}
}

func TestInvalidEmptySchema(t *testing.T) {
	defer func() { _ = recover() }()
	AsInt(&options{})
	t.Errorf("Invalid schema allowed")
}

func TestInvalidComputedRequiredSchema(t *testing.T) {
	defer func() { _ = recover() }()
	AsInt(&options{Required: true, Computed: true})
	t.Errorf("Invalid schema allowed")
}

func TestStringSchema(t *testing.T) {
	s := AsString(SchemaRequired)
	if s.Type != schema.TypeString {
		t.Errorf("Type = %v, want %v", s.Type, schema.TypeString)
	}
}

func TestStringBool(t *testing.T) {
	s := AsBool(SchemaRequired)
	if s.Type != schema.TypeBool {
		t.Errorf("Type = %v, want %v", s.Type, schema.TypeBool)
	}
}

func TestIntSchema(t *testing.T) {
	s := AsInt(SchemaRequired)
	if s.Type != schema.TypeInt {
		t.Errorf("Type = %v, want %v", s.Type, schema.TypeInt)
	}
}

func TestFloatSchema(t *testing.T) {
	s := AsFloat(SchemaRequired)
	if s.Type != schema.TypeFloat {
		t.Errorf("Type = %v, want %v", s.Type, schema.TypeFloat)
	}
}

func TestSidSchema(t *testing.T) {
	s := AsSid(&Sid{}, SchemaRequired)
	if s.Type != schema.TypeString {
		t.Errorf("Type = %v, want %v", s.Type, schema.TypeString)
	}

	err := s.ValidateDiagFunc("XX00112233445566778899aabbccddeeff", nil)
	if len(err) != 0 {
		t.Errorf("Sid errored: %v", err)
	}

	err = s.ValidateDiagFunc("abc", nil)
	if len(err) == 0 {
		t.Error("Sid validate did not error")
	}
}

func TestListSchemaScalar(t *testing.T) {
	s := AsList(AsInt(SchemaRequired), SchemaRequired)
	if s.Type != schema.TypeList {
		t.Errorf("Type = %v, want %v", s.Type, schema.TypeList)
	}
	if got := s.Elem.(*schema.Schema).Type; got != schema.TypeInt {
		t.Errorf("Elem.Type = %v, want %v", got, schema.TypeInt)
	}
}

func TestListSchemaComplex(t *testing.T) {
	mapSchema := map[string]*schema.Schema{"test1": AsInt(SchemaRequired)}
	s := AsList(mapSchema, SchemaRequired)
	if s.Type != schema.TypeList {
		t.Errorf("Type = %v, want %v", s.Type, schema.TypeList)
	}
	if got := s.Elem.(*schema.Resource).Schema; !reflect.DeepEqual(got, mapSchema) {
		t.Errorf("Elem.Schema = %#v, want %#v", got, mapSchema)
	}
}
