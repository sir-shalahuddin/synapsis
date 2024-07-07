package service

import (
	"context"

	"github.com/sir-shalahuddin/synapsis/model"
)

type productRepository interface {
	GetProducts(ctx context.Context, category string) ([]model.Product, error)
}

type productService struct {
	repo productRepository
}

func NewproductService(repo productRepository) *productService {
	return &productService{repo: repo}
}

func (svc *productService) GetProducts(ctx context.Context, category string) ([]model.Product, error) {
	return svc.repo.GetProducts(ctx, category)
}
