package grpcHandler_test

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"testing"
	"time"

	"github.com/ernur-eskermes/product-store/internal/core"
	grpcHandler "github.com/ernur-eskermes/product-store/internal/transport/grpc/handlers"
	mock_grpcHandler "github.com/ernur-eskermes/product-store/internal/transport/grpc/mocks"
	pb "github.com/ernur-eskermes/product-store/pkg/domain"
	"github.com/ernur-eskermes/product-store/pkg/pagination"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
)

func dialer(productService *mock_grpcHandler.MockProductService) func(context.Context, string) (net.Conn, error) {
	listener := bufconn.Listen(1024 * 1024)

	server := grpc.NewServer()

	pb.RegisterProductServiceServer(server, grpcHandler.NewProductHandler(productService))

	go func() {
		if err := server.Serve(listener); err != nil {
			log.Fatal(err)
		}
	}()

	return func(context.Context, string) (net.Conn, error) {
		return listener.Dial()
	}
}

func TestProductHandler_Fetch(t *testing.T) {
	type mockBehavior func(r *mock_grpcHandler.MockProductService)

	cases := []struct {
		name         string
		url          string
		errCode      codes.Code
		errMsg       string
		mockBehavior mockBehavior
	}{
		{
			name:    "valid_request",
			url:     "https://some-url.com",
			errCode: codes.OK,
			errMsg:  "",

			mockBehavior: func(r *mock_grpcHandler.MockProductService) {
				r.EXPECT().GetCSVProducts(gomock.Any(), gomock.Any()).Return([]core.Product{}, nil)
				r.EXPECT().UpdateOrCreate(gomock.Any(), gomock.Any()).Return(nil)
			},
		},
		{
			name:    "invalid_url_for_request",
			url:     "some-url.com",
			errCode: codes.InvalidArgument,
			errMsg:  fmt.Sprintf("parse \"%s\": invalid URI for request", "some-url.com"),

			mockBehavior: func(r *mock_grpcHandler.MockProductService) {},
		},
		{
			name:    "error_when_calling_GetCSVProducts_method",
			url:     "https://some-url.com",
			errCode: codes.InvalidArgument,
			errMsg:  "error",

			mockBehavior: func(r *mock_grpcHandler.MockProductService) {
				r.EXPECT().GetCSVProducts(gomock.Any(), gomock.Any()).Return(nil, errors.New("error"))
			},
		},
		{
			name:    "error_when_calling_UpdateOrCreate_method",
			url:     "https://some-url.com",
			errCode: codes.Unknown,
			errMsg:  "error create or update",

			mockBehavior: func(r *mock_grpcHandler.MockProductService) {
				r.EXPECT().GetCSVProducts(gomock.Any(), gomock.Any()).Return([]core.Product{}, nil)
				r.EXPECT().UpdateOrCreate(gomock.Any(), gomock.Any()).Return(errors.New("error create or update"))
			},
		},
	}

	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()

	ctx := context.Background()

	productService := mock_grpcHandler.NewMockProductService(mockCtl)

	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithContextDialer(dialer(productService)))
	require.NoError(t, err)
	defer conn.Close()

	client := pb.NewProductServiceClient(conn)

	for _, s := range cases {
		t.Run(s.name, func(t *testing.T) {
			s.mockBehavior(productService)

			if _, err = client.Fetch(ctx, &pb.FetchRequest{Url: s.url}); err != nil {
				if er, ok := status.FromError(err); ok {
					require.Equal(t, er.Code(), s.errCode)
					require.Equal(t, er.Message(), s.errMsg)
				}
			}
		})
	}
}

func TestProductHandler_List(t *testing.T) {
	type mockBehavior func(r *mock_grpcHandler.MockProductService)

	products := []core.Product{
		{ID: primitive.ObjectID{}, Name: "product 1", Price: 12},
		{ID: primitive.ObjectID{}, Name: "product 2", Price: 121},
		{ID: primitive.ObjectID{}, Name: "product 3", Price: 122},
	}
	metadata, _ := pagination.New(12, 1, 3)

	cases := []struct {
		name        string
		expResp     []core.Product
		expMetadata *pagination.Pagination

		expErr map[string]string
		errMsg string

		code codes.Code

		body         *pb.Filters
		mockBehavior mockBehavior
	}{
		{
			name:        "test_ok",
			expResp:     products,
			expMetadata: metadata,
			code:        codes.OK,

			body: &pb.Filters{Page: 1, PageSize: 3, Sort: "name"},
			mockBehavior: func(r *mock_grpcHandler.MockProductService) {
				r.EXPECT().GetAll(gomock.Any(), gomock.Any()).Return(products, nil)
				r.EXPECT().GetTotalRecords(gomock.Any()).Return(int64(12), nil)
			},
		},
		{
			name:   "error_when_calling_GetAll_method",
			code:   codes.Unknown,
			errMsg: "error when getAll",

			body: &pb.Filters{Page: 1, PageSize: 3, Sort: "name"},
			mockBehavior: func(r *mock_grpcHandler.MockProductService) {
				r.EXPECT().GetAll(gomock.Any(), gomock.Any()).Return(nil, errors.New("error when getAll"))
			},
		},
		{
			name:   "error_when_calling_GetTotalRecords_method",
			code:   codes.Unknown,
			errMsg: "error when getTotalRecords",

			body: &pb.Filters{Page: 1, PageSize: 3, Sort: "name"},
			mockBehavior: func(r *mock_grpcHandler.MockProductService) {
				r.EXPECT().GetAll(gomock.Any(), gomock.Any()).Return(products, nil)
				r.EXPECT().GetTotalRecords(gomock.Any()).Return(int64(0), errors.New("error when getTotalRecords"))
			},
		},
		{
			name:   "invalid_sort_name",
			code:   codes.InvalidArgument,
			expErr: map[string]string{"sort": "invalid sort value", "page": "must be greater than zero", "page_size": "must be greater than zero"},
			errMsg: "invalid filter params",

			body:         &pb.Filters{Page: -1, PageSize: -3, Sort: "namea"},
			mockBehavior: func(r *mock_grpcHandler.MockProductService) {},
		},
		{
			name:   "invalid_page_params",
			code:   codes.InvalidArgument,
			expErr: map[string]string{"page": "must be a maximum of 10 million", "page_size": "must be a maximum of 100"},
			errMsg: "invalid filter params",

			body:         &pb.Filters{Page: 10_000_000, PageSize: 300, Sort: "price"},
			mockBehavior: func(r *mock_grpcHandler.MockProductService) {},
		},
	}

	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()

	ctx := context.Background()

	productService := mock_grpcHandler.NewMockProductService(mockCtl)

	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithContextDialer(dialer(productService)))
	require.NoError(t, err)
	defer conn.Close()

	productClient := pb.NewProductServiceClient(conn)

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	stream, err := productClient.List(ctx)
	if err != nil {
		log.Fatalf("%v.Execute(ctx) = %v, %v: ", productClient, stream, err)
	}

	completed := make(chan struct{})

	go func() {
		completedCases := make(map[*pb.ListResponse]error, len(cases))
		i := 0

		for {
			// receive the second argument(161 line) many times
			result, err := stream.Recv()
			if errors.Is(err, io.EOF) {
				close(completed)

				return
			}

			if v, ok := completedCases[result]; ok && errors.Is(v, err) {
				fmt.Println("already completed case is requested again:", err)
				continue
			} else {
				completedCases[result] = err
			}

			want := cases[i]

			if err != nil {
				if er, ok := status.FromError(err); ok {
					for _, detail := range er.Details() {
						if badReq, ok := detail.(*errdetails.BadRequest); ok {
							require.EqualValues(t, want.expErr, badRequestToMap(badReq))
						}
					}

					require.Equal(t, want.code, er.Code())
					require.Equal(t, want.errMsg, er.Message())
				}
			} else {
				require.Equal(t, want.expResp, PBProductToStruct(result.Results))
				require.Equal(t, want.expMetadata, PBMetadataToStruct(result.Metadata))
			}

			completed <- struct{}{}
			i++
		}
	}()

	for _, c := range cases {
		c.mockBehavior(productService)

		// TODO don't work. For some reason, it processes the second case(161 line) many times.
		if err = stream.Send(c.body); err != nil && !errors.Is(err, io.EOF) {
			log.Fatalf("%v.Send(%v) = %v: ", stream, c.body, err)
		}

		<-completed
	}

	if err = stream.CloseSend(); err != nil {
		log.Fatalf("%v.CloseSend() got error %v, want %v", stream, err, nil)
	}
}

func PBProductToStruct(products []*pb.ListResponse_Product) []core.Product {
	res := make([]core.Product, 0, len(products))

	for _, product := range products {
		res = append(res, core.Product{
			Name:  product.GetName(),
			Price: int(product.GetPrice()),
		})
	}

	return res
}

func PBMetadataToStruct(metadata *pb.ListResponse_MetaData) *pagination.Pagination {
	return &pagination.Pagination{
		CurrentPage:  metadata.CurrentPage,
		PageSize:     metadata.PageSize,
		FirstPage:    metadata.FirstPage,
		LastPage:     metadata.LastPage,
		TotalRecords: metadata.TotalRecords,
	}
}

func badRequestToMap(err *errdetails.BadRequest) map[string]string {
	e := make(map[string]string, len(err.FieldViolations))

	for _, i := range err.FieldViolations {
		e[i.GetField()] = i.GetDescription()
	}

	return e
}
