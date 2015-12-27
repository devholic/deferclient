// +build go1.6

package deferclient

import (
	"runtime"
)

// backtrace grabs the backtrace
func backTrace() (body string) {
	trace := make([]byte, 65536)
	_ = runtime.Stack(trace, false)
	body = string(trace)

	return body
}
