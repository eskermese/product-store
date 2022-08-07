package grpcHandler

import (
	"context"
	"errors"
	"io"
	"net/url"

	"github.com/ernur-eskermes/product-store/pkg/pagination"

	"github.com/ernur-eskermes/product-store/internal/core"
	pb "github.com/ernur-eskermes/product-store/pkg/domain"
	"github.com/ernur-eskermes/product-store/pkg/filters"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ProductService interface {
	UpdateOrCreate(ctx context.Context, products []core.Product) error
	GetAll(ctx context.Context, f *filters.Filters) ([]core.Product, error)
	GetTotalRecords(ctx context.Context) (int64, error)
	GetCSVProducts(ctx context.Context, url string) ([]core.Product, error)
}

type ProductHandler struct {
	service ProductService
	pb.UnimplementedProductServiceServer
}

func NewProductHandler(s ProductService) *ProductHandler {
	return &ProductHandler{service: s}
}

func (h *ProductHandler) Fetch(ctx context.Context, req *pb.FetchRequest) (*empty.Empty, error) {
	if _, err := url.ParseRequestURI(req.GetUrl()); err != nil {
		return &empty.Empty{}, status.Error(codes.InvalidArgument, err.Error())
	}

	products, err := h.service.GetCSVProducts(ctx, req.GetUrl())
	if err != nil {
		return &empty.Empty{}, status.Error(codes.InvalidArgument, err.Error())
	}

	if err = h.service.UpdateOrCreate(ctx, products); err != nil {
		return &empty.Empty{}, status.Error(codes.Unknown, err.Error())
	}

	return &empty.Empty{}, nil
}

func (h *ProductHandler) List(stream pb.ProductService_ListServer) error {
	for {
		req, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			return nil
		}

		if err != nil {
			return err
		}

		f := filters.New(
			req.Page,
			req.PageSize,
			req.Sort,
			core.ProductDefaultSort,
			core.ProductSortSafeList,
		)
		if err = filters.ValidateFilters(f); err != nil {
			return ErrorFilterResponse(err)
		}

		products, err := h.service.GetAll(context.TODO(), f)
		if err != nil {
			return status.Error(codes.Unknown, err.Error())
		}

		totalRecords, err := h.service.GetTotalRecords(context.TODO())
		if err != nil {
			return status.Error(codes.Unknown, err.Error())
		}

		res := make([]*pb.ListResponse_Product, 0, len(products))

		for _, product := range products {
			res = append(res, &pb.ListResponse_Product{Name: product.Name, Price: int64(product.Price)})
		}

		if err = stream.Send(&pb.ListResponse{
			Results:  res,
			Metadata: calculateMetadata(totalRecords, f.Page, f.PageSize),
		}); err != nil {
			return err
		}
	}
}

func calculateMetadata(totalRecords, page, pageSize int64) *pb.ListResponse_MetaData {
	p, err := pagination.New(totalRecords, page, pageSize)
	if err != nil {
		return nil
	}

	return &pb.ListResponse_MetaData{
		CurrentPage:  p.CurrentPage,
		PageSize:     p.PageSize,
		FirstPage:    p.FirstPage,
		LastPage:     p.LastPage,
		TotalRecords: p.TotalRecords,
	}
}
