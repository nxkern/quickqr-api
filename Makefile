.PHONY: build dev clean test tidy

OPTION ?= pkgname

build:
	go build -ldflags="-s -w" -o bin/quickqr .

dev:
	go run . $(ARGS)

clear:
	rm -rf ./bin

test:
	gotestsum --format $(OPTION)   

tidy: 
	go mod tidy
