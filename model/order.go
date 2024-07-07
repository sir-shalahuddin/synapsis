package model

import "time"

type Order struct {
	ID         string     `db:"id"`
	UserID     string     `db:"user_id"`
	TotalPrice float64    `db:"total_price"`
	Status     string     `db:"status"`
	CreatedAt  *time.Time `db:"created_at"`
	UpdatedAt  *time.Time `db:"updated_at"`
}

type OrderProduct struct {
	OrderID      string     `db:"order_id"`
	ProductID    string     `db:"product_id"`
	Quantity     uint       `db:"quantity"`
	ProductName  string     `db:"product_name"`
	ProductPrice float64    `db:"product_price"`
	CreatedAt    *time.Time `db:"created_at"`
}
