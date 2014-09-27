package deferclient

import (
	"testing"
)

func TestCleanTrace(t *testing.T) {

	var body = `
some text
with a linebreak and a	tab
`

	nbody := cleanTrace(body)

	if nbody != "\\nsome text\\nwith a linebreak and a\\ttab\\n" {
		t.Error("not escaping line breaks and tabs")
	}
}
