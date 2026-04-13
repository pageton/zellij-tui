package ui

import (
	"fmt"
	"strings"

	"github.com/clipperhouse/displaywidth"

	"github.com/pageton/zellij-tui/internal/session"
)

// State represents the current UI mode.
type State int

// State values represent the current UI mode.
const (
	StateList State = iota
	StateInput
	StateConfirmDelete
	StateConfirmKill
	StateConfirmDeleteAll
)

// cardWidth is the content width inside the border (border + padding add more).
const cardWidth = 68
const cardMinHeight = 12
const cardMaxHeight = 24

// ViewData carries all model state needed by the renderer.
type ViewData struct {
	Sessions      []session.Session
	Cursor        int
	State         State
	InputView     string
	ErrMsg        string
	QuitPending   bool
	InsideSession bool
	TermWidth     int
}

// Render builds the full view string from the model state.
func Render(d ViewData) string {

	var b strings.Builder

	// Warning banner when already inside a session
	if d.InsideSession {
		b.WriteString(WarningStyle.Render("⚠ Already inside a Zellij session"))
		b.WriteString("\n\n")
	}

	// Title
	b.WriteString(TitleStyle.Render("◈ Zellij Sessions"))
	b.WriteString("\n\n")

	// Session list or empty state
	if len(d.Sessions) == 0 {
		renderEmpty(&b)
	} else {
		renderSessionList(&b, d.Sessions, d.Cursor)
	}

	// Overlays
	switch d.State {
	case StateInput:
		b.WriteString("\n")
		b.WriteString(d.InputView)
		b.WriteString("\n")
		b.WriteString(StatusBarStyle.Render(fmt.Sprintf(
			"%s create · %s cancel", FmtKey("enter"), FmtKey("esc"),
		)))

	case StateConfirmDelete:
		if d.Cursor >= 0 && d.Cursor < len(d.Sessions) {
			b.WriteString("\n")
			b.WriteString(ConfirmPromptStyle.Render(fmt.Sprintf(
				"Delete %q? %s/%s",
				d.Sessions[d.Cursor].Name, FmtKey("y"), FmtKey("n"),
			)))
		}

	case StateConfirmKill:
		if d.Cursor >= 0 && d.Cursor < len(d.Sessions) {
			b.WriteString("\n")
			b.WriteString(ConfirmPromptStyle.Render(fmt.Sprintf(
				"Kill %q? %s/%s",
				d.Sessions[d.Cursor].Name, FmtKey("y"), FmtKey("n"),
			)))
		}

	case StateConfirmDeleteAll:
		b.WriteString("\n")
		b.WriteString(ConfirmPromptStyle.Render(fmt.Sprintf(
			"Delete ALL sessions? %s/%s", FmtKey("y"), FmtKey("n"),
		)))

	case StateList:
		if d.ErrMsg != "" {
			b.WriteString("\n")
			b.WriteString(ErrorStyle.Render(d.ErrMsg))
		}
		if d.QuitPending {
			b.WriteString("\n")
			b.WriteString(QuitHintStyle.Render("Press q again to quit"))
		}
	}

	// Status bar — always show except in input mode
	if d.State != StateInput {
		b.WriteString("\n")
		b.WriteString(SeparatorStyle.Render(strings.Repeat("─", cardWidth-10)))
		b.WriteString("\n")
		renderStatusBar(&b)
	}

	// Wrap in bordered container
	content := b.String()
	// Adaptive height: grows with content, min 12, max 24
	h := cardMinHeight
	// Rough estimate: title(3) + sessions(N) + overlays(2) + separator(2) + status(1) + padding(4)
	needed := 12 + len(d.Sessions)
	if d.State == StateInput {
		needed += 3
	}
	if needed > h {
		h = needed
	}
	if h > cardMaxHeight {
		h = cardMaxHeight
	}

	return ContainerStyle.Width(cardWidth).Height(h).Render(content)
}

func renderEmpty(b *strings.Builder) {
	b.WriteString(EmptyMessageStyle.Render("No active sessions found."))
	b.WriteString("\n\n")
	fmt.Fprintf(b, "Press %s to create a new session.", FmtKey("n"))
	b.WriteString("\n")
}

func renderSessionList(b *strings.Builder, sessions []session.Session, cursor int) {
	maxName := 0
	for _, s := range sessions {
		w := displaywidth.String(s.Name)
		if w > maxName {
			maxName = w
		}
	}
	if maxName < 8 {
		maxName = 8
	}

	for i, s := range sessions {
		var cursorChar string
		if i == cursor {
			cursorChar = CursorStyle.Render("▸ ")
		} else {
			cursorChar = "  "
		}

		var name string
		if i == cursor {
			name = SelectedNameStyle.Render(s.Name)
		} else {
			name = UnselectedNameStyle.Render(s.Name)
		}

		namePad := maxName - displaywidth.String(s.Name)
		padded := name + strings.Repeat(" ", namePad)

		timeStr := TimeStyle.Render(s.Created + " ago")

		b.WriteString(cursorChar)
		b.WriteString(padded)
		b.WriteString("  ")
		b.WriteString(timeStr)

		if s.IsCurrent {
			b.WriteString("  ")
			b.WriteString(CurrentDotStyle.Render("●"))
		}
		if s.Exited {
			b.WriteString("  ")
			b.WriteString(TimeStyle.Render("exited"))
		}

		b.WriteString("\n")
	}
}

func renderStatusBar(b *strings.Builder) {
	parts := []string{
		FmtKey("n") + " new",
		FmtKey("d") + " del",
		FmtKey("x") + " kill",
		FmtKey("D") + " del all",
		"⏎ open",
		FmtKey("r") + " refresh",
		FmtKey("q") + " quit",
	}
	b.WriteString(StatusBarStyle.Render(strings.Join(parts, " · ")))
}
