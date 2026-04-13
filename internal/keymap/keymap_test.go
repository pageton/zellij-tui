package keymap

import "testing"

func TestKeyBindingsExist(t *testing.T) {
	bindings := []struct {
		name string
		km   KeyMatch
	}{
		{"Up", Up},
		{"Down", Down},
		{"Enter", Enter},
		{"NewSession", NewSession},
		{"Delete", Delete},
		{"Kill", Kill},
		{"DeleteAll", DeleteAll},
		{"Refresh", Refresh},
		{"Quit", Quit},
		{"Escape", Escape},
		{"CtrlC", CtrlC},
	}

	for _, b := range bindings {
		if len(b.km) == 0 {
			t.Errorf("%s has no key patterns", b.name)
		}
	}
}
