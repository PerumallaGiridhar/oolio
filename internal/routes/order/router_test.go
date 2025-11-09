package order

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	enlocales "github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	v10 "github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"

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
			payload := OrderRequest{
				CouponCode: "FIFTYOFF",
				Items: []OrderItem{
					{
						ProductID: testCase.productId,
						Quantity:  testCase.quantity,
					},
				},
			}
			b, _ := json.Marshal(payload)
			req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(b))
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)

			if rr.Code != testCase.reqStatus {
				t.Fatalf("expected status %d got %d", testCase.reqStatus, rr.Code)
			}

			if testCase.name == "valid input" {
				var res OrderResponse
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
