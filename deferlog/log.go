// Package deferlog implements deferpanic error logging.
package deferlog

import (
	"fmt"
	"github.com/deferpanic/deferclient/deferclient"
	"runtime"
)

// Token is your deferpanic token available in settings
var Token string

// Environment sets an environment tag to differentiate between separate
// environments - default is production.
var Environment = "production"

// Backtrace grabs the backtrace
func BackTrace() (body string) {

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

// Wrap wraps an error and ships the backtrace to deferpanic
func Wrap(err error) {
	stack := BackTrace()
	deferclient.Token = Token
	deferclient.Environment = Environment

	go deferclient.ShipTrace(stack, err.Error())
}
