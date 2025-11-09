package product

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	enlocales "github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	v10 "github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"

	"github.com/PerumallaGiridhar/oolio/internal/data"
	"github.com/PerumallaGiridhar/oolio/internal/validation"
)

func init() {
	validation.Validator = v10.New()
	uni := ut.New(enlocales.New())
	tr, _ := uni.GetTranslator("en")
	validation.Translator = tr
	_ = enTranslations.RegisterDefaultTranslations(validation.Validator, validation.Translator)
	_ = validation.Validator.RegisterValidation("promocode", func(fl v10.FieldLevel) bool {
		return true
	})
}

func TestListProducts(t *testing.T) {
	r := NewRouter()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	// ListProducts uses StatusOk (200) in handlers
	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200 got %d", rr.Code)
	}
	var data []data.Product
	if err := json.Unmarshal(rr.Body.Bytes(), &data); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}
	if len(data) == 0 {
		t.Fatalf("expected non-empty product list")
	}
}

func TestFindProductById_SuccessInvalidAndNotFound(t *testing.T) {
	r := NewRouter()

	// existing id
	req := httptest.NewRequest(http.MethodGet, "/1", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200 got %d", rr.Code)
	}
	var p data.Product
	if err := json.Unmarshal(rr.Body.Bytes(), &p); err != nil {
		t.Fatalf("failed to unmarshal product: %v", err)
	}
	if p.ID != "1" {
		t.Fatalf("expected product id 1 got %v", p.ID)
	}

	// non-existing id
	req2 := httptest.NewRequest(http.MethodGet, "/999", nil)
	rr2 := httptest.NewRecorder()
	r.ServeHTTP(rr2, req2)
	if rr2.Code != http.StatusNotFound {
		t.Fatalf("expected status 404 got %d", rr2.Code)
	}

	// invalid id
	req3 := httptest.NewRequest(http.MethodGet, "/abc", nil)
	rr3 := httptest.NewRecorder()
	r.ServeHTTP(rr3, req3)
	if rr3.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected status 404 got %d", rr3.Code)
	}
}
