.PHONY: build run test clean

build:
	go build -o bin/agent ./cmd/agent

run:
	go run ./cmd/agent

test:
	go test -v ./...

clean:
	rm -rf bin/
