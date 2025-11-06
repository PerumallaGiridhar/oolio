package order

import "github.com/PerumallaGiridhar/oolio/internal/data"

type OrderItem struct {
	ProductID string `json:"productId" validate:"required"`
	Quantity  int    `json:"quantity" validate:"required,min=1"`
}

type OrderRequest struct {
	CouponCode string      `json:"couponCode" validate:"omitempty,promocode"`
	Items      []OrderItem `json:"items" validate:"required,dive,required"`
}

type OrderResponse struct {
	ID         string         `json:"id"`
	CouponCode string         `json:"couponCode"`
	Items      []OrderItem    `json:"items"`
	Products   []data.Product `json:"products"`
}
