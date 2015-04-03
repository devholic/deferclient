// Package deferlog implements deferpanic error logging.
package deferstats

import (
	"fmt"
	"github.com/deferpanic/deferclient/deferclient"
	"runtime"
)

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
	deferclient.AppGroup = AppGroup

	go deferclient.ShipTrace(stack, err.Error())
}
