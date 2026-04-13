package zellij

import (
	"os"
	"os/exec"
	"regexp"
	"strings"
	"syscall"

	"github.com/sadiq/zellij-tui/internal/types"
)

// sessionLineRe matches lines like:
//   main [Created 1day 21h 5m 55s ago] (current)
//   work [Created 3h 12m ago]
//   hello [Created 10m ago] (EXITED - attach to resurrect)
var sessionLineRe = regexp.MustCompile(`^(\S+)\s+\[Created\s+(.+?)\s+ago\](?:\s+\(EXITED[^)]*\))?(\s+\(current\))?$`)

// binary is the path to the zellij executable.
var binary = "zellij"

// SetBinary sets the path to the zellij binary.
func SetBinary(path string) {
	binary = path
}

// ListSessions runs `zellij list-sessions --no-formatting` and parses the output.
// Returns an empty slice (no error) if there are no sessions.
func ListSessions() ([]types.Session, error) {
	cmd := exec.Command(binary, "list-sessions", "--no-formatting")
	out, err := cmd.Output()
	if err != nil {
		if cmd.ProcessState != nil && cmd.ProcessState.ExitCode() == 1 {
			return nil, nil
		}
		return nil, err
	}

	var sessions []types.Session
	for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		matches := sessionLineRe.FindStringSubmatch(line)
		if matches == nil {
			continue
		}
		sessions = append(sessions, types.Session{
			Name:      matches[1],
			Created:   matches[2],
			IsCurrent: matches[3] != "",
			Exited:    strings.Contains(line, "(EXITED"),
		})
	}

	return sessions, nil
}

// DeleteSession kills and deletes a session by name.
func DeleteSession(name string) error {
	return exec.Command(binary, "delete-session", "-f", name).Run()
}

// CreateSession creates a new session fully detached from this process.
// It spawns zellij in a new session group so it survives terminal close.
func CreateSession(name string) error {
	cmd := exec.Command(binary, "attach", "-b", "-c", name)
	cmd.Stdin = nil
	cmd.Stdout = nil
	cmd.Stderr = nil
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setsid: true,
	}
	return cmd.Start()
}

// KillSession kills a running session by name (without deleting).
func KillSession(name string) error {
	return exec.Command(binary, "kill-session", name).Run()
}

// DeleteAllSessions kills and deletes all sessions.
func DeleteAllSessions() error {
	return exec.Command(binary, "delete-all-sessions").Run()
}

// IsInsideSession checks if we're currently running inside a Zellij session.
func IsInsideSession() bool {
	for _, key := range []string{"ZELLIJ_SESSION_NAME", "ZELLIJ"} {
		if os.Getenv(key) != "" {
			return true
		}
	}
	return false
}
