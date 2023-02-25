all: fmt build

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: build
build:
	go build ./cmd/cpcli/

.PHONY: install
install:
	go install ./cmd/cpcli/
