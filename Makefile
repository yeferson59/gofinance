.PHONY: lint tidy test fmt bench release

TAG ?=
DESCRIPTION ?=

lint:
	@golangci-lint --config=.golangci.yaml run ./... -v

tidy:
	@go mod tidy
	@cd charts && go mod tidy
	@cd examples && go mod tidy

PKGS := $(shell go list ./...)

test:
	@go test -race -coverpkg=$(shell echo $(PKGS) | tr ' ' ',') -coverprofile=coverage.out -covermode=atomic $(PKGS)
	@cd charts && go test -race ./...
	@cd examples && go build ./...
	@go tool cover -html=coverage.out

fmt:
	@go fmt ./...
	@golangci-lint --config=.golangci.yaml fmt ./... -v
	@goimports -w  -v .

bench:
	@go test -bench=. -benchmem ./...

release:
	@make fmt
	@make lint
	@make test
	@if [ -z "$(TAG)" ] || [ -z "$(DESCRIPTION)" ]; then \
			echo 'Usage: make release TAG=<tag_name> [DESCRIPTION="description of changes"]'; \
			exit 1; \
	fi

	git tag -a v$(TAG) -m "$(DESCRIPTION)"
	git push origin v$(TAG)
