# Variables
PROJECT_ROOT := $(shell pwd)
BINARY_PATH := $(PROJECT_ROOT)/build/lighting_user_vault
E2E_DIR := ./core/e2e

.PHONY: e2e
e2e: build
	BINARY_PATH=$(BINARY_PATH) go test $(E2E_DIR)/...

.PHONY: lint
lint:
	golangci-lint run --config .golangci.yml

.PHONY: build
build:
	go build -o build/lighting_user_vault ./core/main.go

.PHONY: unit
unit:
	go test $$(go list ./... | grep -v $(E2E_DIR))

.PHONY: fuzz
fuzz:
	chmod +x ./scripts/run_all_fuzz_tests.sh
	./scripts/run_all_fuzz_tests.sh