export PATH := $(shell go env GOPATH)/bin:$(PATH)

.DEFAULT_GOAL := build

.PHONY: fmt vet build
fmt:
	@go fmt ./...

vet: fmt
	@go vet ./...

build: vet
	@go build ./...

clean:
	@go mod tidy
	@go clean

test:
	@go test ./...

unit:
	@go test ./internal/...

coverage:
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html

integration:
	@go test -v ./integration_test.go
