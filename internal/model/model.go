package model

import (
	"fmt"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/bubbles/v2/textinput"
	lipgloss "charm.land/lipgloss/v2"
	uv "github.com/charmbracelet/ultraviolet"

	"github.com/sadiq/zellij-tui/internal/keymap"
	"github.com/sadiq/zellij-tui/internal/types"
	"github.com/sadiq/zellij-tui/internal/ui"
	"github.com/sadiq/zellij-tui/internal/zellij"
)

const quitTimeout = 2 * time.Second

// matchKey checks if a tea.KeyPressMsg matches any of the given key patterns.
func matchKey(key tea.KeyPressMsg, patterns []string) bool {
	return uv.Key(key.Key()).MatchString(patterns...)
}

// tickMsg is sent when the quit timer expires.
type tickMsg time.Time

// sessionsLoadedMsg is sent when session listing completes.
type sessionsLoadedMsg struct {
	sessions []types.Session
	err      error
}

// deleteResultMsg is sent when a session deletion/kill completes.
type deleteResultMsg struct {
	err error
}

// createBackgroundResultMsg is sent when a background session creation completes.
type createBackgroundResultMsg struct {
	name string
	err  error
}

// TUI is the Bubble Tea model for the session manager.
type TUI struct {
	state         ui.State
	sessions      []types.Session
	cursor        int
	input         textinput.Model
	err           string
	quitPending   bool
	insideSession bool
	autoAttach    bool
	width         int
	height        int
	result        Action
}

// New creates a new TUI model.
func New() *TUI {
	ti := textinput.New()
	ti.Prompt = "Session name: "
	ti.Placeholder = "my-session"
	ti.CharLimit = 64

	return &TUI{
		state:      ui.StateList,
		input:      ti,
		autoAttach: true,
	}
}

func (m *TUI) Init() tea.Cmd {
	return tea.Batch(
		func() tea.Msg {
			sessions, err := zellij.ListSessions()
			return sessionsLoadedMsg{sessions: sessions, err: err}
		},
		tea.RequestWindowSize,
	)
}

func (m *TUI) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case sessionsLoadedMsg:
		if msg.err != nil {
			m.err = fmt.Sprintf("Failed to list sessions: %v", msg.err)
		} else {
			m.sessions = msg.sessions
			m.err = ""
		}
		if m.cursor >= len(m.sessions) {
			m.cursor = len(m.sessions) - 1
		}
		if m.cursor < 0 {
			m.cursor = 0
		}
		return m, nil

	case deleteResultMsg:
		if msg.err != nil {
			m.err = fmt.Sprintf("Failed to delete session: %v", msg.err)
		} else {
			m.err = ""
			return m, func() tea.Msg {
				sessions, err := zellij.ListSessions()
				return sessionsLoadedMsg{sessions: sessions, err: err}
			}
		}
		return m, nil

	case createBackgroundResultMsg:
		if msg.err != nil {
			m.err = fmt.Sprintf("Failed to create session: %v", msg.err)
		} else {
			m.err = ""
			return m, func() tea.Msg {
				sessions, err := zellij.ListSessions()
				return sessionsLoadedMsg{sessions: sessions, err: err}
			}
		}
		return m, nil

	case tickMsg:
		if m.quitPending {
			m.quitPending = false
		}
		return m, nil
	}

	switch m.state {
	case ui.StateList:
		return m.updateList(msg)
	case ui.StateInput:
		return m.updateInput(msg)
	case ui.StateConfirmDelete:
		return m.updateConfirmDelete(msg)
	case ui.StateConfirmKill:
		return m.updateConfirmKill(msg)
	case ui.StateConfirmDeleteAll:
		return m.updateConfirmDeleteAll(msg)
	}

	return m, nil
}

func (m *TUI) updateList(msg tea.Msg) (tea.Model, tea.Cmd) {
	key, ok := msg.(tea.KeyPressMsg)
	if !ok {
		return m, nil
	}

	switch {
	case matchKey(key, keymap.Up):
		if m.cursor > 0 {
			m.cursor--
		}
		m.quitPending = false
		m.err = ""
		return m, nil

	case matchKey(key, keymap.Down):
		if m.cursor < len(m.sessions)-1 {
			m.cursor++
		}
		m.quitPending = false
		m.err = ""
		return m, nil

	case matchKey(key, keymap.Enter):
		if len(m.sessions) == 0 || m.cursor < 0 || m.cursor >= len(m.sessions) {
			return m, nil
		}
		// Prevent attaching to the current session — zellij crashes
		if m.sessions[m.cursor].IsCurrent {
			m.err = "Already in this session"
			return m, nil
		}
		m.result = Action{Op: ActionAttach, Name: m.sessions[m.cursor].Name}
		return m, tea.Quit

	case matchKey(key, keymap.NewSession):
		m.state = ui.StateInput
		m.input.SetValue("")
		cmd := m.input.Focus()
		m.quitPending = false
		m.err = ""
		return m, cmd

	case matchKey(key, keymap.Delete):
		if len(m.sessions) == 0 || m.cursor < 0 || m.cursor >= len(m.sessions) {
			return m, nil
		}
		m.state = ui.StateConfirmDelete
		m.quitPending = false
		m.err = ""
		return m, nil

	case matchKey(key, keymap.Kill):
		if len(m.sessions) == 0 || m.cursor < 0 || m.cursor >= len(m.sessions) {
			return m, nil
		}
		m.state = ui.StateConfirmKill
		m.quitPending = false
		m.err = ""
		return m, nil

	case matchKey(key, keymap.DeleteAll):
		if len(m.sessions) == 0 {
			return m, nil
		}
		m.state = ui.StateConfirmDeleteAll
		m.quitPending = false
		m.err = ""
		return m, nil

	case matchKey(key, keymap.Refresh):
		m.quitPending = false
		m.err = ""
		return m, func() tea.Msg {
			sessions, err := zellij.ListSessions()
			return sessionsLoadedMsg{sessions: sessions, err: err}
		}

	case matchKey(key, keymap.Quit):
		if m.quitPending {
			m.result = Action{Op: ActionNone}
			return m, tea.Quit
		}
		m.quitPending = true
		return m, tea.Tick(quitTimeout, func(t time.Time) tea.Msg {
			return tickMsg(t)
		})

	case matchKey(key, keymap.CtrlC):
		m.result = Action{Op: ActionNone}
		return m, tea.Quit

	default:
		if m.quitPending {
			m.quitPending = false
		}
	}

	return m, nil
}

func (m *TUI) updateInput(msg tea.Msg) (tea.Model, tea.Cmd) {
	if key, ok := msg.(tea.KeyPressMsg); ok {
		switch {
		case matchKey(key, keymap.Escape):
			m.state = ui.StateList
			m.input.Blur()
			return m, nil

		case matchKey(key, keymap.Enter):
			name := m.input.Value()
			if name == "" {
				return m, nil
			}
			if m.autoAttach {
				m.result = Action{Op: ActionCreate, Name: name}
				return m, tea.Quit
			}
			// Create in background, stay in TUI
			m.state = ui.StateList
			m.input.Blur()
			return m, func() tea.Msg {
				err := zellij.CreateSession(name)
				return createBackgroundResultMsg{name: name, err: err}
			}
		}
	}

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m *TUI) updateConfirmDelete(msg tea.Msg) (tea.Model, tea.Cmd) {
	key, ok := msg.(tea.KeyPressMsg)
	if !ok {
		return m, nil
	}

	switch {
	case matchKey(key, keymap.ConfirmYes):
		if m.cursor < 0 || m.cursor >= len(m.sessions) {
			m.state = ui.StateList
			return m, nil
		}
		name := m.sessions[m.cursor].Name
		m.state = ui.StateList
		return m, func() tea.Msg {
			err := zellij.DeleteSession(name)
			return deleteResultMsg{err: err}
		}

	case matchKey(key, keymap.ConfirmNo), matchKey(key, keymap.Escape):
		m.state = ui.StateList
		return m, nil
	}

	return m, nil
}

func (m *TUI) updateConfirmKill(msg tea.Msg) (tea.Model, tea.Cmd) {
	key, ok := msg.(tea.KeyPressMsg)
	if !ok {
		return m, nil
	}

	switch {
	case matchKey(key, keymap.ConfirmYes):
		if m.cursor < 0 || m.cursor >= len(m.sessions) {
			m.state = ui.StateList
			return m, nil
		}
		name := m.sessions[m.cursor].Name
		m.state = ui.StateList
		return m, func() tea.Msg {
			err := zellij.KillSession(name)
			return deleteResultMsg{err: err}
		}

	case matchKey(key, keymap.ConfirmNo), matchKey(key, keymap.Escape):
		m.state = ui.StateList
		return m, nil
	}

	return m, nil
}

func (m *TUI) updateConfirmDeleteAll(msg tea.Msg) (tea.Model, tea.Cmd) {
	key, ok := msg.(tea.KeyPressMsg)
	if !ok {
		return m, nil
	}

	switch {
	case matchKey(key, keymap.ConfirmYes):
		m.state = ui.StateList
		return m, func() tea.Msg {
			err := zellij.DeleteAllSessions()
			return deleteResultMsg{err: err}
		}

	case matchKey(key, keymap.ConfirmNo), matchKey(key, keymap.Escape):
		m.state = ui.StateList
		return m, nil
	}

	return m, nil
}

func (m *TUI) View() tea.View {
	uiSessions := make([]ui.SessionData, len(m.sessions))
	for i, s := range m.sessions {
		uiSessions[i] = ui.SessionData{
			Name:      s.Name,
			Created:   s.Created,
			IsCurrent: s.IsCurrent,
			Exited:    s.Exited,
		}
	}

	content := ui.Render(
		uiSessions,
		m.cursor,
		m.state,
		m.input.View(),
		m.err,
		m.quitPending,
		m.insideSession,
		m.width,
	)

	// Center the content on screen
	if m.width > 0 && m.height > 0 {
		content = lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, content)
	}

	var view tea.View
	view.AltScreen = true
	view.SetContent(content)
	return view
}

// Result returns the final action chosen by the user.
func (m *TUI) Result() Action {
	return m.result
}

// SetInsideSession sets the flag indicating we're inside a Zellij session.
func (m *TUI) SetInsideSession(v bool) {
	m.insideSession = v
}

// SetAutoAttach sets whether creating a session auto-attaches.
func (m *TUI) SetAutoAttach(v bool) {
	m.autoAttach = v
}
