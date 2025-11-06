package product

import (
	"net/http"
	"strconv"

	"github.com/PerumallaGiridhar/oolio/internal/data"
	"github.com/PerumallaGiridhar/oolio/internal/response"
	"github.com/go-chi/chi/v5"
)

func ListProducts(w http.ResponseWriter, r *http.Request) {
	products := data.GetAllProducts()
	response.JSONResponse(w, http.StatusOK, products)
}

func FindProductById(w http.ResponseWriter, r *http.Request) {
	productIdParam := chi.URLParam(r, "productId")
	_, err := strconv.Atoi(productIdParam)
	if err != nil {
		errorMsg := map[string]string{"error": "invalid product Id, Id must be an integer"}
		response.JSONValidationErrorResponse(w, errorMsg)
		return
	}

	product, found := data.GetProductByID(productIdParam)
	if !found {
		response.JSONErrorResponse(w, http.StatusNotFound, "product not found")
		return
	}

	response.JSONResponse(w, http.StatusOK, product)

}
