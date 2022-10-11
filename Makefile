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

.PHONY: cover
cover:
	GO111MODULE=off go get github.com/axw/gocov/gocov
	GO111MODULE=off go get -u gopkg.in/matm/v1/gocov-html
	${GOPATH}/bin/gocov test ./... | ${GOPATH}/bin/gocov-html > coverage.html
	open coverage.html

.PHONY: run
run:
	go run main.go _input/sample-script.txt _input/sample-script2.txt