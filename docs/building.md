# Building zellij-tui

## Prerequisites

| Dependency | Version   | Notes                              |
|------------|-----------|------------------------------------|
| Go         | 1.25+     | Required for `charm.land/bubbletea/v2` |
| Zellij     | any       | Runtime only; not needed to build  |
| Git        | any       | For cloning the repo               |

### Installing Go

If you don't have Go 1.25+:

```bash
# Direct download (Linux amd64)
curl -Lo /tmp/go.tar.gz https://go.dev/dl/go1.25.7.linux-amd64.tar.gz
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf /tmp/go.tar.gz

# Add to PATH (in .bashrc / .zshrc)
export PATH=$PATH:/usr/local/go/bin
```

Or use your system's package manager if it provides Go 1.25+.

## Clone

```bash
git clone https://github.com/sadiq/zellij-tui.git
cd zellij-tui
```

## Build

```bash
go build -o zellij-tui .
```

This produces a single static binary named `zellij-tui` in the current directory.

### Cross-compilation

Build for a different OS/architecture:

```bash
# Linux ARM64 (e.g. Raspberry Pi)
GOOS=linux GOARCH=arm64 go build -o zellij-tui-arm64 .

# macOS Apple Silicon
GOOS=darwin GOARCH=arm64 go build -o zellij-tui-macos .

# macOS Intel
GOOS=darwin GOARCH=amd64 go build -o zellij-tui-macos-intel .
```

### Build flags

```bash
# Strip debug info for a smaller binary
go build -ldflags="-s -w" -o zellij-tui .

# Embed version information
VERSION=$(git describe --tags --always --dirty)
go build -ldflags="-X main.version=$VERSION" -o zellij-tui .
```

> **Note:** The `main.version` variable doesn't exist in the current codebase. Add it to `main.go` if you want version tracking:
> ```go
> var version = "dev"
> ```

## Install

```bash
# Install to $GOPATH/bin (or $HOME/go/bin if GOPATH is unset)
go install .

# Or copy the built binary manually
cp zellij-tui /usr/local/bin/
```

Make sure `$GOPATH/bin` or the install target is in your `PATH`.

## Run tests

```bash
go test ./...
```

Only the `internal/zellij` package has tests currently (session output parsing).

## Run without installing

```bash
go run .
```

## Verify

After building or installing, confirm it works:

```bash
# Should print an error about zellij not found (if zellij isn't installed)
# or open the TUI (if zellij is installed)
zellij-tui
```

## Dependencies

Dependencies are managed via Go modules. To update:

```bash
# Download dependencies
go mod download

# Tidy unused dependencies
go mod tidy
```

Key runtime dependencies (imported by the project):

| Package                     | Purpose                    |
|-----------------------------|----------------------------|
| `charm.land/bubbletea/v2`  | TUI framework              |
| `charm.land/bubbles/v2`    | Text input component       |
| `charm.land/lipgloss/v2`   | Terminal styling/layout    |
| `github.com/pelletier/go-toml/v2` | Config file parsing  |

## Troubleshooting

### `go: module charm.land/bubbletea/v2: unrecognized import path`

Go 1.25+ is required. Check your version:

```bash
go version
```

### `Error: cannot open /dev/tty`

The TUI renders to `/dev/tty`. This fails in non-interactive contexts (CI, piped commands). Run it from an interactive terminal.

### `Error: zellij not found in PATH`

Either:
- Install zellij: follow instructions at [zellij.dev](https://zellij.dev)
- Or set `zellij_path` in the config file (see [configuration docs](configuration.md))
