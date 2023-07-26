# Makefile for building a Go program for Linux or Mac
# Usage: make [target]
# Default target: build

# Variables
GO = go
BIN = datamover

# Targets
.PHONY: build-linux build-mac clean

build:
	$(GO) build -o $(BIN)
linux:
	GOOS=linux GOARCH=amd64 $(GO) build -o $(BIN)-linux

mac:
	GOOS=darwin GOARCH=amd64 $(GO) build -o $(BIN)-mac

clean:
	rm -f $(BIN) $(BIN)-linux $(BIN)-mac *.db *.sql