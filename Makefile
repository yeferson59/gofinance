.PHONY: lint tidy test cover fuzz fuzz-long fmt bench release

TAG ?=
DESCRIPTION ?=

# Minimum total statement coverage. A ratchet, not a target: raise it as
# coverage rises so a change can never quietly give ground.
COVERAGE_MIN ?= 90

# How long each fuzz target runs in the short sweep.
FUZZTIME ?= 10s

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
	@./scripts/check_coverage.sh coverage.out $(COVERAGE_MIN)
	@go tool cover -html=coverage.out

# cover runs the same tests as test but reports the numbers instead of opening
# the browser, which is what CI wants.
cover:
	@go test -race -coverpkg=$(shell echo $(PKGS) | tr ' ' ',') -coverprofile=coverage.out -covermode=atomic $(PKGS)
	@go tool cover -func=coverage.out
	@./scripts/check_coverage.sh coverage.out $(COVERAGE_MIN)

# fuzz runs every fuzz target briefly, enough to catch a regression that the
# seed corpus alone would miss. Run it in CI on every change.
fuzz:
	@for pkg in $(PKGS); do \
		for target in $$(go test -list 'Fuzz.*' $$pkg 2>/dev/null | grep '^Fuzz' || true); do \
			echo "fuzzing $$target in $$pkg"; \
			go test $$pkg -run '^$$' -fuzz "^$$target$$" -fuzztime=$(FUZZTIME) || exit 1; \
		done; \
	done

# fuzz-long is the same sweep with a much longer budget, for a nightly or
# weekly run rather than every push.
fuzz-long:
	@$(MAKE) fuzz FUZZTIME=5m

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
