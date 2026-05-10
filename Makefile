.PHONY: test lint lint-fix fmt fmt-diff

GOLANGCI_LINT=go tool -modfile tools/go.mod golangci-lint

test:
	go test -v -race ./...

lint:
	$(GOLANGCI_LINT) run ./...

lint-fix:
	$(GOLANGCI_LINT) run --fix ./...

fmt:
	$(GOLANGCI_LINT) fmt ./...

fmt-diff:
	$(GOLANGCI_LINT) fmt ./... --diff

tools:
	cd tools && go mod tidy

mod:
	go mod tidy
