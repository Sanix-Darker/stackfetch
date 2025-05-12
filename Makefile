.PHONY: build run

run:
	go run ./cmd/stackfetch/main.go

build:
	rm -rf ./stackfetch
	go build -ldflags "-s -w -X main.version=$(git describe --tags)" -o stackfetch ./cmd/stackfetch
