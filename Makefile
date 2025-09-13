.PHONY: dep

## Ensure all dependencies are up to date
deps:
	@go mod tidy && go mod download