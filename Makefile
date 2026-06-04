.PHONY: build install test run

build:
	go build -o bin/termtube ./cmd/termtube

install:
	go install ./cmd/termtube

test:
	go test ./...

run: build
	./bin/termtube
