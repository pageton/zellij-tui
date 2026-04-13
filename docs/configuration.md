# Configuration

zellij-tui uses a single TOML config file. If no config file exists, sensible defaults are used.

## Config file location

The config is loaded from the first path that exists:

1. `$XDG_CONFIG_HOME/zellij-tui/config.toml`
2. `~/.config/zellij-tui/config.toml`

If neither file exists, the default config is used (everything works out of the box).

## Creating the config

```bash
mkdir -p ~/.config/zellij-tui
cat > ~/.config/zellij-tui/config.toml << 'EOF'
# See below for available options
EOF
```

## Options

### `zellij_path`

Path to the zellij binary.

- **Type:** string
- **Default:** auto-detected from `PATH`

When unset or empty, zellij-tui runs `exec.LookPath("zellij")` to find the binary automatically. Set this only if zellij is installed in a non-standard location.

```toml
zellij_path = "/usr/local/bin/zellij"
```

### `auto_attach`

Controls what happens after creating a new session.

- **Type:** boolean
- **Default:** `true`

| Value  | Behavior                                        |
|--------|-------------------------------------------------|
| `true` | Immediately attach to the new session (default) |
| `false`| Create the session in the background, stay in the TUI |

```toml
auto_attach = true
```

## Full example

```toml
# ~/.config/zellij-tui/config.toml

# Path to zellij binary (optional)
zellij_path = "/usr/bin/zellij"

# Auto-attach on session creation (default: true)
auto_attach = true
```

## Defaults (no config file needed)

If you skip the config file entirely, the defaults are:

| Option         | Default            |
|----------------|--------------------|
| `zellij_path`  | PATH lookup        |
| `auto_attach`  | `true`             |

## Validation

- Unknown keys in the config file are silently ignored (no error).
- If the config file contains invalid TOML, defaults are used.
- If `zellij_path` points to a non-existent binary, zellij-tui exits with an error at startup.

## Zellij compatibility

zellij-tui requires that zellij's `on_force_close` is set to `"detach"` (not `"quit"`).

When set to `"quit"`, closing the terminal while attached to a session kills the zellij server process and destroys the session. With `"detach"`, the server survives and the session can be re-attached from zellij-tui or `zellij attach`.

Check your zellij config (`~/.config/zellij/config.kdl` or Home Manager):

```kdl
// kdl config
on_force_close "detach"
```

```nix
# Home Manager
programs.zellij.settings.on_force_close = "detach";
```
