# Contributing

Contributions are welcome. Please follow these guidelines.

## Development Setup

1. Clone the repository:

```bash
git clone https://github.com/anaryk/metabase-mcp-server.git
cd metabase-mcp-server
```

2. Install dependencies:

```bash
go mod download
```

3. Run tests:

```bash
go test -v -race ./...
```

4. Run linter:

```bash
golangci-lint run
```

## Making Changes

1. Fork the repository and create a feature branch.
2. Write tests for any new functionality.
3. Ensure all tests pass with `go test -race ./...`.
4. Ensure code passes linting with `golangci-lint run`.
5. Submit a pull request with a clear description of the changes.

## Code Style

- Follow standard Go conventions and `gofmt` formatting.
- Use structured logging with zerolog.
- Write table-driven tests where applicable.
- Keep functions focused and concise.

## Reporting Issues

Open an issue on GitHub with a clear description of the problem, including steps to reproduce and expected behavior.
