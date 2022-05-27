package service

import (
	"context"

	"github.com/ernur-eskermes/product-store/pkg/filters"

	"github.com/ernur-eskermes/product-store/internal/core"
)

type Service struct {
	Product *Product
}

type ProductStorage interface {
	GetAll(ctx context.Context, f *filters.Filters) ([]core.Product, error)
	UpdateOrCreate(ctx context.Context, products []core.Product) error
	GetTotalRecords(ctx context.Context) (int64, error)
}

type Deps struct {
	ProductStorage ProductStorage
}

func New(deps Deps) *Service {
	return &Service{
		Product: NewProduct(deps.ProductStorage),
	}
}
