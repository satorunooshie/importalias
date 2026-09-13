.PHONY: all build test vet fmt check

all: check build

build:
	go build ./...

test:
	go test ./...

vet:
	go vet ./...

fmt:
	test -z "$$(gofmt -l importalias.go importalias_test.go cmd/importalias/main.go)"

check: fmt test vet
