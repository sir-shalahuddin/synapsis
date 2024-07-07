package service

import (
	"context"
	"fmt"

	"github.com/sir-shalahuddin/synapsis/dto"
	"github.com/sir-shalahuddin/synapsis/model"
	error_helper "github.com/sir-shalahuddin/synapsis/pkg/helper"
)

type cartRepository interface {
	AddProduct(ctx context.Context, cart *model.Cart) (*model.Cart, error)
	GetProducts(ctx context.Context, userID string) ([]model.Product, error)
	DeleteProduct(ctx context.Context, userID string, productID string) error
	ValidProduct(ctx context.Context, productID string) (bool, error)
	DuplicateProduct(ctx context.Context, userID, productID string) (bool, error)
}

type cartService struct {
	repo cartRepository
}

func NewCartService(repo cartRepository) *cartService {
	return &cartService{repo: repo}
}

func (svc *cartService) AddProduct(ctx context.Context, req *dto.CartRequest, userID string) (*model.Cart, error) {

	cart := model.Cart{
		ProductID: req.ProductID,
		Quantity:  req.Quantity,
		UserID:    userID,
	}

	validProduct, err := svc.repo.ValidProduct(ctx, cart.ProductID)
	if err != nil {
		return nil, err
	}

	if !validProduct {
		return nil, error_helper.ErrorNotFound
	}

	duplicateProduct, err := svc.repo.DuplicateProduct(ctx, cart.UserID, cart.ProductID)
	if err != nil {
		return nil, err
	}

	fmt.Println(duplicateProduct)

	if duplicateProduct {
		return nil, error_helper.ErrorDuplicate
	}

	return svc.repo.AddProduct(ctx, &cart)

}

func (svc *cartService) GetProducts(ctx context.Context, userID string) ([]model.Product, error) {

	return svc.repo.GetProducts(ctx, userID)

}

func (svc *cartService) DeleteProduct(ctx context.Context, userID string, productID string) error {

	validProduct, err := svc.repo.ValidProduct(ctx, productID)

	if err != nil {
		return err
	}

	if !validProduct {
		return error_helper.ErrorNotFound
	}

	return svc.repo.DeleteProduct(ctx, userID, productID)
}
