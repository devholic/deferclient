// Package errors implements deferpanic error logging.
package errors

import (
	"fmt"
	"github.com/deferpanic/deferclient/deferclient"
	"runtime"
)

// DeferPanicError is an interface impmenting the error interface and
// adding the backtrace
type DeferPanicError interface {
	Error() string

	GetBackTrace() string
}

// DeferPanicBaseError contains the error msg, the backtrace, and the
// original error value
type DeferPanicBaseError struct {
	Msg       string
	BackTrace string
	orig      error
}

// GetBackTrace is a getter for obtaining the backtrace for a deferpanic
// error
func (e *DeferPanicBaseError) GetBackTrace() string {
	return e.BackTrace
}

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
func Wrap(err error, msg string) DeferPanicError {
	stack := BackTrace()
	deferclient.ShipTrace(stack, msg)

	return &DeferPanicBaseError{
		Msg:       msg,
		BackTrace: stack,
		orig:      err,
	}
}

// new instantiates a new error and ships the backtrace to deferpanic
func New(msg string) DeferPanicError {
	stack := BackTrace()
	go deferclient.ShipTrace(stack, msg)

	return &DeferPanicBaseError{
		Msg:       msg,
		BackTrace: stack,
	}
}

// Error implments the error interface
func (e *DeferPanicBaseError) Error() string {
	return e.Msg
}
