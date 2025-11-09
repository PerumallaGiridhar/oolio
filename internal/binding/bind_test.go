package binding

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	enlocales "github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	v10 "github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"

	"github.com/PerumallaGiridhar/oolio/internal/validation"
)

type testDTO struct {
	Name  string `json:"name" validate:"required"`
	Count int    `json:"count" validate:"min=1"`
}

func init() {
	// initialize validator and translations for tests
	validation.Validator = v10.New()
	uni := ut.New(enlocales.New())
	tr, _ := uni.GetTranslator("en")
	validation.Translator = tr
	_ = enTranslations.RegisterDefaultTranslations(validation.Validator, validation.Translator)
}

func TestBindAndValidateJSONRequest_Success(t *testing.T) {
	dto := testDTO{}
	payload := map[string]any{"name": "foo", "count": 2}
	b, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(b))

	errs := BindAndValidateJSONRequest(req, &dto)
	if errs != nil {
		t.Fatalf("expected no errors, got %v", errs)
	}
	if dto.Name != "foo" || dto.Count != 2 {
		t.Fatalf("unexpected dto values: %+v", dto)
	}
}

func TestBindAndValidateJSONRequest_UnknownField(t *testing.T) {
	dto := testDTO{}
	// payload contains unknown field `bad` -> should trigger unknown fields error
	payload := map[string]any{"name": "foo", "count": 1, "bad": "x"}
	b, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(b))

	errs := BindAndValidateJSONRequest(req, &dto)
	if errs == nil {
		t.Fatalf("expected error for unknown fields, got nil")
	}
	if _, ok := errs["error"]; !ok {
		t.Fatalf("expected error key in response, got %v", errs)
	}
}

func TestBindAndValidateJSONRequest_ValidationError(t *testing.T) {
	dto := testDTO{}
	// missing required name and count < min
	payload := map[string]any{"count": 0}
	b, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(b))

	errs := BindAndValidateJSONRequest(req, &dto)
	if errs == nil {
		t.Fatalf("expected validation errors, got nil")
	}
	// should contain at least one field error
	if len(errs) == 0 {
		t.Fatalf("expected field errors, got empty map")
	}
}
