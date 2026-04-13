package model

import (
	"fmt"
	"testing"

	tea "charm.land/bubbletea/v2"

	"github.com/pageton/zellij-tui/internal/session"
	"github.com/pageton/zellij-tui/internal/ui"
)

// keyText creates a KeyPressMsg for a printable character.
func keyText(s string) tea.KeyPressMsg {
	return tea.KeyPressMsg{Text: s}
}

// keyCode creates a KeyPressMsg for a special key by code.
func keyCode(code rune) tea.KeyPressMsg {
	return tea.KeyPressMsg{Code: code}
}

// keyCtrl creates a KeyPressMsg for a ctrl+key combo.
func keyCtrl(code rune) tea.KeyPressMsg {
	return tea.KeyPressMsg{Code: code, Mod: tea.ModCtrl}
}

// isQuitCmd executes the cmd and returns true if it produces a tea.QuitMsg.
func isQuitCmd(cmd tea.Cmd) bool {
	if cmd == nil {
		return false
	}
	_, ok := cmd().(tea.QuitMsg)
	return ok
}

// freshModel returns a TUI with test sessions preloaded in StateList.
func freshModel() *TUI {
	m := New()
	m.sessions = []session.Session{
		{Name: "alpha", Created: "1h"},
		{Name: "beta", Created: "2h", IsCurrent: true},
		{Name: "gamma", Created: "3h"},
	}
	m.cursor = 0
	return m
}

// --- Cursor movement ---

func TestCursorDown(t *testing.T) {
	m := freshModel()
	_, cmd := m.Update(keyText("j"))
	if cmd != nil {
		t.Fatal("expected nil cmd for cursor move")
	}
	if m.cursor != 1 {
		t.Fatalf("cursor = %d, want 1", m.cursor)
	}
}

func TestCursorUp(t *testing.T) {
	m := freshModel()
	m.cursor = 1
	_, _ = m.Update(keyText("k"))
	if m.cursor != 0 {
		t.Fatalf("cursor = %d, want 0", m.cursor)
	}
}

func TestCursorClampsTop(t *testing.T) {
	m := freshModel()
	_, _ = m.Update(keyText("k"))
	if m.cursor != 0 {
		t.Fatalf("cursor = %d, want 0 (clamped at top)", m.cursor)
	}
}

func TestCursorClampsBottom(t *testing.T) {
	m := freshModel()
	m.cursor = 2
	_, _ = m.Update(keyText("j"))
	if m.cursor != 2 {
		t.Fatalf("cursor = %d, want 2 (clamped at bottom)", m.cursor)
	}
}

func TestCursorResetOnMovement(t *testing.T) {
	m := freshModel()
	m.quitPending = true
	m.err = "some error"
	_, _ = m.Update(keyText("j"))
	if m.quitPending {
		t.Error("quitPending should be cleared on cursor move")
	}
	if m.err != "" {
		t.Error("err should be cleared on cursor move")
	}
}

// --- Attach ---

func TestAttachQuit(t *testing.T) {
	m := freshModel()
	m.cursor = 0 // alpha
	_, cmd := m.Update(keyCode(tea.KeyEnter))
	if !isQuitCmd(cmd) {
		t.Fatal("expected quit cmd on attach")
	}
	if m.result.Op != ActionAttach || m.result.Name != "alpha" {
		t.Fatalf("result = %+v, want ActionAttach/alpha", m.result)
	}
}

func TestAttachCurrentSessionBlocked(t *testing.T) {
	m := freshModel()
	m.cursor = 1 // beta (IsCurrent)
	_, cmd := m.Update(keyCode(tea.KeyEnter))
	if isQuitCmd(cmd) {
		t.Fatal("should not quit when attaching to current session")
	}
	if m.err == "" {
		t.Error("expected error message for current session attach")
	}
}

func TestAttachEmptySessionsNoop(t *testing.T) {
	m := New()
	_, cmd := m.Update(keyCode(tea.KeyEnter))
	if cmd != nil {
		t.Fatal("expected nil cmd with no sessions")
	}
}

// --- State transitions: List -> Input ---

func TestNewSessionTransition(t *testing.T) {
	m := freshModel()
	_, cmd := m.Update(keyText("n"))
	if m.state != ui.StateInput {
		t.Fatalf("state = %d, want StateInput", m.state)
	}
	if cmd == nil {
		t.Error("expected non-nil cmd (input focus)")
	}
}

// --- State transitions: Input -> List (escape) ---

func TestInputEscapeBackToList(t *testing.T) {
	m := freshModel()
	_, _ = m.Update(keyText("n")) // enter input mode
	if m.state != ui.StateInput {
		t.Fatal("prerequisite: should be in StateInput")
	}
	_, _ = m.Update(keyCode(tea.KeyEscape))
	if m.state != ui.StateList {
		t.Fatalf("state = %d, want StateList after escape", m.state)
	}
}

// --- Input validation ---

func TestInputEmptyNameNoop(t *testing.T) {
	m := freshModel()
	_, _ = m.Update(keyText("n")) // enter input mode
	_, cmd := m.Update(keyCode(tea.KeyEnter))
	if isQuitCmd(cmd) {
		t.Fatal("empty name should not quit")
	}
}

func TestInputInvalidNameRejected(t *testing.T) {
	m := freshModel()
	_, _ = m.Update(keyText("n"))
	m.input.SetValue("--bad")
	_, cmd := m.Update(keyCode(tea.KeyEnter))
	if isQuitCmd(cmd) {
		t.Fatal("invalid name should not quit")
	}
	if m.err == "" {
		t.Error("expected error for invalid name")
	}
}

func TestInputValidNameAutoAttach(t *testing.T) {
	m := freshModel()
	m.autoAttach = true
	_, _ = m.Update(keyText("n"))
	m.input.SetValue("my-session")
	_, cmd := m.Update(keyCode(tea.KeyEnter))
	if !isQuitCmd(cmd) {
		t.Fatal("valid name + autoAttach should quit")
	}
	if m.result.Op != ActionCreate || m.result.Name != "my-session" {
		t.Fatalf("result = %+v, want ActionCreate/my-session", m.result)
	}
}

func TestInputValidNameNoAutoAttach(t *testing.T) {
	m := freshModel()
	m.autoAttach = false
	_, _ = m.Update(keyText("n"))
	m.input.SetValue("bg-session")
	_, cmd := m.Update(keyCode(tea.KeyEnter))
	if isQuitCmd(cmd) {
		t.Fatal("no autoAttach should not quit")
	}
	if m.state != ui.StateList {
		t.Fatalf("state = %d, want StateList after background create", m.state)
	}
}

// --- State transitions: List -> Confirm ---

func TestDeleteTransitionToConfirm(t *testing.T) {
	m := freshModel()
	_, _ = m.Update(keyText("d"))
	if m.state != ui.StateConfirmDelete {
		t.Fatalf("state = %d, want StateConfirmDelete", m.state)
	}
}

func TestKillTransitionToConfirm(t *testing.T) {
	m := freshModel()
	_, _ = m.Update(keyText("x"))
	if m.state != ui.StateConfirmKill {
		t.Fatalf("state = %d, want StateConfirmKill", m.state)
	}
}

func TestDeleteAllTransitionToConfirm(t *testing.T) {
	m := freshModel()
	_, _ = m.Update(keyText("D"))
	if m.state != ui.StateConfirmDeleteAll {
		t.Fatalf("state = %d, want StateConfirmDeleteAll", m.state)
	}
}

func TestDeleteEmptySessionsNoop(t *testing.T) {
	m := New()
	_, cmd := m.Update(keyText("d"))
	if cmd != nil {
		t.Fatal("expected nil cmd with no sessions")
	}
	if m.state != ui.StateList {
		t.Fatal("state should stay StateList with no sessions")
	}
}

// --- Confirm -> List (yes/no/escape) ---

func TestConfirmYesReturnsAction(t *testing.T) {
	m := freshModel()
	_, _ = m.Update(keyText("d")) // enter confirm delete
	if m.state != ui.StateConfirmDelete {
		t.Fatal("prerequisite: should be in confirm state")
	}
	_, cmd := m.Update(keyText("y"))
	if m.state != ui.StateList {
		t.Fatalf("state = %d, want StateList after confirm yes", m.state)
	}
	if cmd == nil {
		t.Error("confirm yes should return a command (the delete action)")
	}
}

func TestConfirmNoBackToList(t *testing.T) {
	m := freshModel()
	_, _ = m.Update(keyText("d"))
	_, cmd := m.Update(keyText("n"))
	if m.state != ui.StateList {
		t.Fatalf("state = %d, want StateList after confirm no", m.state)
	}
	if cmd != nil {
		t.Error("confirm no should return nil cmd")
	}
	if m.confirmAction != nil {
		t.Error("confirmAction should be nil after cancel")
	}
}

func TestConfirmEscapeBackToList(t *testing.T) {
	m := freshModel()
	_, _ = m.Update(keyText("x")) // confirm kill
	_, _ = m.Update(keyCode(tea.KeyEscape))
	if m.state != ui.StateList {
		t.Fatalf("state = %d, want StateList after confirm escape", m.state)
	}
}

// --- Quit behavior ---

func TestQuitFirstPressSetsPending(t *testing.T) {
	m := freshModel()
	_, cmd := m.Update(keyText("q"))
	if !m.quitPending {
		t.Error("first q should set quitPending")
	}
	if cmd == nil {
		t.Error("first q should return a tick cmd")
	}
}

func TestQuitSecondPressQuits(t *testing.T) {
	m := freshModel()
	_, _ = m.Update(keyText("q"))    // first press
	_, cmd := m.Update(keyText("q")) // second press
	if !isQuitCmd(cmd) {
		t.Fatal("second q should produce quit")
	}
	if m.result.Op != ActionNone {
		t.Error("result should be ActionNone on quit")
	}
}

func TestCtrlCQuitsImmediately(t *testing.T) {
	m := freshModel()
	_, cmd := m.Update(keyCtrl('c'))
	if !isQuitCmd(cmd) {
		t.Fatal("ctrl+c should produce quit")
	}
	if m.result.Op != ActionNone {
		t.Error("result should be ActionNone on ctrl+c")
	}
}

// --- sessionsLoadedMsg handler ---

func TestSessionsLoadedSuccess(t *testing.T) {
	m := New()
	sessions := []session.Session{
		{Name: "work", Created: "5m"},
	}
	_, _ = m.Update(sessionsLoadedMsg{sessions: sessions})
	if len(m.sessions) != 1 {
		t.Fatalf("sessions len = %d, want 1", len(m.sessions))
	}
	if m.err != "" {
		t.Error("err should be cleared on success")
	}
}

func TestSessionsLoadedError(t *testing.T) {
	m := New()
	_, _ = m.Update(sessionsLoadedMsg{err: fmt.Errorf("boom")})
	if m.err == "" {
		t.Error("expected error message from failed load")
	}
}

func TestSessionsLoadedCursorClamp(t *testing.T) {
	m := New()
	m.cursor = 5
	_, _ = m.Update(sessionsLoadedMsg{sessions: []session.Session{
		{Name: "only", Created: "1s"},
	}})
	if m.cursor != 0 {
		t.Fatalf("cursor = %d, want 0 (clamped after load)", m.cursor)
	}
}

// --- WindowSizeMsg ---

func TestWindowSize(t *testing.T) {
	m := New()
	_, _ = m.Update(tea.WindowSizeMsg{Width: 100, Height: 40})
	if m.width != 100 || m.height != 40 {
		t.Fatalf("width=%d height=%d, want 100x40", m.width, m.height)
	}
}

// --- Setters ---

func TestSetInsideSession(t *testing.T) {
	m := New()
	m.SetInsideSession(true)
	if !m.insideSession {
		t.Error("insideSession should be true")
	}
}

func TestSetAutoAttach(t *testing.T) {
	m := New()
	m.SetAutoAttach(false)
	if m.autoAttach {
		t.Error("autoAttach should be false")
	}
}

// --- Result ---

func TestResultDefault(t *testing.T) {
	m := New()
	r := m.Result()
	if r.Op != ActionNone {
		t.Errorf("default result op = %d, want ActionNone", r.Op)
	}
}
