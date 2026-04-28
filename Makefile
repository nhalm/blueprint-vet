.PHONY: build test install fmt lint tidy

build:
	go build ./...

test:
	go test ./...

install:
	go install ./cmd/blueprint-vet

fmt:
	go fmt ./...

lint:
	go vet ./...

tidy:
	go mod tidy
