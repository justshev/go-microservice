.PHONY: tidy test run-api run-worker build-api build-worker

tidy:
	go mod tidy

test:
	go test ./... -count=1

run-api:
	HTTP_PORT=8080 LOG_LEVEL=info go run ./cmd/api

run-worker:
	LOG_LEVEL=info go run ./cmd/worker

build-api:
	mkdir -p bin
	go build -o bin/api ./cmd/api

build-worker:
	mkdir -p bin
	go build -o bin/worker ./cmd/worker
