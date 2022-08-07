package service

import (
	"context"
	"net/http"

	"github.com/ernur-eskermes/product-store/internal/core"
	"github.com/ernur-eskermes/product-store/pkg/filters"
	"github.com/gocarina/gocsv"
)

type ProductStorage interface {
	GetAll(ctx context.Context, f *filters.Filters) ([]core.Product, error)
	UpdateOrCreate(ctx context.Context, products []core.Product) error
	GetTotalRecords(ctx context.Context) (int64, error)
}

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type ProductService struct {
	repo ProductStorage

	httpClient HTTPClient
}

func NewProductService(repo ProductStorage, httpClient HTTPClient) *ProductService {
	return &ProductService{
		repo: repo,

		httpClient: httpClient,
	}
}

func (s *ProductService) GetAll(ctx context.Context, f *filters.Filters) ([]core.Product, error) {
	return s.repo.GetAll(ctx, f)
}

func (s *ProductService) GetTotalRecords(ctx context.Context) (int64, error) {
	return s.repo.GetTotalRecords(ctx)
}

func (s *ProductService) UpdateOrCreate(ctx context.Context, products []core.Product) error {
	return s.repo.UpdateOrCreate(ctx, products)
}

func (s *ProductService) GetCSVProducts(ctx context.Context, url string) ([]core.Product, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	products := make([]core.Product, 0)
	if err = gocsv.Unmarshal(resp.Body, &products); err != nil {
		return nil, err
	}

	return products, nil
}
