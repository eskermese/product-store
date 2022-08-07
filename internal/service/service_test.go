package service_test

import (
	"testing"

	"github.com/ernur-eskermes/product-store/internal/service"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	productService, productStorage := mockProductService(t, nil)

	s := service.New(service.Deps{ProductStorage: productStorage})

	require.Equal(t, productService, s.Product)
}
