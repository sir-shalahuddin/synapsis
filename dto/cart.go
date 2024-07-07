package dto

type CartRequest struct {
	ProductID string `json:"product_id" validate:"required"`
	Quantity  uint   `json:"quantity" validate:"required"`
}
