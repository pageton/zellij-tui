// Package ui provides styles and rendering for the zellij-tui interface.
package ui

import (
	lipgloss "charm.land/lipgloss/v2"
)

var (
	// ContainerStyle wraps everything in a rounded border box.
	ContainerStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#3C3C3C")).
			Padding(2, 4)

	// TitleStyle for the heading.
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFFFFF"))

	// CursorStyle for the cursor indicator "▸".
	CursorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6CDFFF"))

	// SelectedNameStyle for the session name when selected.
	SelectedNameStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#FFFFFF"))

	// UnselectedNameStyle for session names that aren't selected.
	UnselectedNameStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#6C6C6C"))

	// TimeStyle for the "created X ago" text.
	TimeStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#4A4A4A"))

	// CurrentDotStyle for the "●" indicator on the current session.
	CurrentDotStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#5FFF87"))

	// StatusBarStyle for the bottom key hints.
	StatusBarStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#5A5A5A"))

	// KeyStyle for key hints in the status bar.
	KeyStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#6CDFFF"))

	// EmptyMessageStyle for "No active sessions found".
	EmptyMessageStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#6C6C6C"))

	// InputBoxStyle for the text input.
	InputBoxStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#AAAAAA"))

	// ErrorStyle for inline error messages.
	ErrorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF6B6B"))

	// WarningStyle for the "already in a session" banner.
	WarningStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFD700"))

	// QuitHintStyle for the "Press q again to quit" message.
	QuitHintStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#5A5A5A"))

	// SeparatorStyle for the dashed line above the status bar.
	SeparatorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#3C3C3C"))

	// ConfirmPromptStyle for delete confirmation text.
	ConfirmPromptStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FF6B6B"))
)

// FmtKey renders a key label in the accent color.
func FmtKey(key string) string {
	return KeyStyle.Render(key)
}
