package model

import "time"

type Product struct {
	ID        string     `db:"id"`
	Name      string     `db:"name"`
	Stock     uint       `db:"stock"`
	Price     float64    `db:"price"`
	Category  string     `db:"category"`
	CreatedAt *time.Time `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at"`
}
