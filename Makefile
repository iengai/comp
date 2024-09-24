help:
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

functions := $(shell find functions -name \*main.go | awk -F'/' '{print $$2}')

build: ## Build golang binaries
	@for function in $(functions) ; do \
		env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/bootstrap functions/$$function/main.go ; \
		zip -j bin/$$function.zip bin/bootstrap; \
	done
