package deferclient

import (
	"testing"
)

func TestNewCommand(t *testing.T) {
	c := NewCommand(1, true, false, true)

	if c.Id != 1 {
		t.Error("not creating Id field")
	}
	if c.GenerateTrace != true {
		t.Error("not creating Generate trace field")
	}
	if c.Requested != false {
		t.Error("not creating Requested field")
	}
	if c.Executed != true {
		t.Error("not creating Executed field")
	}
}
