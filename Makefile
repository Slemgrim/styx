PROJECT := styx

all: help

help:
	@echo "make build - build styx-api and styx-worker"
	@echo "make run-api - run the styx-api with go run"
	@echo "make run-worker - run the styx-worker with go run"

build:
	GOOS=darwin GOARCH=amd64 go build -o build/darwin/styx-api ./cmd/styx-api
	GOOS=darwin GOARCH=amd64 go build -o build/darwin/styx-worker ./cmd/styx-worker
	GOOS=linux GOARCH=amd64 go build -o build/linux/styx-api ./cmd/styx-api
	GOOS=linux GOARCH=amd64 go build -o build/linux/styx-worker ./cmd/styx-worker

run-api:
	go run ./cmd/styx-api/main.go

run-worker:
	go run ./cmd/styx-worker/main.go