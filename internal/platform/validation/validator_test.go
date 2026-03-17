package validation

import "testing"

func TestValidatorUsesExpectedMessages(t *testing.T) {
	validator := New()
	errs := FieldErrors{}

	validator.Email("email", "invalid", errs)
	validator.MinLength("password", "123", 8, errs)
	validator.RequiredString("name", "   ", errs)

	if got := errs["email"]; got != "must be a valid email" {
		t.Fatalf("expected email validation message, got %q", got)
	}

	if got := errs["password"]; got != "min length is 8" {
		t.Fatalf("expected password validation message, got %q", got)
	}

	if got := errs["name"]; got != "required" {
		t.Fatalf("expected required validation message, got %q", got)
	}
}
