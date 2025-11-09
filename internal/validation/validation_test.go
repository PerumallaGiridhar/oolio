package validation

import (
	"testing"

	"github.com/go-playground/validator/v10"
)

func TestRegisterTranslations_SetsGlobalTranslator(t *testing.T) {
	Validator = validator.New()
	Translator = nil

	if err := RegisterTranslations(); err != nil {
		t.Fatalf("RegisterTranslations() error = %v", err)
	}

	if Translator == nil {
		t.Fatalf("expected Translator to be set, got nil")
	}

	type Req struct {
		Field string `validate:"required"`
	}

	err := Validator.Struct(Req{})
	if err == nil {
		t.Fatalf("expected validation error, got nil")
	}

	ve, ok := err.(validator.ValidationErrors)
	if !ok || len(ve) == 0 {
		t.Fatalf("expected ValidationErrors with at least one error, got %T %v", err, err)
	}

	msg := ve[0].Translate(Translator)
	if msg == "" {
		t.Fatalf("expected non-empty translated message, got %q", msg)
	}
}
