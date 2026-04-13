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
- **Full Nix support** — flake, dev shell, Home Manager module, overlay

## Install

### Nix (recommended)

```bash
# Run without installing
nix run github:pageton/zellij-tui

# Install to your profile
nix profile install github:pageton/zellij-tui
```

For declarative config via Home Manager, see [Nix setup](docs/nix.md).

### go install

```bash
go install github.com/pageton/zellij-tui@latest
```

### Build from source

```bash
git clone https://github.com/pageton/zellij-tui.git
cd zellij-tui
just build
```

Or without [just](https://github.com/casey/just):

```bash
go build -o zellij-tui .
```

## Usage

```bash
zellij-tui
```

Or use the shell helper — add to `.bashrc` / `.zshrc`:

```bash
source /path/to/zellij-tui/shell/zt.sh
```

Then run `zt`.

### Auto-launch on terminal open

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

### Nix (Home Manager)

```nix
programs.zellij-tui = {
  enable = true;
  settings = {
    zellij_path = "${pkgs.zellij}/bin/zellij";
    auto_attach = true;
  };
};
```

Generates the config file and adds `zellij-tui` to your PATH. See the [full Nix docs](docs/nix.md) for setup.

## How it works

1. Launches a Bubble Tea TUI rendered to `/dev/tty` (works even when stdout is piped)
2. Calls `zellij list-sessions --no-formatting` to populate the session list
3. On exit, uses `syscall.Exec` to replace the process with the selected `zellij attach` command — no shell wrapper overhead

## Documentation

- [Keybindings](docs/keybindings.md) — full key reference for all modes
- [Building from source](docs/building.md) — prerequisites, build flags, cross-compilation, troubleshooting
- [Configuration](docs/configuration.md) — config file location, options, defaults
- [Nix setup](docs/nix.md) — flakes, Home Manager module, overlay, direnv

## License

[MIT](LICENSE)
