GO_BIN = $(GOPATH)/bin
GO_PKG = $(GOPATH)/pkg
GO_SRC = $(GOPATH)/src

GO = go
GOLINT = $(GO_BIN)/revive

PROJECT_BASE = .

.PHONY: lint
lint: check-gosetup ## Run linter on a project
ifeq (, $(shell which ${GOLINT}))
        $(GO) get -u github.com/mgechev/revive
endif
	@for source in $(shell find ${PROJECT_BASE} -type f -name '*.go' -not -path '*/vendor/*'); do \
		${GOLINT} -config config.toml -formatter stylish $$source; \
	done

.PHONY: fmt
fmt: ## Run go fmt on a project
	$(GO) fmt ./...

.PHONY: unit-test
unit-test: ## Run unit  tests
	$(GO) clean -testcache
	$(GO) test -cover ./internal/pkg/db/... -v

.PHONY: end-to-end-test
end-to-end-test: ## Run end-to-end tests within docker with complete infrastructure
	$(MAKE) -C test/end-to-end

.PHONY: check-gosetup
check-gosetup:
ifndef GOPATH
	$(error GOPATH is undefined)
endif

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := help
