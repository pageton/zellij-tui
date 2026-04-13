// Package main is the entrypoint for the zellij-tui TUI session manager.
package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"

	tea "charm.land/bubbletea/v2"

	"github.com/pageton/zellij-tui/internal/config"
	"github.com/pageton/zellij-tui/internal/model"
	"github.com/pageton/zellij-tui/internal/session"
	"github.com/pageton/zellij-tui/internal/zellij"
)

// version is set at build time via -ldflags="-X main.version=...".
var version = "dev"

func main() {
	os.Exit(run())
}

func run() int {
	if len(os.Args) > 1 && os.Args[1] == "--version" {
		fmt.Printf("zellij-tui %s\n", version)
		return 0
	}

	cfg := config.Load()

	// Resolve zellij binary path.
	// TRUST BOUNDARY: if zellij_path is set in config, it points to an
	// arbitrary binary that is executed with full user privileges. The
	// --version check below is a heuristic, not a cryptographic guarantee.
	// Users must ensure the config file and the target binary are not
	// writable by untrusted parties.
	zellijPath := cfg.ZellijPath
	if zellijPath == "" {
		found, err := exec.LookPath("zellij")
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error: zellij not found in PATH.")
			fmt.Fprintln(os.Stderr, "Install from https://zellij.dev or set zellij_path in config.")
			return 1
		}
		zellijPath = found
	}

	// Verify the binary actually exists
	if _, err := os.Stat(zellijPath); err != nil {
		fmt.Fprintf(os.Stderr, "Error: zellij binary not found at %s\n", zellijPath)
		return 1
	}

	// Verify the binary is actually zellij
	out, err := exec.Command(zellijPath, "--version").Output()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to run %s --version: %v\n", zellijPath, err)
		return 1
	}
	if !strings.Contains(strings.ToLower(string(out)), "zellij") {
		fmt.Fprintf(os.Stderr, "Error: %s does not appear to be zellij (version output: %s)\n", zellijPath, strings.TrimSpace(string(out)))
		return 1
	}

	// Tell the client to use this binary path
	zellij.SetBinary(zellijPath)

	m := model.New()
	m.SetInsideSession(zellij.IsInsideSession())
	m.SetAutoAttach(cfg.AutoAttach)

	// Render TUI to /dev/tty so it works even when stdout is captured
	tty, err := os.OpenFile("/dev/tty", os.O_RDWR, 0)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error: cannot open /dev/tty")
		return 1
	}
	defer func() { _ = tty.Close() }()

	p := tea.NewProgram(m,
		tea.WithInput(tty),
		tea.WithOutput(tty),
	)

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return 1
	}

	// Exec into zellij directly — replaces this process
	result := m.Result()

	// Defense-in-depth: re-validate session name at the exec boundary.
	// This guards against future regressions in the validation pipeline.
	if result.Name != "" && !session.ValidSessionName(result.Name) {
		fmt.Fprintf(os.Stderr, "Error: invalid session name %q\n", result.Name)
		return 1
	}

	switch result.Op {
	case model.ActionAttach:
		_ = syscall.Exec(zellijPath, []string{"zellij", "attach", result.Name}, os.Environ())
	case model.ActionCreate:
		// Create the session in the background first. This starts a
		// persistent zellij server that survives terminal closure.
		// Then exec into a client-only attach (not -c) so the session
		// persists even if the terminal is closed.
		if err := zellij.CreateSession(result.Name); err != nil {
			fmt.Fprintf(os.Stderr, "Error: failed to create session %q: %v\n", result.Name, err)
			return 1
		}
		_ = syscall.Exec(zellijPath, []string{"zellij", "attach", result.Name}, os.Environ())
	}
	return 0
}
