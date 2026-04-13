// Package session defines session types and validation for zellij-tui.
package session

import "regexp"

// Session represents a single Zellij session.
type Session struct {
	Name      string // e.g. "main"
	Created   string // e.g. "1day 21h 5m 55s"
	IsCurrent bool   // marked as "(current)" by zellij
	Exited    bool   // session has exited, can be resurrected
}

// sessionNameRe defines the allowlist for valid session names.
// Matches: alphanumeric, underscore, hyphen, dot — 1 to 64 chars.
// Rejects names that could be interpreted as CLI flags (leading hyphens).
var sessionNameRe = regexp.MustCompile(`^[a-zA-Z0-9_][a-zA-Z0-9_.-]{0,63}$`)

// ValidSessionName returns true if the name matches the allowlist.
func ValidSessionName(name string) bool {
	return sessionNameRe.MatchString(name)
}
