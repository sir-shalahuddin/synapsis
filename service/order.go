package service

import (
	"context"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/sir-shalahuddin/synapsis/dto"
	"github.com/sir-shalahuddin/synapsis/model"
	error_helper "github.com/sir-shalahuddin/synapsis/pkg/helper"
)

type orderRepository interface {
	GetStock(ctx context.Context, tx *sqlx.Tx, productIDs []string) ([]model.Product, error)
	UpdateStock(ctx context.Context, tx *sqlx.Tx, products []model.OrderProduct) error
	CreateOrder(ctx context.Context, tx *sqlx.Tx, userID string, totalPrice float64) (string, error)
	InsertOrderProduct(ctx context.Context, tx *sqlx.Tx, orderID string, products []model.OrderProduct) error
	PayOrder(ctx context.Context, orderID string) (bool, error)
	ValidUser(ctx context.Context, userID string, orderID string) (bool, error)
}

type txRepository interface {
	BeginTx(ctx context.Context) (*sqlx.Tx, error)
	Commit(tx *sqlx.Tx) error
	Rollback(tx *sqlx.Tx) error
}

type orderService struct {
	repo   orderRepository
	txRepo txRepository
}

func NewOrderService(repo orderRepository, txRepo txRepository) *orderService {
	return &orderService{repo: repo, txRepo: txRepo}
}

func (svc *orderService) CreateOrder(ctx context.Context, req *dto.OrderRequest, userID string) (string, error) {

	tx, err := svc.txRepo.BeginTx(ctx)
	if err != nil {
		return "", err
	}

	defer func() {
		if err != nil {
			if errTx := svc.txRepo.Rollback(tx); errTx != nil {
				log.Println("Error rolback : ", errTx)
			}
			return
		}
		if errTx := svc.txRepo.Commit(tx); errTx != nil {
			log.Println("Error Commit : ", errTx)
		}
	}()

	productIDs := make([]string, len(req.Orders))
	for i, order := range req.Orders {
		productIDs[i] = order.ProductID
	}

	stocks, err := svc.repo.GetStock(ctx, tx, productIDs)
	if err != nil {
		return "", err
	}

	productPriceMap := make(map[string]float64)
	for _, p := range stocks {
		productPriceMap[p.ID] = p.Price
	}

	totalPrice := 0.0
	for _, item := range req.Orders {
		price, exists := productPriceMap[item.ProductID]
		if !exists {
			return "", fmt.Errorf("product with ID %s not found", item.ProductID)
		}
		totalPrice += price * float64(item.Quantity)
	}

	var products []model.OrderProduct

	for _, v := range stocks {
		products = append(products, model.OrderProduct{
			ProductID:    v.ID,
			Quantity:     v.Stock,
			ProductName:  v.Name,
			ProductPrice: v.Price,
		})
	}

	err = svc.repo.UpdateStock(ctx, tx, products)
	if err != nil {
		return "", err
	}

	orderID, err := svc.repo.CreateOrder(ctx, tx, userID, totalPrice)
	if err != nil {
		return "", err
	}

	err = svc.repo.InsertOrderProduct(ctx, tx, orderID, products)
	if err != nil {
		return "", err
	}

	return orderID, nil
}

func (svc *orderService) PayOrder(ctx context.Context, userID string, orderID string) (bool, error) {

	validUser, err := svc.repo.ValidUser(ctx, userID, orderID)

	if err != nil {
		return false, err
	}

	if !validUser {
		return false , error_helper.ErrorUnAuthorized
	}

	return svc.repo.PayOrder(ctx, orderID)
}
