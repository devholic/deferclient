package errors

import (
	"github.com/deferpanic/deferclient/deferclient"
	"testing"
)

func TestErrorsNew(t *testing.T) {
	deferclient.NoPost = true

	err := New("Testy McTest")
	if err == nil {
		t.Error("error is not being returned")
	}

	if err.Error() != "Testy McTest" {
		t.Error("error Msg is not being set")
	}

	if err.GetBackTrace() == "" {
		t.Error("BackTrace is empty")
	}
}
