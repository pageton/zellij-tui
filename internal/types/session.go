package types

// Session represents a single Zellij session.
type Session struct {
	Name      string // e.g. "main"
	Created   string // e.g. "1day 21h 5m 55s"
	IsCurrent bool   // marked as "(current)" by zellij
	Exited    bool   // session has exited, can be resurrected
}
