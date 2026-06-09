.PHONY: help default githooks build fmt goimports govet golint terrafmt install test testacc test-cover

TEST?=$$(go list ./... |grep -v 'vendor')
# filesystem-mirror layout: <REGISTRY>/<NAMESPACE>/<NAME> = source address registry.terraform.io/layerx/twilio
REGISTRY=registry.terraform.io
NAMESPACE=layerx
NAME=twilio
BINARY=terraform-provider-${NAME}
VERSION=0.18.46
# Resolve the host OS/arch dynamically (e.g. darwin_arm64) instead of hard-coding.
OS_ARCH=$(shell go env GOOS)_$(shell go env GOARCH)
GOLANGCI_LINT_VERSION=v2.12.2

.DEFAULT_GOAL := help

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

default: build

githooks: ## Install the pre-commit git hook
	ln -sf ../../githooks/pre-commit .git/hooks/pre-commit

build: ## Build the provider binary (no source mutation; run `make fmt` to format)
	go build -o ${BINARY}

fmt: goimports terrafmt ## Format Go imports/modules and Terraform example files

goimports: ## Format imports and tidy modules
	go install golang.org/x/tools/cmd/goimports@v0.24.0
	goimports -w .
	go mod tidy

govet: goimports ## Run go vet
	go vet ./...

golint: govet ## Run golangci-lint
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@${GOLANGCI_LINT_VERSION}
	golangci-lint run

terrafmt: ## Format Terraform example files
	terraform fmt -recursive

install: build ## Build and install into the local filesystem mirror
	mkdir -p ~/.terraform.d/plugins/${REGISTRY}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}
	mv ${BINARY} ~/.terraform.d/plugins/${REGISTRY}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}/${BINARY}_v${VERSION}

test: build ## Run unit tests
	go test $(TEST) || exit 1
	echo $(TEST) | \
		xargs -t -n4 go test $(TESTARGS) -timeout=30s -parallel=4

testacc: ## Run acceptance tests (creates real, billable Twilio resources)
	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 120m

GO_DIRS = $(shell go list ./... | grep -v /resources/)
test-cover: ## Run tests with coverage output
	go test ${GO_DIRS} -coverprofile coverage.out
	go test ${GO_DIRS} -json > test-report.out
