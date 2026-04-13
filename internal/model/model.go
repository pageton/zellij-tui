// Package model implements the Bubble Tea TUI model for session management.
package model

import (
	"fmt"
	"time"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	lipgloss "charm.land/lipgloss/v2"

	"github.com/pageton/zellij-tui/internal/keymap"
	"github.com/pageton/zellij-tui/internal/session"
	"github.com/pageton/zellij-tui/internal/ui"
	"github.com/pageton/zellij-tui/internal/zellij"
)

const quitTimeout = 2 * time.Second

// refreshSessions returns a command that fetches the current session list.
func refreshSessions() tea.Cmd {
	return func() tea.Msg {
		sessions, err := zellij.ListSessions()
		return sessionsLoadedMsg{sessions: sessions, err: err}
	}
}

// resetStatus clears transient UI state (quit pending flag and error message).
func (m *TUI) resetStatus() {
	m.quitPending = false
	m.err = ""
}

// cursorValid returns true if the cursor points to an existing session.
func (m *TUI) cursorValid() bool {
	return len(m.sessions) > 0 && m.cursor >= 0 && m.cursor < len(m.sessions)
}

// tickMsg is sent when the quit timer expires.
type tickMsg time.Time

// sessionsLoadedMsg is sent when session listing completes.
type sessionsLoadedMsg struct {
	sessions []session.Session
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
	sessions      []session.Session
	cursor        int
	input         textinput.Model
	err           string
	quitPending   bool
	insideSession bool
	autoAttach    bool
	width         int
	height        int
	result        Action
	confirmAction func() tea.Cmd // stored action to run on confirm
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

// Init returns the initial Bubble Tea commands.
func (m *TUI) Init() tea.Cmd {
	return tea.Batch(
		refreshSessions(),
		tea.RequestWindowSize,
	)
}

// Update handles Bubble Tea messages and returns the updated model.
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
			return m, nil
		}
		m.err = ""
		return m, refreshSessions()

	case createBackgroundResultMsg:
		if msg.err != nil {
			m.err = fmt.Sprintf("Failed to create session: %v", msg.err)
			return m, nil
		}
		m.err = ""
		return m, tea.Tick(500*time.Millisecond, func(_ time.Time) tea.Msg {
			return refreshSessions()()
		})

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
	case ui.StateConfirmDelete, ui.StateConfirmKill, ui.StateConfirmDeleteAll:
		return m.updateConfirm(msg)
	}

	return m, nil
}

func (m *TUI) updateList(msg tea.Msg) (tea.Model, tea.Cmd) {
	key, ok := msg.(tea.KeyPressMsg)
	if !ok {
		return m, nil
	}

	switch {
	case keymap.Up.Match(key):
		if m.cursor > 0 {
			m.cursor--
		}
		m.resetStatus()
		return m, nil

	case keymap.Down.Match(key):
		if m.cursor < len(m.sessions)-1 {
			m.cursor++
		}
		m.resetStatus()
		return m, nil

	case keymap.Enter.Match(key):
		if !m.cursorValid() {
			return m, nil
		}
		if m.sessions[m.cursor].IsCurrent {
			m.err = "Already in this session"
			return m, nil
		}
		m.result = Action{Op: ActionAttach, Name: m.sessions[m.cursor].Name}
		return m, tea.Quit

	case keymap.NewSession.Match(key):
		m.state = ui.StateInput
		m.input.SetValue("")
		cmd := m.input.Focus()
		m.resetStatus()
		return m, cmd

	case keymap.Delete.Match(key):
		if !m.cursorValid() {
			return m, nil
		}
		name := m.sessions[m.cursor].Name
		m.confirmAction = func() tea.Cmd {
			return func() tea.Msg {
				err := zellij.DeleteSession(name)
				return deleteResultMsg{err: err}
			}
		}
		m.state = ui.StateConfirmDelete
		m.resetStatus()
		return m, nil

	case keymap.Kill.Match(key):
		if !m.cursorValid() {
			return m, nil
		}
		name := m.sessions[m.cursor].Name
		m.confirmAction = func() tea.Cmd {
			return func() tea.Msg {
				err := zellij.KillSession(name)
				return deleteResultMsg{err: err}
			}
		}
		m.state = ui.StateConfirmKill
		m.resetStatus()
		return m, nil

	case keymap.DeleteAll.Match(key):
		if len(m.sessions) == 0 {
			return m, nil
		}
		m.confirmAction = func() tea.Cmd {
			return func() tea.Msg {
				err := zellij.DeleteAllSessions()
				return deleteResultMsg{err: err}
			}
		}
		m.state = ui.StateConfirmDeleteAll
		m.resetStatus()
		return m, nil

	case keymap.Refresh.Match(key):
		m.resetStatus()
		return m, refreshSessions()

	case keymap.Quit.Match(key):
		if m.quitPending {
			m.result = Action{Op: ActionNone}
			return m, tea.Quit
		}
		m.quitPending = true
		return m, tea.Tick(quitTimeout, func(t time.Time) tea.Msg {
			return tickMsg(t)
		})

	case keymap.CtrlC.Match(key):
		m.result = Action{Op: ActionNone}
		return m, tea.Quit

	default:
		m.quitPending = false
	}

	return m, nil
}

func (m *TUI) updateInput(msg tea.Msg) (tea.Model, tea.Cmd) {
	if key, ok := msg.(tea.KeyPressMsg); ok {
		switch {
		case keymap.Escape.Match(key):
			m.state = ui.StateList
			m.input.Blur()
			return m, nil

		case keymap.Enter.Match(key):
			name := m.input.Value()
			if name == "" {
				return m, nil
			}
			if !session.ValidSessionName(name) {
				m.err = "Invalid name: use letters, numbers, _ . - (no leading -)"
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

func (m *TUI) updateConfirm(msg tea.Msg) (tea.Model, tea.Cmd) {
	key, ok := msg.(tea.KeyPressMsg)
	if !ok {
		return m, nil
	}

	switch {
	case keymap.ConfirmYes.Match(key):
		action := m.confirmAction
		m.state = ui.StateList
		m.confirmAction = nil
		if action == nil {
			return m, nil
		}
		return m, action()

	case keymap.ConfirmNo.Match(key), keymap.Escape.Match(key):
		m.state = ui.StateList
		m.confirmAction = nil
		return m, nil
	}

	return m, nil
}

// View renders the current TUI state.
func (m *TUI) View() tea.View {
	content := ui.Render(ui.ViewData{
		Sessions:      m.sessions,
		Cursor:        m.cursor,
		State:         m.state,
		InputView:     m.input.View(),
		ErrMsg:        m.err,
		QuitPending:   m.quitPending,
		InsideSession: m.insideSession,
		TermWidth:     m.width,
	})

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
