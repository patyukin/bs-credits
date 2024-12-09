.PHONY: start stop rebuild gen up down restart deps tidy gen-auth-api vendor-proto

LOCAL_BIN:=$(CURDIR)/bin

up:
	docker compose up -d

down:
	docker compose down

start:
	docker compose start

stop:
	docker compose stop

restart:
	docker compose restart

rebuild:
	docker compose down -v --remove-orphans
	docker compose up -d --build

install-deps:
	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.34.2
	GOBIN=$(LOCAL_BIN) go install -mod=mod google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.5.1
	GOBIN=$(LOCAL_BIN) go install github.com/envoyproxy/protoc-gen-validate@v1.1.0
	GOBIN=$(LOCAL_BIN) go install github.com/pressly/goose/v3/cmd/goose@latest

tidy:
	GOBIN=$(LOCAL_BIN) go mod tidy

test:
	go test ./...

cov:
	GOBIN=$(LOCAL_BIN) go test ./... -coverprofile=coverage.out
	GOBIN=$(LOCAL_BIN) go tool cover -html=coverage.out -o coverage.html