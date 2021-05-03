.PHONY: setup
setup: ## Setup
	go get github.com/golangci/golangci-lint/cmd/golangci-lint
	go get golang.org/x/tools/cmd/goimports

.PHONY: lint
lint: ## Lint
	golangci-lint run --config configs/golangci.yml

.PHONY: fmt
fmt: ## Format
	gofmt -w -s .
	goimports -w .
	go clean ./...

.PHONY: clean
clean: ## Clean
	go clean ./...
