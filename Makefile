.PHONY: all
all: vendor

.PHONY: vet
vet:
	go vet ./...

.PHONY: build
build:
	go build ./cmd/gowet

.PHONY: test
test:
	go test -v -cover -race ./...

.PHONY: run
run:
	go run ./cmd/gowet

