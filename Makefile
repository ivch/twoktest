SHELL=/bin/sh
export GO111MODULE=on

.PHONY: lint
lint:
	golangci-lint run

.PHONY: inastall-lint
install-lint:
	go get github.com/golangci/golangci-lint/cmd/golangci-lint

.PHONY: deps
deps:
	rm -rf vendor
	go mod download
	go mod tidy
