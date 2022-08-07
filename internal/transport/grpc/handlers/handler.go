package grpcHandler

type Handler struct {
	Product *ProductHandler
}

type Deps struct {
	ProductService ProductService
}

func New(deps Deps) *Handler {
	return &Handler{
		Product: NewProductHandler(deps.ProductService),
	}
}
