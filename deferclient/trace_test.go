package deferclient

import (
	"testing"
)

func TestNewTrace(t *testing.T) {
	tr := NewTrace([]byte("Out"), []byte("Pkg"), 1, 2, 3)

	if string(tr.Out) != "Out" {
		t.Error("not creating Out field")
	}
	if string(tr.Pkg) != "Pkg" {
		t.Error("not creating Pkg field")
	}
	if tr.CRC32 != 1 {
		t.Error("not creating CRC32 field")
	}
	if tr.Size != 2 {
		t.Error("not creating Size field")
	}
	if tr.CommandId != 3 {
		t.Error("not creating CommandId field")
	}
}
