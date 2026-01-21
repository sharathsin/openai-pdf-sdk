.PHONY: all build test lint clean run

all: build

build:
	go build -o bin/app cmd/main.go

test:
	go test -race -v ./...

lint:
	golangci-lint run ./...

clean:
	rm -rf bin/
	go clean

run:
	go run cmd/main.go
