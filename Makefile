.PHONY: build dev clean test tidy fmt lint debug vulnerability-check dive

OPTION ?= pkgname

build:
	go build -ldflags="-s -w" -o bin/quickqr .

dev:
	air -c .air.toml

clear:
	rm -rf ./bin

test:
	gotestsum --format $(OPTION)   

tidy: 
	go mod tidy

fmt:
	golangci-lint fmt

lint:
	golangci-lint run

debug:
	dlv debug .

vulnerability-check:
	govulncheck ./...

dive:
	dive $(DOCKER_IMAGE)
