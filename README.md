# zellij-tui

A terminal session manager for [Zellij](https://zellij.dev), built with [Bubble Tea v2](https://charm.land/bubbletea/v2).

Browse, create, attach, delete, and kill Zellij sessions from a single keyboard-driven interface.

## Features

- **Session list** — view all sessions with creation times and current-session indicator
- **Create sessions** — press `n`, type a name, hit enter
- **Attach** — select a session and press enter to `exec` directly into it
- **Delete / Kill** — delete (`d`) or force-kill (`x`) individual sessions, or wipe all (`D`)
- **Exited session support** — see EXITED sessions and resurrect them by attaching
- **Auto-attach** — creating a session immediately attaches (configurable)
- **Works inside sessions** — shows a warning banner; prevents attaching to the current session
- **No wrapper layer** — uses `syscall.Exec` so zellij replaces the process directly

## Requirements

- [Go](https://go.dev) 1.25+
- [Zellij](https://zellij.dev) installed and in `PATH`

## Install

```bash
go install github.com/sadiq/zellij-tui@latest
```

## Usage

Run the binary directly:

```bash
zellij-tui
```

Or use the shell helper for a shorter command — add this to your `.bashrc` / `.zshrc`:

```bash
source /path/to/zellij-tui/shell/zt.sh
```

Then run:

```bash
zt
```

### Auto-launch on terminal open

To start the session manager every time you open a terminal (when not already in a Zellij session):

```bash
if [[ -z "$ZELLIJ" ]] && [[ $- == *i* ]]; then
    zt
fi
```

## Keybindings

| Key       | Action                            |
|-----------|-----------------------------------|
| `↑` / `k` | Move cursor up                    |
| `↓` / `j` | Move cursor down                  |
| `Enter`   | Attach to selected session        |
| `n`       | Create a new session              |
| `d`       | Delete selected session           |
| `x`       | Kill selected session             |
| `D`       | Delete all sessions               |
| `r`       | Refresh session list              |
| `q`       | Quit (press twice to confirm)     |
| `Ctrl+C`  | Quit immediately                  |
| `Esc`     | Cancel / go back                  |

## Configuration

Config file: `~/.config/zellij-tui/config.toml` (or `$XDG_CONFIG_HOME/zellij-tui/config.toml`)

```toml
# Path to the zellij binary (optional, defaults to PATH lookup)
zellij_path = "/usr/bin/zellij"

# Auto-attach when creating a new session (default: true)
# true  = immediately attach to the new session
# false = create in background, stay in TUI
auto_attach = true
```

## How it works

1. Launches a Bubble Tea TUI rendered to `/dev/tty` (works even when stdout is piped)
2. Calls `zellij list-sessions --no-formatting` to populate the session list
3. On exit, uses `syscall.Exec` to replace the process with the selected `zellij attach` command — no shell wrapper overhead

## Build from source

```bash
git clone https://github.com/sadiq/zellij-tui.git
cd zellij-tui
go build -o zellij-tui .
```

## License

MIT
