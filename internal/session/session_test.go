package session_test

import (
	"testing"

	"github.com/pageton/zellij-tui/internal/session"
)

func TestSessionFields(t *testing.T) {
	s := session.Session{
		Name:      "main",
		Created:   "1day 21h 5m 55s",
		IsCurrent: true,
		Exited:    false,
	}

	if s.Name != "main" {
		t.Errorf("Name = %q, want %q", s.Name, "main")
	}
	if s.Created != "1day 21h 5m 55s" {
		t.Errorf("Created = %q, want %q", s.Created, "1day 21h 5m 55s")
	}
	if !s.IsCurrent {
		t.Error("IsCurrent should be true")
	}
	if s.Exited {
		t.Error("Exited should be false")
	}
}

func TestValidSessionName(t *testing.T) {
	tests := []struct {
		name  string
		input string
		valid bool
	}{
		{"simple", "main", true},
		{"with hyphen", "my-session", true},
		{"with underscore", "my_session", true},
		{"with dot", "session.v2", true},
		{"numbers", "session123", true},
		{"single char", "a", true},
		{"64 chars", "a123456789012345678901234567890123456789012345678901234567890123", true},
		{"empty", "", false},
		{"leading hyphen", "-flag", false},
		{"double leading hyphen", "--help", false},
		{"spaces", "my session", false},
		{"shell chars", "foo;bar", false},
		{"backtick", "foo`whoami`", false},
		{"dollar", "foo$HOME", false},
		{"pipe", "foo|bar", false},
		{"too long", "a12345678901234567890123456789012345678901234567890123456789012345", false},
		{"slash", "foo/bar", false},
		{"null byte", "foo\x00bar", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := session.ValidSessionName(tt.input)
			if got != tt.valid {
				t.Errorf("ValidSessionName(%q) = %v, want %v", tt.input, got, tt.valid)
			}
		})
	}
}
