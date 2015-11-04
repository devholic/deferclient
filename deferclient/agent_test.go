package deferclient

import (
	"testing"
)

func TestNewAgent(t *testing.T) {

	a := NewAgent()

	if a.Name == "" {
		t.Error("not creating agent name")
	}

	if a.Cpucores != 0 {
		t.Error("not creating cpu cores")
	}

	if a.Goarch != "" {
		t.Error("not creating go architecture")
	}

	if a.Goos != "" {
		t.Error("not creating go os")
	}

	if a.Totalmem == 0 {
		t.Error("not creating memory")
	}

	if a.Govers != "" {
		t.Error("not creating go version")
	}

	if a.ApiVersion != "" {
		t.Error("not creating api version")
	}

	if a.CRC32 != 0 {
		t.Error("not creating crc32")
	}

	if a.Size != 0 {
		t.Error("not creating size")
	}
}
