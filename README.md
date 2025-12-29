# Rusky - Technical Debt Manager

[![CI](https://github.com/titobsala/rusky/workflows/CI/badge.svg)](https://github.com/titobsala/rusky/actions/workflows/ci.yml)
[![CodeQL](https://github.com/titobsala/rusky/workflows/CodeQL/badge.svg)](https://github.com/titobsala/rusky/actions/workflows/codeql.yml)
[![Release](https://github.com/titobsala/rusky/workflows/Release/badge.svg)](https://github.com/titobsala/rusky/actions/workflows/release.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/titobsala/rusky)](https://goreportcard.com/report/github.com/titobsala/rusky)

A simple, language-agnostic TUI/CLI tool for tracking technical debt in your projects.

## Features

- **Manual Debt Tracking**: Add and manage technical debt items via CLI commands
- **Interactive TUI**: Browse and manage debt items with a beautiful terminal interface
- **Mark Complete**: Toggle completion status of debt items
- **Local Storage**: Stores debt items in `.rusky.json` in your project root (VCS-friendly)
- **Language-Agnostic**: Works with any codebase

## Installation

### Download Pre-built Binaries (Recommended)

Download the latest release for your platform from the [releases page](https://github.com/titobsala/rusky/releases).

**Linux (amd64):**
```bash
curl -LO https://github.com/titobsala/rusky/releases/latest/download/rusky-linux-amd64.tar.gz
tar xzf rusky-linux-amd64.tar.gz
sudo mv rusky /usr/local/bin/
```

**macOS (Apple Silicon):**
```bash
curl -LO https://github.com/titobsala/rusky/releases/latest/download/rusky-darwin-arm64.tar.gz
tar xzf rusky-darwin-arm64.tar.gz
sudo mv rusky /usr/local/bin/
```

**Windows:**
Download the zip file from releases and extract to your desired location.

### Using Go Install

```bash
go install github.com/tito-sala/rusky/cmd/rusky@latest
```

### From Source

```bash
git clone https://github.com/titobsala/rusky.git
cd rusky
make build
sudo mv rusky /usr/local/bin/
```

## Quick Start

```bash
# Add your first technical debt item
rusky add "refactor authentication module"

# Add more items
rusky add "fix memory leak in data processor"
rusky add "update deprecated API endpoints"

# List all items (non-interactive)
rusky list

# Mark an item as completed (by index or UUID)
rusky complete 1

# Launch interactive TUI (default when run without arguments)
rusky
```

## Usage

### Commands

#### `rusky` (Interactive TUI)

Launch the interactive terminal UI to browse and manage your technical debt.

**Keyboard shortcuts:**
- `↑/↓` or `k/j` - Navigate between items
- `Enter` or `Space` - Toggle completion status
- `q` or `Esc` - Quit

#### `rusky add <description>`

Add a new technical debt item.

```bash
rusky add "refactor authentication module"
```

#### `rusky complete <id|index>`

Mark a debt item as completed. You can use either the UUID or the 1-based index.

```bash
# Complete by index
rusky complete 1

# Complete by UUID
rusky complete 47085ae2-3240-4fac-a853-5c1400109580
```

#### `rusky list`

Display all technical debt items in a simple text format.

```bash
rusky list
```

#### `rusky version`

Show version information.

```bash
rusky --version
```

## Storage

Rusky stores technical debt items in a `.rusky.json` file in your current working directory (typically your project root).

**Example `.rusky.json`:**

```json
{
  "version": "0.1.0",
  "items": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "description": "refactor auth module",
      "status": "open",
      "created_at": "2025-12-29T10:30:00Z"
    },
    {
      "id": "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
      "description": "update deprecated API endpoints",
      "status": "completed",
      "created_at": "2025-12-28T14:20:00Z",
      "completed_at": "2025-12-29T09:15:00Z"
    }
  ]
}
```

This file is:
- **VCS-friendly**: Commit it to track technical debt alongside your code
- **Human-readable**: JSON format makes it easy to read and edit manually if needed
- **Project-specific**: Each project has its own `.rusky.json` file

## Roadmap

### v0.2.0
- Automatic scanning for TODO/FIXME/HACK comments in codebase
- Filter and sort functionality (by status, date, priority)
- Add priority/tags to debt items

### v0.3.0
- Export to other formats (CSV, Markdown)
- Integration with popular issue trackers (GitHub Issues, Jira)

### Future
- Team sync capabilities
- Dashboard with metrics (debt velocity, completion rate)
- VS Code extension

## Development

### Prerequisites

- Go 1.25 or later
- golangci-lint (for linting)

### Building from Source

```bash
# Clone the repository
git clone https://github.com/titobsala/rusky.git
cd rusky

# Build with version injection
make build

# Or build manually
go build -ldflags "-X github.com/tito-sala/rusky/internal/cli.version=$(git describe --tags --always)" -o rusky cmd/rusky/main.go
```

### Makefile Commands

The project includes a Makefile for common development tasks:

```bash
make build      # Build optimized binary with version injection
make dev        # Build development binary (with debug info)
make test       # Run tests with race detection
make coverage   # Generate HTML coverage report
make lint       # Run golangci-lint
make fmt        # Format code with gofmt and goimports
make install    # Install binary to GOPATH/bin
make clean      # Remove build artifacts
make build-all  # Build for all platforms (Linux, macOS, Windows)
make help       # Show all available commands
```

### Running Tests

```bash
# Run all tests
make test

# Run tests with coverage
make coverage

# Run tests manually
go test -v -race ./...
```

### Code Quality

Before submitting a pull request, ensure your code passes all checks:

```bash
# Run linter
make lint

# Format code
make fmt

# Run tests
make test
```

The CI pipeline will automatically run these checks on all pull requests.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

### Development Workflow

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Run tests and linter (`make test lint`)
5. Commit your changes (`git commit -m 'Add amazing feature'`)
6. Push to the branch (`git push origin feature/amazing-feature`)
7. Open a Pull Request

### CI/CD

All pull requests must pass the following checks:
- **Linting**: Code must pass golangci-lint checks
- **Tests**: All tests must pass with race detection
- **Build**: Code must build successfully for all platforms
- **CodeQL**: Security analysis must pass

You can run these checks locally using the Makefile commands above.

## License

MIT License - see [LICENSE](LICENSE) file for details.

## Built With

- [Go](https://golang.org/)
- [Bubbletea](https://github.com/charmbracelet/bubbletea) - TUI framework
- [Lipgloss](https://github.com/charmbracelet/lipgloss) - Terminal styling
- [Cobra](https://github.com/spf13/cobra) - CLI framework

## Author

Tito Sala
