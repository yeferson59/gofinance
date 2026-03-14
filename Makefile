.PHONY: lint tidy test fmt bench

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
