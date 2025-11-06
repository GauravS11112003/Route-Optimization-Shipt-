.PHONY: run build clean test

run:
	go run cmd/main.go

build:
	go build -o bin/shipt-route-optimizer cmd/main.go

clean:
	rm -rf bin/

test:
	go test ./...

deps:
	go mod download
	go mod tidy

