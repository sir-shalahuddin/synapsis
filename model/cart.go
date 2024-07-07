package model

import "time"

type Cart struct {
	UserID    string     `db:"user_id"`
	ProductID string     `db:"product_id"`
	Quantity  uint       `db:"quantity"`
	CreatedAt *time.Time `db:"created_at"`
}
