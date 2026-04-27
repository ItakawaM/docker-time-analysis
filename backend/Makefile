BINARY ?= server.exe

.PHONY: run build help

run:
	go run ./cmd/server/main.go

build:
	go build -o $(BINARY) ./cmd/server/main.go 

clean:
	rm -f ./$(BINARY)

help:
	@echo "Commands:"
	@echo "  make run   - Runs the CLI for testing purposes"
	@echo "  make build - Builds the CLI into a binary"
	@echo "  make clean - Removes the built CLI binary"