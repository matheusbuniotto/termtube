.PHONY: build install test run

build:
	go build -o bin/ytune ./cmd/ytune

install:
	go install ./cmd/ytune

test:
	go test ./...

run: build
	./bin/ytune
