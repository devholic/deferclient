// Package errors implements deferpanic error logging.
// graciously stolen from dropbox
package errors

import (
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

func BackTrace() string {
	buf := make([]byte, 1<<16)
	runtime.Stack(buf, false)

	sz := len(buf) - 1
	body := string(buf[:sz])

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
