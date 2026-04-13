// Package keymap defines keyboard bindings for the TUI.
package keymap

import (
	tea "charm.land/bubbletea/v2"
	uv "github.com/charmbracelet/ultraviolet"
)

// KeyMatch defines a key binding as a set of string patterns (e.g. "up", "k").
// Using the MatchString approach from Bubble Tea v2's KeyPressEvent.
type KeyMatch []string

// Match returns true if the key press matches any of the patterns.
func (km KeyMatch) Match(key tea.KeyPressMsg) bool {
	return uv.Key(key.Key()).MatchString(km...)
}

var (
	// Up moves the cursor up.
	Up = KeyMatch{"up", "k"}
	// Down moves the cursor down.
	Down = KeyMatch{"down", "j"}
	// Enter confirms selection.
	Enter = KeyMatch{"enter"}
	// NewSession opens the session creation input.
	NewSession = KeyMatch{"n"}
	// Delete prompts to delete the selected session.
	Delete = KeyMatch{"d"}
	// Kill prompts to kill the selected session.
	Kill = KeyMatch{"x"}
	// DeleteAll prompts to delete all sessions.
	DeleteAll = KeyMatch{"D"}
	// Refresh reloads the session list.
	Refresh = KeyMatch{"r"}
	// Quit exits the TUI.
	Quit = KeyMatch{"q"}
	// Escape cancels the current operation.
	Escape = KeyMatch{"esc"}
	// ConfirmYes accepts a confirmation prompt.
	ConfirmYes = KeyMatch{"y"}
	// ConfirmNo rejects a confirmation prompt.
	ConfirmNo = KeyMatch{"n"} // same key as NewSession; safe because confirm mode intercepts first
	// CtrlC forces an immediate quit.
	CtrlC = KeyMatch{"ctrl+c"}
)
