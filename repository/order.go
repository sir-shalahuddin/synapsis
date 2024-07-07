package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/sir-shalahuddin/synapsis/model"
)

type orderRepository struct {
	db *sqlx.DB
}

func NewOrderRepository(db *sqlx.DB) *orderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) GetStock(ctx context.Context, tx *sqlx.Tx, productIDs []string) ([]model.Product, error) {
	query := `SELECT id, stock, price, name FROM products WHERE id IN (?) FOR UPDATE`
	query, args, err := sqlx.In(query, productIDs)
	if err != nil {
		log.Println("Error building query in GetStock:", err)
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	query = tx.Rebind(query)

	var products []model.Product

	log.Println("GetStock Query:", query)
	log.Println("GetStock Arguments:", args)

	err = tx.SelectContext(ctx, &products, query, args...)
	if err != nil {
		log.Println("Error executing query in GetStock:", err)
		return nil, fmt.Errorf("failed to get stocks and prices: %w", err)
	}

	return products, nil
}

func (r *orderRepository) UpdateStock(ctx context.Context, tx *sqlx.Tx, products []model.OrderProduct) error {
	if len(products) == 0 {
		return nil
	}

	query := `UPDATE products SET stock = stock - CASE `
	args := []interface{}{}
	i := 1

	for _, product := range products {
		productID := product.ProductID
		quantity := product.Quantity

		query += fmt.Sprintf("WHEN id = $%d THEN $%d::INTEGER ", i, i+1)
		args = append(args, productID, quantity)
		i += 2
	}

	query += `END WHERE id IN (`
	for j := 1; j <= len(products); j++ {
		if j > 1 {
			query += ", "
		}
		query += fmt.Sprintf("$%d", j*2-1)
	}
	query += `)`

	log.Println("UpdateStock Query:", query)
	log.Println("UpdateStock Arguments:", args)

	_, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		log.Println("Error executing query in UpdateStock:", err)
		return fmt.Errorf("failed to update stocks: %w", err)
	}
	return nil
}

func (r *orderRepository) CreateOrder(ctx context.Context, tx *sqlx.Tx, userID string, totalPrice float64) (string, error) {
	var id string

	query := `INSERT INTO orders (user_id, total_price, status, created_at) VALUES ($1, $2, 'menunggu pembayaran', NOW()) RETURNING id`

	err := tx.QueryRowContext(ctx, query, userID, totalPrice).Scan(&id)
	if err != nil {
		log.Println("Error executing query in CreateOrder:", err)
		return "", fmt.Errorf("failed to create order: %w", err)
	}

	return id, nil
}

func (r *orderRepository) InsertOrderProduct(ctx context.Context, tx *sqlx.Tx, orderID string, products []model.OrderProduct) error {
	query := `INSERT INTO orders_products 
						(order_id, product_id, quantity, product_name, product_price, created_at) VALUES 
						(:order_id, :product_id, :quantity, :product_name, :product_price, NOW())`

	_, err := tx.NamedExecContext(ctx, query, products)
	if err != nil {
		log.Println("Error executing query in InsertOrderProduct:", err)
		return fmt.Errorf("failed to insert order products: %w", err)
	}

	return nil
}

func (r *orderRepository) PayOrder(ctx context.Context, orderID string) (bool, error) {
	query := `UPDATE orders SET status = 'pembayaran diterima' where id = $1`

	_, err := r.db.ExecContext(ctx, query, orderID)

	if err != nil {
		log.Println("Error executing query in Pay Order:", err)
		return false, fmt.Errorf("failed to pay order: %w", err)
	}

	return true, nil
}

func (r *orderRepository) ValidUser(ctx context.Context, userID string, orderID string) (bool, error) {
	var exists bool

	query := `select exists (select 1 from orders where id = $1 and user_id=$2)`

	err := r.db.QueryRowContext(ctx, query, orderID, userID).Scan(&exists)

	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		log.Println("Error executing query in ValidUser():", err)
		return false, fmt.Errorf("failed to check valid user: %w", err)
	}

	return exists, nil
}
