package model

import "testing"

func TestActionString(t *testing.T) {
	tests := []struct {
		action Action
		want   string
	}{
		{Action{Op: ActionAttach, Name: "work"}, "ATTACH:work"},
		{Action{Op: ActionCreate, Name: "dev"}, "CREATE:dev"},
		{Action{Op: ActionNone}, ""},
	}

	for _, tt := range tests {
		got := tt.action.String()
		if got != tt.want {
			t.Errorf("Action{%d, %q}.String() = %q, want %q", tt.action.Op, tt.action.Name, got, tt.want)
		}
	}
}
