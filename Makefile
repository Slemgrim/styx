PROJECT := styx

all: help

help:
	@echo "make build - build styx-api and styx-worker"
	@echo "make run-api - run the styx-api with go run"
	@echo "make run-worker - run the styx-worker with go run"

build:
	go build -o styx-api ./cmd/styx-api
	go build -o styx-worker ./cmd/styx-worker

run-api:
	go run ./cmd/styx-api/main.go

run-worker:
	go run ./cmd/styx-worker/main.go