package repository

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/sir-shalahuddin/synapsis/model"
)

type cartRepository struct {
	db *sqlx.DB
}

func NewCartRepository(db *sqlx.DB) *cartRepository {
	return &cartRepository{db: db}
}

func (r *cartRepository) AddProduct(ctx context.Context, cart *model.Cart) (*model.Cart, error) {

	query := `insert into carts (product_id, user_id, quantity, created_at) values (:product_id, :user_id, :quantity, NOW())`

	_, err := r.db.NamedExecContext(ctx, query, cart)
	if err != nil {
		log.Println("cart Repository Error Add products : ", err)
		return nil, errors.New("failed to add prouct to carts from the database")
	}

	return cart, nil
}

func (r *cartRepository) GetProducts(ctx context.Context, userID string) ([]model.Product, error) {
	var products []model.Product

	query := `SELECT p.* FROM products p INNER JOIN carts c ON p.id = c.product_id WHERE c.user_id =$1`

	if err := r.db.SelectContext(ctx, &products, query, userID); err != nil {
		log.Println("cart Repository Error Get products : ", err)
		return nil, errors.New("failed to get product from carts from the database")
	}

	return products, nil
}

func (r *cartRepository) DeleteProduct(ctx context.Context, userID string, productID string) error {

	cart := model.Cart{
		UserID:    userID,
		ProductID: productID,
	}

	query := `delete from carts where user_id = :user_id and product_id = :product_id`

	_, err := r.db.NamedExecContext(ctx, query, cart)

	if err != nil {
		log.Println("cart Repository Error Delete products : ", err)
		return errors.New("failed to delete product from carts from the database")
	}
	return nil
}

func (r *cartRepository) ValidProduct(ctx context.Context, productID string) (bool, error) {
	var exists bool
	query := `SELECT exists (SELECT 1 FROM products WHERE id = $1)`

	err := r.db.QueryRowxContext(ctx, query, productID).Scan(&exists)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return exists, nil
}

func (r *cartRepository) DuplicateProduct(ctx context.Context, userID, productID string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS (SELECT 1 FROM carts WHERE user_id = $1 AND product_id = $2)`

	err := r.db.QueryRowxContext(ctx, query, userID, productID).Scan(&exists)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return exists, nil
}
