package zellij

import (
	"testing"
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
