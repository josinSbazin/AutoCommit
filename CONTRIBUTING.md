# Contributing to AutoCommit

Thank you for your interest in contributing to AutoCommit! This document provides guidelines and information for contributors.

## Getting Started

1. Fork the repository
2. Clone your fork: `git clone https://github.com/YOUR_USERNAME/autocommit.git`
3. Create a branch: `git checkout -b feature/my-feature`
4. Make your changes
5. Run tests: `make test`
6. Commit your changes (use autocommit!)
7. Push and create a Pull Request

## Development Setup

### Prerequisites

- Go 1.22 or later
- Git
- Make (optional, but recommended)

### Building

```bash
# Build
make build

# Run
./bin/autocommit

# Or run directly
go run ./cmd/autocommit
```

### Testing

```bash
# Run all tests
make test

# Run with coverage
make test-coverage
```

### Linting

```bash
# Install golangci-lint
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Run linter
make lint
```

## Project Structure

```
autocommit/
├── cmd/autocommit/      # Entry point
├── internal/
│   ├── cli/             # CLI commands (Cobra)
│   ├── config/          # Configuration management
│   ├── git/             # Git operations
│   ├── prompt/          # Prompt building for LLMs
│   ├── provider/        # LLM providers
│   └── ui/              # Terminal UI
├── install.sh           # Linux/macOS installer
└── install.ps1          # Windows installer
```

## Adding a New Provider

1. Create a new file in `internal/provider/`:

```go
// internal/provider/myprovider.go
package provider

type MyProvider struct {
    apiKey string
    model  string
}

func NewMyProvider(cfg *config.Config) (*MyProvider, error) {
    apiKey := os.Getenv("MYPROVIDER_API_KEY")
    if apiKey == "" {
        return nil, fmt.Errorf("MYPROVIDER_API_KEY not set")
    }
    return &MyProvider{apiKey: apiKey, model: cfg.Model}, nil
}

func (p *MyProvider) Generate(ctx context.Context, prompt string) (string, error) {
    // Implement API call
}

func (p *MyProvider) Name() string { return "myprovider" }
func (p *MyProvider) Validate() error { return nil }
```

2. Register in `internal/provider/provider.go`:

```go
case "myprovider":
    return NewMyProvider(cfg)
```

3. Add to auto-detection if needed

4. Update documentation

## Code Style

- Follow standard Go conventions
- Use `gofmt` and `goimports`
- Write descriptive commit messages
- Add tests for new functionality
- Keep functions small and focused

## Commit Messages

We use Conventional Commits. Examples:

- `feat(provider): add Mistral support`
- `fix(git): handle binary files in diff`
- `docs: update README with new providers`
- `test(provider): add tests for GigaChat`

## Pull Request Process

1. Ensure all tests pass
2. Update documentation if needed
3. Add entry to CHANGELOG.md
4. Request review from maintainers

## Reporting Issues

- Use GitHub Issues
- Include OS, Go version, and autocommit version
- Provide steps to reproduce
- Include error messages and logs

## License

By contributing, you agree that your contributions will be licensed under the MIT License.
