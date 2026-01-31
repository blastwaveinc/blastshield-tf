HOSTNAME=registry.terraform.io
NAMESPACE=blastwave
NAME=blastshield
BINARY=terraform-provider-${NAME}
OS_ARCH=$(shell go env GOOS)_$(shell go env GOARCH)

# OpenAPI version to use for code generation
OPENAPI_VERSION ?= latest
OPENAPI_SPEC := openapi/$(OPENAPI_VERSION).json

# Provider version matches API version from OpenAPI spec
VERSION=$(shell jq -r '.info.version' $(OPENAPI_SPEC))

# Test configuration - override these via environment variables
BLASTSHIELD_HOST ?= http://localhost:4999
BLASTSHIELD_TOKEN ?= dev

default: install

build: generate
	go build -o ${BINARY}

release:
	GOOS=darwin GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_darwin_amd64
	GOOS=darwin GOARCH=arm64 go build -o ./bin/${BINARY}_${VERSION}_darwin_arm64
	GOOS=linux GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_linux_amd64
	GOOS=linux GOARCH=arm64 go build -o ./bin/${BINARY}_${VERSION}_linux_arm64
	GOOS=windows GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_windows_amd64.exe

install: build
	mkdir -p ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}
	mv ${BINARY} ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}

# Run unit tests (no API required)
test:
	OPENAPI_VERSION=$(OPENAPI_VERSION) go test -v ./...

# Run all acceptance tests against a live API
testacc:
	OPENAPI_VERSION=$(OPENAPI_VERSION) BLASTSHIELD_HOST=$(BLASTSHIELD_HOST) BLASTSHIELD_TOKEN=$(BLASTSHIELD_TOKEN) \
		TF_ACC=1 go test -v ./internal/provider/... -timeout 120m

# Run acceptance tests for a specific resource (e.g., make testacc-node)
testacc-%:
	OPENAPI_VERSION=$(OPENAPI_VERSION) BLASTSHIELD_HOST=$(BLASTSHIELD_HOST) BLASTSHIELD_TOKEN=$(BLASTSHIELD_TOKEN) \
		TF_ACC=1 go test -v ./internal/provider/... -run 'TestAcc$*' -timeout 30m

# Download OpenAPI spec from a Blastshield orchestrator
# Automatically extracts version and saves to openapi/vX.X.X.json
fetch-openapi:
	@echo "Downloading OpenAPI spec from $(BLASTSHIELD_HOST)..."
	@curl -s "$(BLASTSHIELD_HOST)/openapi.json" | jq . > openapi-temp.json && \
		VERSION=$$(jq -r '.info.version' openapi-temp.json) && \
		if [ -z "$$VERSION" ] || [ "$$VERSION" = "null" ]; then \
			echo "Error: Could not extract version from OpenAPI spec" && \
			rm openapi-temp.json && \
			exit 1; \
		fi && \
		TARGET="openapi/v$$VERSION.json" && \
		mv openapi-temp.json "$$TARGET" && \
		echo "✓ Successfully downloaded and saved to $$TARGET" && \
		echo "" && \
		echo "Version: $$VERSION" && \
		echo "" && \
		echo "Next steps:" && \
		echo "  1. Update symlink (if this is the latest):  cd openapi && ln -sf v$$VERSION.json latest.json" && \
		echo "  2. Generate provider:                       make generate OPENAPI_VERSION=v$$VERSION"

# Generate code from OpenAPI spec
# Usage: make generate OPENAPI_VERSION=v1.13
generate:
	@echo "Generating provider code for API version $(OPENAPI_VERSION)"
	@if [ ! -f "$(OPENAPI_SPEC)" ]; then \
		echo "Error: OpenAPI spec not found: $(OPENAPI_SPEC)"; \
		echo "Available versions:"; \
		ls -1 openapi/*.json | xargs -n1 basename; \
		exit 1; \
	fi
	python3 -m venv .venv
	.venv/bin/pip install --quiet jinja2
	.venv/bin/python generate.py --spec $(OPENAPI_SPEC)
	rm -rf .venv

fmt:
	go fmt ./...

lint:
	golangci-lint run

# Generate documentation from provider schemas
docs: build
	go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs@latest generate --provider-name blastshield

clean:
	rm -f ${BINARY}
	rm -rf bin/
	rm -rf internal/provider/generated/
	rm -f internal/provider/provider.go
	rm -f provider.test
	rm -f .venv
	rm -rf __pycache__

# Test API connectivity and authentication
test-api:
	@echo "Testing API connection to $(BLASTSHIELD_HOST)..."
	@curl -s -f -H "Authorization: Bearer $(BLASTSHIELD_TOKEN)" "$(BLASTSHIELD_HOST)/nodes/" > /dev/null && \
		echo "✓ API connection successful - authentication working" || \
		(echo "✗ API connection failed - check BLASTSHIELD_HOST and BLASTSHIELD_TOKEN" && exit 1)

# Cleanup test entities from the API (useful after test failures)
# Requires curl and jq to be installed
cleanup-test-entities:
	@echo "Cleaning up test entities with tag 'blastshield_tf_testing_entity'..."
	@echo "This requires a running API server and proper authentication."
	@echo "Set BLASTSHIELD_HOST and BLASTSHIELD_TOKEN environment variables if needed."

.PHONY: build release install test testacc test-api fetch-openapi generate fmt lint docs clean cleanup-test-entities
