package validation

import (
	"fmt"
	"strings"

	playgroundvalidator "github.com/go-playground/validator/v10"
)

type FieldErrors map[string]string

type Validator struct {
	validate *playgroundvalidator.Validate
}

func New() *Validator {
	return &Validator{
		validate: playgroundvalidator.New(),
	}
}

func (errs FieldErrors) Add(field string, message string) {
	if errs == nil {
		return
	}

	if _, exists := errs[field]; exists {
		return
	}

	errs[field] = message
}

func (errs FieldErrors) HasAny() bool {
	return len(errs) > 0
}

func (v *Validator) RequiredString(field string, value string, errs FieldErrors) {
	if err := v.validate.Var(strings.TrimSpace(value), "required"); err != nil {
		errs.Add(field, "required")
	}
}

func (v *Validator) MinLength(field string, value string, min int, errs FieldErrors) {
	if err := v.validate.Var(strings.TrimSpace(value), fmt.Sprintf("min=%d", min)); err != nil {
		errs.Add(field, fmt.Sprintf("min length is %d", min))
	}
}

func (v *Validator) Email(field string, value string, errs FieldErrors) {
	trimmed := strings.TrimSpace(value)
	if err := v.validate.Var(trimmed, "required"); err != nil {
		errs.Add(field, "required")
		return
	}

	if err := v.validate.Var(trimmed, "email"); err != nil {
		errs.Add(field, "must be a valid email")
	}
}
