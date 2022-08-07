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

test:
	go test -v -count=1 ./...

test100:
	go test -v -count=100 ./...

race:
	go test -v -race -count=1 ./...

.PHONY=cover
cover:
	go test -short -count=1 -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out
	rm coverage.out

gen:
	mockgen -source=internal/transport/grpc/handlers/product.go -destination=internal/transport/grpc/mocks/mock.go
	mockgen -source=internal/service/product.go -destination=internal/service/mocks/mock.go
