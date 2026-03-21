.PHONY: lint tidy test fmt bench release

TAG ?=
DESCRIPTION ?=

lint:
	@golangci-lint --config=.golangci.yaml run ./... -v

tidy:
	@go mod tidy

test:
	@go test -race -coverpkg=./... -coverprofile=coverage.out -covermode=atomic ./...
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

	git tag -a $(TAG) -m "$(DESCRIPTION)"
