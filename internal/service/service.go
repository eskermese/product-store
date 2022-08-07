package service

type Service struct {
	Product *ProductService
}

type Deps struct {
	ProductStorage ProductStorage

	HTTPClient HTTPClient
}

func New(deps Deps) *Service {
	return &Service{
		Product: NewProductService(deps.ProductStorage, deps.HTTPClient),
	}
}
