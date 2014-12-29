// Package errors implements deferpanic error logging.
// graciously stolen from dropbox
package errors

import (
	"fmt"
	"github.com/deferpanic/deferclient/deferclient"
	"runtime"
)

type DeferPanicError interface {
	Error() string

	GetBackTrace() string
}

type DeferPanicBaseError struct {
	Msg       string
	BackTrace string
	orig      error
}

func (e *DeferPanicBaseError) GetBackTrace() string {
	return e.BackTrace
}

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

func Wrap(err error, msg string) DeferPanicError {
	stack := BackTrace()
	deferclient.ShipTrace(stack, msg)

	return &DeferPanicBaseError{
		Msg:       msg,
		BackTrace: stack,
		orig:      err,
	}
}

func New(msg string) DeferPanicError {
	stack := BackTrace()
	go deferclient.ShipTrace(stack, msg)

	return &DeferPanicBaseError{
		Msg:       msg,
		BackTrace: stack,
	}
}

func (e *DeferPanicBaseError) Error() string {
	return e.Msg
}
