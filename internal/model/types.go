package model

// Action represents the output protocol between the TUI binary and the shell wrapper.
type Action struct {
	Op   ActionOp
	Name string
}

// ActionOp represents the type of action the user chose.
type ActionOp int

// ActionOp values represent the type of action the user chose.
const (
	ActionNone   ActionOp = iota // no action
	ActionAttach                 // attach to existing session
	ActionCreate                 // create and attach
)

func (a Action) String() string {
	switch a.Op {
	case ActionAttach:
		return "ATTACH:" + a.Name
	case ActionCreate:
		return "CREATE:" + a.Name
	default:
		return ""
	}
}
