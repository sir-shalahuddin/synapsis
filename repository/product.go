package repository

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/sir-shalahuddin/synapsis/model"
)

type productRepository struct {
	db *sqlx.DB
}

func NewProductRepository(db *sqlx.DB) *productRepository {
	return &productRepository{db: db}
}

func (r *productRepository) GetProducts(ctx context.Context, category string) ([]model.Product, error) {
	var products []model.Product
	fmt.Println(category)
	var query string
	var args []interface{}

	if category != "" {
		query = `SELECT * FROM products WHERE category = $1`
		args = append(args, category)
	} else {
		query = `SELECT * FROM products`
	}

	err := r.db.SelectContext(ctx, &products, query, args...)
	if err != nil {
		log.Println("product Repository Error Get products : ", err)
		return nil, errors.New("failed to get products from the database")
	}

	return products, nil
}

