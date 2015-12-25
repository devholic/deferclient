// +build !go1.6

package deferclient

import (
	"fmt"
	"runtime"
)

// backtrace grabs the backtrace
func backTrace() (body string) {
	for skip := 1; ; skip++ {
		pc, file, line, ok := runtime.Caller(skip)
		if !ok {
			break
		}
		if file[len(file)-1] == 'c' {
			continue
		}
		f := runtime.FuncForPC(pc)
		body += fmt.Sprintf("%s:%d %s()\n", file, line, f.Name())
	}

	return body
}
