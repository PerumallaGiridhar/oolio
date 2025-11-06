package routes

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	enlocales "github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	v10 "github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"

	"github.com/PerumallaGiridhar/oolio/internal/data"
	"github.com/PerumallaGiridhar/oolio/internal/routes/order"
	"github.com/PerumallaGiridhar/oolio/internal/validation"
)

func TestStatsEndpoint(t *testing.T) {
	r := NewRouter()
	req := httptest.NewRequest(http.MethodGet, "/stats", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200 got %d", rr.Code)
	}
	body, _ := io.ReadAll(rr.Body)
	type stats struct {
		Alloc      string
		TotalAlloc string
		Sys        string
		NumGC      string
	}
	var s stats
	if err := json.Unmarshal(body, &s); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
}

func init() {
	// initialize validator and translator used by binding package
	validation.Validator = v10.New()
	uni := ut.New(enlocales.New())
	tr, _ := uni.GetTranslator("en")
	validation.Translator = tr
	// register default translations to avoid nil usage in Translate
	_ = enTranslations.RegisterDefaultTranslations(validation.Validator, validation.Translator)
	// register a simple promocode validator used by order DTO tags
	_ = validation.Validator.RegisterValidation("promocode", func(fl v10.FieldLevel) bool {
		// accept empty or any string in tests
		return true
	})
}

func TestListProducts(t *testing.T) {
	r := NewRouter()
	req := httptest.NewRequest(http.MethodGet, "/api/product/", nil)
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
	req := httptest.NewRequest(http.MethodGet, "/api/product/1", nil)
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
	req2 := httptest.NewRequest(http.MethodGet, "/api/product/999", nil)
	rr2 := httptest.NewRecorder()
	r.ServeHTTP(rr2, req2)
	if rr2.Code != http.StatusNotFound {
		t.Fatalf("expected status 404 got %d", rr2.Code)
	}

	// invalid id
	req3 := httptest.NewRequest(http.MethodGet, "/api/product/abc", nil)
	rr3 := httptest.NewRecorder()
	r.ServeHTTP(rr3, req3)
	if rr3.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected status 404 got %d", rr3.Code)
	}
}

func TestCreateOrder_SuccessAndValidationError(t *testing.T) {
	r := NewRouter()

	type tc struct {
		name      string
		productId string
		quantity  int
		reqStatus int
	}

	cases := []tc{
		{
			name:      "valid input",
			productId: "9",
			quantity:  10,
			reqStatus: http.StatusOK,
		},
		{
			name:      "product does not exist",
			productId: "999",
			quantity:  1,
			reqStatus: http.StatusBadRequest,
		},
		{
			name:      "invalid product id",
			productId: "aaa",
			quantity:  1,
			reqStatus: http.StatusUnprocessableEntity,
		},
		{
			name:      "invalid quantity",
			productId: "5",
			quantity:  0,
			reqStatus: http.StatusUnprocessableEntity,
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			payload := order.OrderRequest{
				CouponCode: "FIFTYOFF",
				Items: []order.OrderItem{
					{
						ProductID: testCase.productId,
						Quantity:  testCase.quantity,
					},
				},
			}
			b, _ := json.Marshal(payload)
			req := httptest.NewRequest(http.MethodPost, "/api/order/", bytes.NewReader(b))
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)

			if rr.Code != testCase.reqStatus {
				t.Fatalf("expected status %d got %d", testCase.reqStatus, rr.Code)
			}

			if testCase.name == "valid input" {
				var res order.OrderResponse
				if err := json.Unmarshal(rr.Body.Bytes(), &res); err != nil {
					t.Fatalf("Invalid json response from create order")
				}
				if !reflect.DeepEqual(res.Items, payload.Items) {
					t.Fatalf("response items does not match items in order request")
				}
			}
		})

	}
}
