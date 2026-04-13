# Keybindings

All keybindings are vim-friendly (both arrow keys and `hjkl`-style keys work for navigation).

## Session list mode

This is the default mode when you launch zellij-tui.

| Key        | Action                                                      |
|------------|-------------------------------------------------------------|
| `↑` / `k`  | Move cursor up                                              |
| `↓` / `j`  | Move cursor down                                            |
| `Enter`    | Attach to the selected session                              |
| `n`        | Open new session input                                      |
| `d`        | Prompt to delete the selected session                       |
| `x`        | Prompt to kill the selected session                         |
| `D`        | Prompt to delete all sessions                               |
| `r`        | Refresh the session list                                    |
| `q`        | Quit — press **twice** within 2 seconds to confirm         |
| `Ctrl+C`   | Quit immediately                                            |

Any other key cancels a pending quit.

## New session input

Activated by pressing `n` from the session list.

| Key     | Action                                          |
|---------|-------------------------------------------------|
| `Enter` | Create the session (attach or background, per config) |
| `Esc`   | Cancel and return to session list               |

Typing works normally — the input has a 64-character limit.

## Confirmation prompts

Triggered by `d` (delete), `x` (kill), or `D` (delete all).

| Key   | Action                               |
|-------|--------------------------------------|
| `y`   | Confirm the action                   |
| `n`   | Cancel and return to session list    |
| `Esc` | Cancel and return to session list    |

## Behavior notes

- **Double-quit** — pressing `q` once shows "Press q again to quit". A second press within 2 seconds exits. Any other key cancels.
- **Current session** — attempting to attach to the session you're already in shows an error ("Already in this session") instead of crashing.
- **Empty list** — navigation and session-action keys (`Enter`, `d`, `x`) are no-ops when there are no sessions.
