package ui

import (
	"strings"
	"testing"

	"github.com/pageton/zellij-tui/internal/session"
)

func TestRenderEmptyState(t *testing.T) {
	output := Render(ViewData{State: StateList, TermWidth: 80})
	if output == "" {
		t.Fatal("Render returned empty string")
	}
	if !strings.Contains(output, "No active sessions found") {
		t.Error("empty state should mention no sessions")
	}
}

func TestRenderWithSessions(t *testing.T) {
	sessions := []session.Session{
		{Name: "main", Created: "3h 12m", IsCurrent: true},
		{Name: "work", Created: "5s", IsCurrent: false},
	}

	output := Render(ViewData{
		Sessions:  sessions,
		Cursor:    0,
		State:     StateList,
		TermWidth: 80,
	})
	if output == "" {
		t.Fatal("Render returned empty string")
	}
	if !strings.Contains(output, "main") {
		t.Error("output should contain session name 'main'")
	}
	if !strings.Contains(output, "work") {
		t.Error("output should contain session name 'work'")
	}
}

func TestFmtKey(t *testing.T) {
	result := FmtKey("enter")
	if result == "enter" {
		t.Error("FmtKey should apply styling, output should differ from input")
	}
}
