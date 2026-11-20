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
	@if command -v lefthook >/dev/null 2>&1;then \
		lefthook install; \
	elif command -v go >/dev/null 2>&1; then \
		echo "installing lefthook..."; \
		go install github.com/evilmartians/lefthook/v2@v2.2.2 && \
		"$$(go env GOPATH)/bin/lefthook" install; \
	else \
		echo "Go is not installed, please install lefthook manually"; \
		exit 1; \
	fi

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
