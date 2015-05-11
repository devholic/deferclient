package deferclient

import (
	"testing"
)

func TestNewAgent(t *testing.T) {

	a := NewAgent()

	if a.Name == "" {
		t.Error("not creating agent name")
	}

	if a.Totalmem == 0 {
		t.Error("not setting memory")
	}

}
