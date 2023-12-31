-include .env

PACKAGE := "github.com/siyul-park/jsonata-go"

GO_PACKAGE := $(shell go list ${PACKAGE}/...)

.PHONY: init
init:
	@go install -v ${GO_PACKAGE}

.PHONY: init-staticcheck
init-staticcheck:
	@go install honnef.co/go/tools/cmd/staticcheck@latest

.PHONY: generate
generate:
	@go generate ${GO_PACKAGE}

.PHONY: build
build:
	@go clean -cache
	@mkdir -p dist
	@go build -o dist ./...

.PHONY: clean
clean:
	@go clean -cache
	@rm -rf dist

.PHONY: tidy
tidy:
	@go mod tidy

.PHONY: check
check: lint test

.PHONY: test
test:
	@go test $(test-options) ${GO_PACKAGE}

.PHONY: race
race:
	@go test -race $(test-options) ${GO_PACKAGE}

.PHONY: coverage
coverage:
	@go test -coverprofile coverage.out -covermode count ${GO_PACKAGE}
	@go tool cover -func=coverage.out | grep total

.PHONY: benchmark
benchmark:
	@go test -run="-" -bench=".*" -benchmem ${GO_PACKAGE}

.PHONY: lint
lint: fmt vet staticcheck

.PHONY: vet
vet:
	@go vet ${GO_PACKAGE}

.PHONY: fmt
fmt:
	@go fmt ${GO_PACKAGE}

.PHONY: staticcheck
staticcheck: init-staticcheck
	@staticcheck ${GO_PACKAGE}

