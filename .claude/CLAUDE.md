# Project Instructions

## Build

Always use `make build` instead of `go build` directly. The Makefile outputs the binary to `bin/` which is in `.gitignore`. Running `go build` without `-o` produces a `./docker-cloudflare-ddns` binary in the project root that can accidentally be committed.

Available make targets:

- `make build` — Build binary to `bin/`
- `make run` — Build and run with `.env`
- `make lint` — Run golangci-lint
- `make clean` — Remove `bin/` directory
