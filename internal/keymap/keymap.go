package keymap

// KeyMatch defines a key binding as a set of string patterns (e.g. "up", "k").
// Using the MatchString approach from Bubble Tea v2's KeyPressEvent.
type KeyMatch []string

var (
	Up           = KeyMatch{"up", "k"}
	Down         = KeyMatch{"down", "j"}
	Enter        = KeyMatch{"enter"}
	NewSession   = KeyMatch{"n"}
	Delete       = KeyMatch{"d"}
	Kill         = KeyMatch{"x"}
	DeleteAll    = KeyMatch{"D"}
	Refresh      = KeyMatch{"r"}
	Quit         = KeyMatch{"q"}
	Escape       = KeyMatch{"esc"}
	ConfirmYes   = KeyMatch{"y"}
	ConfirmNo    = KeyMatch{"n"}
	CtrlC        = KeyMatch{"ctrl+c"}
)
