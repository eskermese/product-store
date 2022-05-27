package main

import (
	"context"
	"encoding/csv"
	"io"
	"os"
	"os/signal"
	"syscall"

	"github.com/gocarina/gocsv"
	log "github.com/sirupsen/logrus"

	"github.com/ernur-eskermes/product-store/internal/transport/grpc"
	grpcHandler "github.com/ernur-eskermes/product-store/internal/transport/grpc/handlers"

	"github.com/ernur-eskermes/product-store/pkg/database/mongodb"

	"github.com/ernur-eskermes/product-store/internal/config"
	"github.com/ernur-eskermes/product-store/internal/service"
	"github.com/ernur-eskermes/product-store/internal/storage"
)

func init() {
	log.SetReportCaller(true)
	log.SetFormatter(&log.JSONFormatter{})

	gocsv.SetCSVReader(func(in io.Reader) gocsv.CSVReader {
		r := csv.NewReader(in)
		r.Comma = ';'

		return r
	})
}

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatal(err)
	}

	mongoClient, err := mongodb.NewClient(cfg.Mongo.URI, cfg.Mongo.User, cfg.Mongo.Password)
	if err != nil {
		log.Fatal(err)
	}

	db := mongoClient.Database(cfg.Mongo.Database)

	storages := storage.New(db)
	services := service.New(service.Deps{
		ProductStorage: storages.Product,
	})

	grpcHandlers := grpcHandler.New(grpcHandler.Deps{
		ProductService: services.Product,
	})
	grpcSrv := grpc.New(grpcHandlers)

	go func() {
		log.Info("Starting gRPC server")

		if err = grpcSrv.ListenAndServe(cfg.GRPC.Port); err != nil {
			log.Error("gRPC ListenAndServer error", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	log.Info("Shutting down server")

	grpcSrv.Stop()

	if err = mongoClient.Disconnect(context.TODO()); err != nil {
		log.Error(err)
	}
}
