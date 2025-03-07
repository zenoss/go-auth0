include .env
export

GOFUMPT             := $(shell command -v gofumpt 2> /dev/null)
GOLANGCI_LINT       := $(shell command -v golangci-lint 2> /dev/null)
REVIVE              := $(shell command -v revive 2> /dev/null)

M = $(shell printf "\033[34;1m▶\033[0m")
RED = $(shell printf "\033[31;1m▶\033[0m")

default: check test

.PHONY: check
check: fmt lint

.PHONY: fmt
fmt:
ifeq ($(strip $(GOFUMPT)),)
	@echo "$(RED) Warning: gofumpt is not available on this system, please install it"
else
	@echo "$(M) gofumpt: formatting…"
	@if $(GOFUMPT) -l -w . ; \
	then \
		echo "$(M) gofumpt: files look good"; \
	else \
		echo "$(RED) gofumpt: Please commit formatted files"; \
		exit 1; \
	fi
endif

.PHONY: lint
lint: revive golangci-lint

.PHONY: revive
revive:
ifeq ($(strip $(REVIVE)),)
	@echo "$(RED) Warning: revive is not available on this system, please install it"
else
	@echo "$(M) revive: linting…"
	@$(REVIVE) \
		-config ./.revive.toml \
		-formatter=stylish \
		-exclude ./vendor/... \
		./...
endif

.PHONY: golangci-lint
golangci-lint:
ifeq ($(strip $(GOLANGCI_LINT)),)
	@echo "$(RED) Warning: golangci-lint is not available on this system, please install it"
else
	@echo "$(M) golangci-lint: linting…"
	-@$(GOLANGCI_LINT) run \
		--tests=false \
		--sort-results \
		--presets "unused,performance,bugs" \
		--disable "protogetter,contextcheck,musttag" \
		./...
endif

.PHONY: test
test:
	@echo "$(M) go test: running…"
	@go test ./auth0/... -tags=integration -cover
