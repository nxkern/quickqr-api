.PHONY: build dev clean


build:
	go build -ldflags="-s -w" -o bin/quickqr .

dev:
	go run . $(ARGS)

clear:
	rm -rf ./bin
