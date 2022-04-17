VERSION := $(shell cat VERSION)
BUILD_DIR = .build

GO := $(GOROOT)/bin/go
ifeq (, $(GOROOT))
	GO := go
endif

GOFILES_NOVENDOR := $(shell find . -type f -name '*.go' -not -path "*/vendor/*" 2>/dev/null)
GOPACKAGES       := $(shell $(GO) list ./...)
GOIMPORTS_REPO   := golang.org/x/tools/cmd/goimports
GOIMPORTS        := $(GOPATH)/bin/goimports

.PHONY: fmt
fmt: ## Ensures all go files are properly formatted.
	@echo "Formatting..."
	@GO111MODULE=off $(GO) get -u ${GOIMPORTS_REPO}
	@$(GOIMPORTS) -l -w ${GOFILES_NOVENDOR}

.PHONY: vendor
vendor: export GO111MODULE=on
vendor: export GOPRIVATE=github.corp.ebay.com
vendor: # Ensures all go module dependencies are synced and copied to vendor
	@echo "Updating module dependencies..."
	@$(GO) mod tidy
	@$(GO) mod vendor

.PHONY: check
check: ## Performs code hygiene checks and runs tests.
	@echo "Checking for suspicious constructs..."
	@$(GO) vet ${GOPACKAGES}
	@echo "Checking formatting..."
	@GO111MODULE=off $(GO) get -u ${GOIMPORTS_REPO}
	@$(GOIMPORTS) -l ${GOFILES_NOVENDOR} | (! grep .) || (echo "Code differs from goimports' style ^" && false)
	@$(MAKE) test-local

.PHONY: test-local # Run local go test
test-local: ## Runs tests.
	@echo "Running tests..."
	@$(GO) test -short -race -gcflags=-l ${GOPACKAGES}

.PHONY: coverage
coverage: ## Runs tests with coverage.
	@echo "Running tests with coverage..."
	@echo ${GOPACKAGES} | xargs $(GO) test -short -race -coverprofile=coverage.txt -covermode=atomic && $(GO) tool cover -html=coverage.txt -o coverage.html

.PHONY: clean
clean:
	@rm -rf ${BUILD_DIR}

