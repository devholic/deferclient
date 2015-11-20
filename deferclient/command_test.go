package deferclient

import (
	"testing"
)

func TestNewCommand(t *testing.T) {
	c := NewCommand(1, CommandTypeTrace, false, true)

	if c.Id != 1 {
		t.Error("not creating Id field")
	}
	if c.Type != CommandTypeTrace {
		t.Error("not creating Type field")
	}
	if c.Requested != false {
		t.Error("not creating Requested field")
	}
	if c.Executed != true {
		t.Error("not creating Executed field")
	}
}
