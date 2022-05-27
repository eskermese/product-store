package grpcHandler

import (
	"context"

	"github.com/ernur-eskermes/product-store/pkg/filters"

	"github.com/ernur-eskermes/product-store/internal/core"
)

type Handler struct {
	Product *Product
}

type ProductService interface {
	UpdateOrCreate(ctx context.Context, products []core.Product) error
	GetAll(ctx context.Context, f *filters.Filters) ([]core.Product, error)
	GetTotalRecords(ctx context.Context) (int64, error)
}

type Deps struct {
	ProductService ProductService
}

func New(deps Deps) *Handler {
	return &Handler{
		Product: NewProduct(deps.ProductService),
	}
}
