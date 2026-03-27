# Contributing to Boilerworks Go Micro

Thank you for your interest in contributing!

## Getting Started

1. Clone the repository
2. Run `docker compose up -d --build` to start all services
3. Read [bootstrap.md](bootstrap.md) for conventions

## Development Process

1. Fork the repository
2. Create a feature branch from `main`
3. Make your changes
4. Run `make lint` and `make test`
5. Submit a pull request

## Code Style

- Follow standard Go conventions (`gofmt`, `go vet`)
- Use `golangci-lint` for static analysis
- sqlc generates query code -- edit SQL in `db/queries/`, not Go files in `internal/database/queries/`
- All handlers return `ApiResponse` format

## Questions?

Open an issue in this repository.
