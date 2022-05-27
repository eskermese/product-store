package grpc

import (
	"fmt"
	"net"

	grpcHandler "github.com/ernur-eskermes/product-store/internal/transport/grpc/handlers"

	"google.golang.org/grpc/reflection"

	pb "github.com/ernur-eskermes/product-store/pkg/domain"
	"google.golang.org/grpc"
)

type Server struct {
	grpcSrv        *grpc.Server
	productHandler *grpcHandler.Product
}

func New(handlers *grpcHandler.Handler) *Server {
	return &Server{
		grpcSrv:        grpc.NewServer(),
		productHandler: handlers.Product,
	}
}

func (s *Server) ListenAndServe(port int) error {
	addr := fmt.Sprintf(":%d", port)

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	pb.RegisterProductServiceServer(s.grpcSrv, s.productHandler)
	reflection.Register(s.grpcSrv)

	return s.grpcSrv.Serve(lis)
}

func (s *Server) Stop() {
	s.grpcSrv.GracefulStop()
}
