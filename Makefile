HOSTNAME=registry.terraform.io
NAMESPACE=blastwave
NAME=blastshield
BINARY=terraform-provider-${NAME}
VERSION=0.1.0
OS_ARCH=$(shell go env GOOS)_$(shell go env GOARCH)

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
	go test -v ./...

# Run all acceptance tests against a live API
testacc:
	BLASTSHIELD_HOST=$(BLASTSHIELD_HOST) BLASTSHIELD_TOKEN=$(BLASTSHIELD_TOKEN) \
		TF_ACC=1 go test -v ./internal/provider/... -timeout 120m

# Run acceptance tests for a specific resource (e.g., make testacc-node)
testacc-%:
	BLASTSHIELD_HOST=$(BLASTSHIELD_HOST) BLASTSHIELD_TOKEN=$(BLASTSHIELD_TOKEN) \
		TF_ACC=1 go test -v ./internal/provider/... -run 'TestAcc$*' -timeout 30m

# Download OpenAPI spec from the Blastshield orchestrator
fetch-openapi:
	@curl -s "$(BLASTSHIELD_HOST)/openapi.json" -o openapi.json && \
		echo "Successfully downloaded openapi.json from $(BLASTSHIELD_HOST)" && \
		echo 'Note: You may need to manually modify the "version" field to reflect your API version for testing'

# Generate code from OpenAPI spec
generate:
	python3 -m venv .venv
	.venv/bin/pip install --quiet jinja2
	.venv/bin/python generate.py
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

# Cleanup test entities from the API (useful after test failures)
# Requires curl and jq to be installed
cleanup-test-entities:
	@echo "Cleaning up test entities with tag 'blastshield_tf_testing_entity'..."
	@echo "This requires a running API server and proper authentication."
	@echo "Set BLASTSHIELD_HOST and BLASTSHIELD_TOKEN environment variables if needed."

.PHONY: build release install test testacc fetch-openapi generate fmt lint docs clean cleanup-test-entities
