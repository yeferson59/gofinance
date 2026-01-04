.PHONY: lint test fmt

lint:
	@golangci-lint --config=.golangci.yaml run ./... -v

test:
	@go test -coverpkg=./... -coverprofile=coverage.out -covermode=atomic ./...
	@go tool cover -html=coverage.out

fmt:
	@go fmt ./...
	@golangci-lint --config=.golangci.yaml fmt ./... -v
	@goimports -w  -v .
