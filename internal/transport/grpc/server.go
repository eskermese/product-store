package grpc

import (
	"fmt"
	"net"
	"time"

	"github.com/ernur-eskermes/product-store/pkg/logger"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc/reflection"

	pb "github.com/ernur-eskermes/product-store/pkg/domain"
	"google.golang.org/grpc"
)

type Deps struct {
	Logger logger.Logger

	ProductHandler pb.ProductServiceServer
}

type Server struct {
	Deps
	grpcSrv *grpc.Server
}

func New(deps Deps) *Server {
	zapLogger := logger.GetZapLogger(deps.Logger)

	opts := []grpc_zap.Option{
		grpc_zap.WithDurationField(func(duration time.Duration) zapcore.Field {
			return zap.Int64("grpc.time_ns", duration.Nanoseconds())
		}),
	}

	return &Server{
		grpcSrv: grpc.NewServer(
			grpc_middleware.WithUnaryServerChain(
				grpc_ctxtags.UnaryServerInterceptor(),
				grpc_zap.UnaryServerInterceptor(zapLogger, opts...),
			),
			grpc_middleware.WithStreamServerChain(
				grpc_ctxtags.StreamServerInterceptor(),
				grpc_zap.StreamServerInterceptor(zapLogger, opts...),
			),
		),
		Deps: deps,
	}
}

func (s *Server) ListenAndServe(port int) error {
	addr := fmt.Sprintf(":%d", port)

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	pb.RegisterProductServiceServer(s.grpcSrv, s.Deps.ProductHandler)
	reflection.Register(s.grpcSrv)

	return s.grpcSrv.Serve(lis)
}

func (s *Server) Stop() {
	s.grpcSrv.GracefulStop()
}
