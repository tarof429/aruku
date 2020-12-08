# A Self-Documenting Makefile: http://marmelab.com/blog/2016/02/29/auto-documented-makefile.html

.PHONY: all build help test

all: build ## Build all targets

help: ## Display this help
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'


build: ## Build all targets
	env GOOS=linux GOOARCH=amd64 go build -o aruku-bin-amd64
	env GOOS=darwin GOOARCH=amd64 go build -o aruku-bin-darwin64

test: ## Run all tests
	go test -v ./...

clean:
	@rm -f aruku-bin-amd64 aruku-bin-darwin64