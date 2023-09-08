.PHONY: lint
lint:
	golangci-lint run --config .golangci.yml

.PHONY: build
build:
	go build -o build/lighting_user_vault ./core/main.go