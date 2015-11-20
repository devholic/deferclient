package deferclient

import (
	"testing"
)

func TestNewMemProfile(t *testing.T) {
	c := NewMemProfile([]byte("Out"), []byte("Pkg"), 1, true)

	if string(c.Out) != "Out" {
		t.Error("not creating Out field")
	}
	if string(c.Pkg) != "Pkg" {
		t.Error("not creating Pkg field")
	}
	if c.CommandId != 1 {
		t.Error("not creating CommandId field")
	}
	if c.Ignored != true {
		t.Error("not creating Ignored field")
	}
}
