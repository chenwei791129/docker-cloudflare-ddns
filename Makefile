.PHONY: build run clean help

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-10s %s\n", $$1, $$2}'

build: ## Build binary to bin/
	@mkdir -p bin
	CGO_ENABLED=0 go build -o bin/cloudflare-ddns .

run: build ## Build and run with .env
	@set -a && . ./.env && set +a && ./bin/cloudflare-ddns

clean: ## Remove bin/ directory
	rm -rf bin
