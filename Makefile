.PHONY:
.SILENT:
.DEFAULT_GOAL := run

build:
	go mod download && CGO_ENABLED=0 GOOS=linux go build -o ./.bin/app ./cmd/app/main.go

run: build
	docker-compose up --remove-orphans

deploy: build
	docker-compose -f deploy/docker-compose.yml up --remove-orphans

down:
	docker-compose down -v

swag:
	swag init -g cmd/app/main.go

lint:
	golangci-lint run

proto:
	protoc --go_out=proto --go-grpc_out=proto proto/product.proto