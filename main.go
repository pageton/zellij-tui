package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"

	tea "charm.land/bubbletea/v2"

	"github.com/sadiq/zellij-tui/internal/config"
	"github.com/sadiq/zellij-tui/internal/model"
	"github.com/sadiq/zellij-tui/internal/zellij"
)

func main() {
	cfg := config.Load()

	// Resolve zellij binary path
	zellijPath := cfg.ZellijPath
	if zellijPath == "" {
		found, err := exec.LookPath("zellij")
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error: zellij not found in PATH.")
			fmt.Fprintln(os.Stderr, "Install from https://zellij.dev or set zellij_path in config.")
			os.Exit(1)
		}
		zellijPath = found
	}

	// Verify the binary actually exists
	if _, err := os.Stat(zellijPath); err != nil {
		fmt.Fprintf(os.Stderr, "Error: zellij binary not found at %s\n", zellijPath)
		os.Exit(1)
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
		os.Exit(1)
	}
	defer tty.Close()

	p := tea.NewProgram(m,
		tea.WithInput(tty),
		tea.WithOutput(tty),
	)

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Exec into zellij directly — replaces this process
	result := m.Result()
	switch result.Op {
	case model.ActionAttach:
		syscall.Exec(zellijPath, []string{"zellij", "attach", result.Name}, os.Environ())
	case model.ActionCreate:
		syscall.Exec(zellijPath, []string{"zellij", "attach", "-c", result.Name}, os.Environ())
	}
}
