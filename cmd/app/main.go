package main

import (
	"context"
	"encoding/csv"
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/ernur-eskermes/product-store/pkg/logger"
	"github.com/gocarina/gocsv"

	"github.com/ernur-eskermes/product-store/internal/transport/grpc"
	grpcHandler "github.com/ernur-eskermes/product-store/internal/transport/grpc/handlers"

	"github.com/ernur-eskermes/product-store/pkg/database/mongodb"

	"github.com/ernur-eskermes/product-store/internal/config"
	"github.com/ernur-eskermes/product-store/internal/service"
	"github.com/ernur-eskermes/product-store/internal/storage"
)

func init() {
	gocsv.SetCSVReader(func(in io.Reader) gocsv.CSVReader {
		r := csv.NewReader(in)
		r.Comma = ';'

		return r
	})
}

func main() {
	log := logger.New("debug", "product_store")
	defer logger.Cleanup(log)

	cfg, err := config.New()
	if err != nil {
		log.Fatal("error when initializing the config", logger.Error(err))
	}

	mongoClient, err := mongodb.NewClient(cfg.Mongo.URI, cfg.Mongo.User, cfg.Mongo.Password)
	if err != nil {
		log.Fatal("connection to mongodb error", logger.Error(err))
	}

	db := mongoClient.Database(cfg.Mongo.Database)

	storages := storage.New(db)
	services := service.New(service.Deps{
		ProductStorage: storages.Product,
		HTTPClient:     http.DefaultClient,
	})

	grpcHandlers := grpcHandler.New(grpcHandler.Deps{
		ProductService: services.Product,
	})
	grpcSrv := grpc.New(grpc.Deps{
		Logger:         log,
		ProductHandler: grpcHandlers.Product,
	})

	go func() {
		log.Info("Starting gRPC server")

		if err = grpcSrv.ListenAndServe(cfg.GRPC.Port); err != nil {
			log.Error("gRPC ListenAndServer error", logger.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	log.Info("Shutting down server")

	grpcSrv.Stop()

	if err = mongoClient.Disconnect(context.Background()); err != nil {
		log.Error("disconnect mongodb error", logger.Error(err))
	}
}
