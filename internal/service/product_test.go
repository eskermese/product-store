package service_test

import (
	"bytes"
	"context"
	"errors"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/ernur-eskermes/product-store/internal/core"
	"github.com/ernur-eskermes/product-store/internal/service"
	mock_service "github.com/ernur-eskermes/product-store/internal/service/mocks"
	"github.com/ernur-eskermes/product-store/pkg/filters"
	"github.com/gocarina/gocsv"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func mockProductService(t *testing.T, httpClient service.HTTPClient) (*service.ProductService, *mock_service.MockProductStorage) {
	t.Helper()

	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()

	productRepo := mock_service.NewMockProductStorage(mockCtl)

	productService := service.NewProductService(productRepo, httpClient)

	return productService, productRepo
}

func TestProduct_GetCSVProducts(t *testing.T) {
	type mockBehavior func(r *mock_service.MockHttpClient)

	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()

	httpClient := mock_service.NewMockHttpClient(mockCtl)
	productService, _ := mockProductService(t, httpClient)

	ctx := context.Background()

	products := []core.Product{
		{ID: primitive.ObjectID{}, Name: "Test Product", Price: 1000},
		{ID: primitive.ObjectID{}, Name: "Test Product2", Price: 2538},
	}
	b, err := gocsv.MarshalBytes(products)
	require.NoError(t, err)

	cases := []struct {
		name         string
		url          string
		expResp      []core.Product
		expErr       string
		mockBehavior mockBehavior
	}{
		{
			name:    "test_ok",
			url:     "https://some-url.com",
			expResp: products,
			mockBehavior: func(r *mock_service.MockHttpClient) {
				httpResp := ioutil.NopCloser(bytes.NewReader(b))
				r.EXPECT().Do(gomock.Any()).Return(&http.Response{Body: httpResp}, nil)
			},
		},
		{
			name:   "error_when_requesting",
			url:    "https://some-url.com",
			expErr: "error1",
			mockBehavior: func(r *mock_service.MockHttpClient) {
				r.EXPECT().Do(gomock.Any()).Return(nil, errors.New("error1"))
			},
		},
		{
			name:   "empty_byte_given",
			url:    "https://some-url.com",
			expErr: "empty csv file given",
			mockBehavior: func(r *mock_service.MockHttpClient) {
				httpResp := ioutil.NopCloser(bytes.NewReader([]byte("")))
				r.EXPECT().Do(gomock.Any()).Return(&http.Response{Body: httpResp}, nil)
			},
		},
		{
			name:         "empty_url",
			url:          "://some-url.com",
			expErr:       "parse \"://some-url.com\": missing protocol scheme",
			mockBehavior: func(r *mock_service.MockHttpClient) {},
		},
	}

	for _, s := range cases {
		t.Run(s.name, func(t *testing.T) {
			s.mockBehavior(httpClient)

			p, err := productService.GetCSVProducts(ctx, s.url)
			if err != nil {
				require.EqualError(t, err, s.expErr)
			} else {
				require.Equal(t, p, s.expResp)
			}
		})
	}
}

func TestProduct_GetAll(t *testing.T) {
	type mockBehavior func(r *mock_service.MockProductStorage)

	productService, productRepo := mockProductService(t, nil)

	ctx := context.Background()

	products := []core.Product{{ID: primitive.ObjectID{}, Name: "Test Product", Price: 1000}}

	cases := []struct {
		name         string
		expErr       string
		expResp      []core.Product
		mockBehavior mockBehavior
	}{
		{
			name:    "test_ok",
			expResp: products,
			mockBehavior: func(r *mock_service.MockProductStorage) {
				r.EXPECT().GetAll(gomock.Any(), gomock.Any()).Return(products, nil)
			},
		},
		{
			name:   "error_when_calling_GetAll",
			expErr: "error1",
			mockBehavior: func(r *mock_service.MockProductStorage) {
				r.EXPECT().GetAll(gomock.Any(), gomock.Any()).Return(nil, errors.New("error1"))
			},
		},
	}

	for _, s := range cases {
		t.Run(s.name, func(t *testing.T) {
			s.mockBehavior(productRepo)

			p, err := productService.GetAll(ctx, &filters.Filters{})
			if err != nil {
				require.EqualError(t, err, s.expErr)
			} else {
				require.Equal(t, p, s.expResp)
			}
		})
	}
}

func TestProduct_GetTotalRecords(t *testing.T) {
	type mockBehavior func(r *mock_service.MockProductStorage)

	productService, productRepo := mockProductService(t, nil)

	ctx := context.Background()

	cases := []struct {
		name         string
		expErr       string
		expResp      int64
		mockBehavior mockBehavior
	}{
		{
			name:    "test_ok",
			expResp: 12,
			mockBehavior: func(r *mock_service.MockProductStorage) {
				r.EXPECT().GetTotalRecords(ctx).Return(int64(12), nil)
			},
		},
		{
			name:   "error_when_calling_GetTotalRecords",
			expErr: "error1",
			mockBehavior: func(r *mock_service.MockProductStorage) {
				r.EXPECT().GetTotalRecords(gomock.Any()).Return(int64(0), errors.New("error1"))
			},
		},
	}

	for _, s := range cases {
		t.Run(s.name, func(t *testing.T) {
			s.mockBehavior(productRepo)

			totalRecords, err := productService.GetTotalRecords(ctx)
			if err != nil {
				require.EqualError(t, err, s.expErr)
			} else {
				require.Equal(t, totalRecords, s.expResp)
			}
		})
	}
}

func TestProduct_UpdateOrCreate(t *testing.T) {
	type mockBehavior func(r *mock_service.MockProductStorage)

	productService, productRepo := mockProductService(t, nil)

	ctx := context.Background()

	products := []core.Product{{ID: primitive.ObjectID{}, Name: "Test Product", Price: 1000}}

	cases := []struct {
		name         string
		expErr       string
		mockBehavior mockBehavior
	}{
		{
			name: "test_ok",
			mockBehavior: func(r *mock_service.MockProductStorage) {
				r.EXPECT().UpdateOrCreate(gomock.Any(), gomock.Any()).Return(nil)
			},
		},
		{
			name:   "error_when_calling_UpdateOrCreate",
			expErr: "error1",
			mockBehavior: func(r *mock_service.MockProductStorage) {
				r.EXPECT().UpdateOrCreate(gomock.Any(), gomock.Any()).Return(errors.New("error1"))
			},
		},
	}

	for _, s := range cases {
		t.Run(s.name, func(t *testing.T) {
			s.mockBehavior(productRepo)

			if err := productService.UpdateOrCreate(ctx, products); err != nil {
				require.EqualError(t, err, s.expErr)
			}
		})
	}
}
