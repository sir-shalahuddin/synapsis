package dto

type OrderRequest struct {
	Orders []Order `json:"orders" validate:"required,dive"`
}

type Order struct {
	ProductID string `json:"product_id" validate:"required,uuid4"`
	Quantity  uint   `json:"quantity" validate:"required,gt=0"`
}
