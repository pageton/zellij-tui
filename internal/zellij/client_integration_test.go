package zellij

import (
	"os"
	"path/filepath"
	"testing"
)

// writeFakeBinary creates an executable script that prints the given output
// and sets SetBinary to point at it. Returns a cleanup function.
func writeFakeBinary(t *testing.T, output string) func() {
	t.Helper()
	dir := t.TempDir()
	script := filepath.Join(dir, "zellij")
	content := "#!/bin/sh\n" + output
	if err := os.WriteFile(script, []byte(content), 0o755); err != nil {
		t.Fatal(err)
	}
	origBinary := binary
	SetBinary(script)
	return func() {
		binary = origBinary
	}
}

func TestListSessionsParsesOutput(t *testing.T) {
	cleanup := writeFakeBinary(t, `echo 'main [Created 1day 21h 5m 55s ago] (current)'
echo 'work [Created 3h 12m ago]'
`)
	defer cleanup()

	sessions, err := ListSessions()
	if err != nil {
		t.Fatalf("ListSessions error: %v", err)
	}
	if len(sessions) != 2 {
		t.Fatalf("len(sessions) = %d, want 2", len(sessions))
	}
	if sessions[0].Name != "main" {
		t.Errorf("sessions[0].Name = %q, want %q", sessions[0].Name, "main")
	}
	if !sessions[0].IsCurrent {
		t.Error("sessions[0].IsCurrent should be true")
	}
	if sessions[1].Name != "work" {
		t.Errorf("sessions[1].Name = %q, want %q", sessions[1].Name, "work")
	}
	if sessions[1].IsCurrent {
		t.Error("sessions[1].IsCurrent should be false")
	}
}

func TestListSessionsEmpty(t *testing.T) {
	// Exit code 1 means no sessions (zellij behavior)
	cleanup := writeFakeBinary(t, "exit 1\n")
	defer cleanup()

	sessions, err := ListSessions()
	if err != nil {
		t.Fatalf("ListSessions error: %v", err)
	}
	if sessions != nil {
		t.Fatalf("expected nil slice for no sessions, got %v", sessions)
	}
}

func TestListSessionsExited(t *testing.T) {
	cleanup := writeFakeBinary(t, `echo 'hello [Created 10m ago] (EXITED - attach to resurrect)'
`)
	defer cleanup()

	sessions, err := ListSessions()
	if err != nil {
		t.Fatalf("ListSessions error: %v", err)
	}
	if len(sessions) != 1 {
		t.Fatalf("len(sessions) = %d, want 1", len(sessions))
	}
	if !sessions[0].Exited {
		t.Error("Exited should be true")
	}
}

func TestListSessionsRejectsInvalidNames(t *testing.T) {
	cleanup := writeFakeBinary(t, `echo '--help [Created 5s ago]'
echo 'valid-session [Created 1h ago]'
`)
	defer cleanup()

	sessions, err := ListSessions()
	if err != nil {
		t.Fatalf("ListSessions error: %v", err)
	}
	if len(sessions) != 1 {
		t.Fatalf("len(sessions) = %d, want 1 (invalid name filtered)", len(sessions))
	}
	if sessions[0].Name != "valid-session" {
		t.Errorf("sessions[0].Name = %q, want %q", sessions[0].Name, "valid-session")
	}
}

func TestListSessionsSkipsUnparseableLines(t *testing.T) {
	cleanup := writeFakeBinary(t, `echo 'garbage line without brackets'
echo 'work [Created 3h 12m ago]'
`)
	defer cleanup()

	sessions, err := ListSessions()
	if err != nil {
		t.Fatalf("ListSessions error: %v", err)
	}
	if len(sessions) != 1 {
		t.Fatalf("len(sessions) = %d, want 1 (unparseable line skipped)", len(sessions))
	}
}

func TestDeleteSessionCallsBinary(t *testing.T) {
	called := filepath.Join(t.TempDir(), "called")
	dir := t.TempDir()
	script := filepath.Join(dir, "zellij")
	content := "#!/bin/sh\necho \"$@\" > " + called + "\nexit 0\n"
	if err := os.WriteFile(script, []byte(content), 0o755); err != nil {
		t.Fatal(err)
	}
	origBinary := binary
	SetBinary(script)
	defer func() { binary = origBinary }()

	err := DeleteSession("test-session")
	if err != nil {
		t.Fatalf("DeleteSession error: %v", err)
	}
	data, err := os.ReadFile(called)
	if err != nil {
		t.Fatal(err)
	}
	got := string(data)
	want := "delete-session -f test-session\n"
	if got != want {
		t.Errorf("binary called with %q, want %q", got, want)
	}
}

func TestKillSessionCallsBinary(t *testing.T) {
	called := filepath.Join(t.TempDir(), "called")
	dir := t.TempDir()
	script := filepath.Join(dir, "zellij")
	content := "#!/bin/sh\necho \"$@\" > " + called + "\nexit 0\n"
	if err := os.WriteFile(script, []byte(content), 0o755); err != nil {
		t.Fatal(err)
	}
	origBinary := binary
	SetBinary(script)
	defer func() { binary = origBinary }()

	err := KillSession("my-session")
	if err != nil {
		t.Fatalf("KillSession error: %v", err)
	}
	data, err := os.ReadFile(called)
	if err != nil {
		t.Fatal(err)
	}
	got := string(data)
	want := "kill-session my-session\n"
	if got != want {
		t.Errorf("binary called with %q, want %q", got, want)
	}
}

func TestCreateSessionCallsBinary(t *testing.T) {
	called := filepath.Join(t.TempDir(), "called")
	dir := t.TempDir()
	script := filepath.Join(dir, "zellij")
	content := "#!/bin/sh\necho \"$@\" > " + called + "\nexit 0\n"
	if err := os.WriteFile(script, []byte(content), 0o755); err != nil {
		t.Fatal(err)
	}
	origBinary := binary
	SetBinary(script)
	defer func() { binary = origBinary }()

	err := CreateSession("new-session")
	if err != nil {
		t.Fatalf("CreateSession error: %v", err)
	}
	data, err := os.ReadFile(called)
	if err != nil {
		t.Fatal(err)
	}
	got := string(data)
	want := "attach -b new-session\n"
	if got != want {
		t.Errorf("binary called with %q, want %q", got, want)
	}
}

func TestDeleteAllSessionsCallsBinary(t *testing.T) {
	called := filepath.Join(t.TempDir(), "called")
	dir := t.TempDir()
	script := filepath.Join(dir, "zellij")
	content := "#!/bin/sh\necho \"$@\" > " + called + "\nexit 0\n"
	if err := os.WriteFile(script, []byte(content), 0o755); err != nil {
		t.Fatal(err)
	}
	origBinary := binary
	SetBinary(script)
	defer func() { binary = origBinary }()

	err := DeleteAllSessions()
	if err != nil {
		t.Fatalf("DeleteAllSessions error: %v", err)
	}
	data, err := os.ReadFile(called)
	if err != nil {
		t.Fatal(err)
	}
	got := string(data)
	want := "delete-all-sessions\n"
	if got != want {
		t.Errorf("binary called with %q, want %q", got, want)
	}
}
