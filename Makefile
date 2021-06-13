.PHONY: all
all: vendor

.PHONY: vendor
vendor:
	go mod vendor

.PHONY: run
run:
	go run ./cmd/gowet

.PHONY: build
build:
	go build ./cmd/gowet

.PHONY: test
test:
	go test -v -race -run ./...

