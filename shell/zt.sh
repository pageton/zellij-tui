#!/usr/bin/env bash
# zt — Zellij session manager TUI
#
# Usage:
#   Source this file in your shell config (.bashrc / .zshrc):
#     source /path/to/zellij-tui/shell/zt.sh
#
#   Then run:
#     zt
#
# The binary exec's directly into zellij, so no wrapper logic is needed.
# This function just ensures a clean shell context.
#
# Optional auto-launch on terminal open (add to .bashrc / .zshrc):
#   if [[ -z "$ZELLIJ" ]] && [[ $- == *i* ]]; then
#     zt
#   fi
#
# Config file: ~/.config/zellij-tui/config.toml
#   zellij_path = "/usr/bin/zellij"   # optional, defaults to PATH lookup

zt() {
    zellij-tui "$@"
}
