package zellij

import (
	"strings"
	"testing"

	"github.com/pageton/zellij-tui/internal/session"
)

func TestSessionLineRe(t *testing.T) {
	tests := []struct {
		line        string
		name        string
		created     string
		isCurrent   bool
		exited      bool
		shouldMatch bool
	}{
		{
			line:        "main [Created 1day 21h 5m 55s ago] (current)",
			name:        "main",
			created:     "1day 21h 5m 55s",
			isCurrent:   true,
			shouldMatch: true,
		},
		{
			line:        "work [Created 3h 12m ago]",
			name:        "work",
			created:     "3h 12m",
			shouldMatch: true,
		},
		{
			line:        "scratch [Created 12m ago] (current)",
			name:        "scratch",
			created:     "12m",
			isCurrent:   true,
			shouldMatch: true,
		},
		{
			line:        "my-session [Created 5s ago]",
			name:        "my-session",
			created:     "5s",
			shouldMatch: true,
		},
		{
			line:        "hello [Created 10m 50s ago] (EXITED - attach to resurrect)",
			name:        "hello",
			created:     "10m 50s",
			exited:      true,
			shouldMatch: true,
		},
		{
			line:        "",
			shouldMatch: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.line, func(t *testing.T) {
			matches := sessionLineRe.FindStringSubmatch(tt.line)
			if !tt.shouldMatch {
				if matches != nil {
					t.Errorf("expected no match for %q, got %v", tt.line, matches)
				}
				return
			}
			if matches == nil {
				t.Fatalf("expected match for %q, got none", tt.line)
			}
			if matches[1] != tt.name {
				t.Errorf("name: got %q, want %q", matches[1], tt.name)
			}
			if matches[2] != tt.created {
				t.Errorf("created: got %q, want %q", matches[2], tt.created)
			}
			gotCurrent := matches[3] != ""
			if gotCurrent != tt.isCurrent {
				t.Errorf("isCurrent: got %v, want %v", gotCurrent, tt.isCurrent)
			}
		})
	}
}

func TestParseSessionsRejectsInvalidNames(t *testing.T) {
	tests := []struct {
		name      string
		line      string
		wantValid bool
	}{
		{"valid name", "my-session [Created 5s ago]", true},
		{"leading hyphen", "--help [Created 5s ago]", false},
		{"spaces in name", "bad session [Created 5s ago]", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Simulate the parse+validate pipeline from ListSessions
			for _, line := range strings.Split(tt.line, "\n") {
				line = strings.TrimSpace(line)
				if line == "" {
					continue
				}
				matches := sessionLineRe.FindStringSubmatch(line)
				if matches == nil {
					if tt.wantValid {
						t.Errorf("regex should match %q", line)
					}
					continue
				}
				valid := session.ValidSessionName(matches[1])
				if valid != tt.wantValid {
					t.Errorf("ValidSessionName(%q) = %v, want %v", matches[1], valid, tt.wantValid)
				}
			}
		})
	}
}
