.DEFAULT_GOAL := build
.PHONY: test lint vet race bench

test:
	go test ./... -cover

lint:
	golangci-lint run ./...

vet:
	go vet ./...

race:
	go test -race ./...

bench:
	go test -bench=.
