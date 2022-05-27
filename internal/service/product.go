package service

import (
	"context"

	"github.com/ernur-eskermes/product-store/internal/core"
	"github.com/ernur-eskermes/product-store/pkg/filters"
)

type Product struct {
	repo ProductStorage
}

func NewProduct(repo ProductStorage) *Product {
	return &Product{
		repo: repo,
	}
}

func (s *Product) GetAll(ctx context.Context, f *filters.Filters) ([]core.Product, error) {
	return s.repo.GetAll(ctx, f)
}

func (s *Product) GetTotalRecords(ctx context.Context) (int64, error) {
	return s.repo.GetTotalRecords(ctx)
}

func (s *Product) UpdateOrCreate(ctx context.Context, products []core.Product) error {
	return s.repo.UpdateOrCreate(ctx, products)
}
