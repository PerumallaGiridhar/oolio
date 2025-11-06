package order

import (
	"net/http"
	"strconv"

	"github.com/PerumallaGiridhar/oolio/internal/binding"
	"github.com/PerumallaGiridhar/oolio/internal/data"
	"github.com/PerumallaGiridhar/oolio/internal/response"
	"github.com/google/uuid"
)

func CreateOrderRequest(w http.ResponseWriter, r *http.Request) {
	var req OrderRequest
	if err := binding.BindAndValidateJSONRequest(r, &req); err != nil {
		response.JSONValidationErrorResponse(w, err)
		return
	}

	var products []data.Product
	for _, item := range req.Items {
		_, err := strconv.Atoi(item.ProductID)
		if err != nil {
			errorMsg := map[string]string{"error": "invalid product Id, Id must be an integer"}
			response.JSONValidationErrorResponse(w, errorMsg)
			return
		}
		product, found := data.GetProductByID(item.ProductID)
		if found {
			products = append(products, product)
		} else {
			response.JSONErrorResponse(w, http.StatusBadRequest, "ProductId does not exists")
			return
		}
	}

	respData := OrderResponse{
		ID:         uuid.New().String(),
		CouponCode: req.CouponCode,
		Items:      req.Items,
		Products:   products,
	}
	response.JSONResponse(w, http.StatusOK, respData)
}
