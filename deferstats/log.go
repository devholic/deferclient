// Package deferlog implements deferpanic error logging.
package deferstats

import (
	"fmt"
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
func (c *Client) Wrap(err error) {
	stack := BackTrace()

	go c.BaseClient.ShipTrace(stack, err.Error(), 0)
}

// WrapHTTPError wraps an error that occurs w/in a http request and
// sends the error to deferpanic w/the spanId
func WrapHTTPError(spanId int64, err error) {
}

// Persists ensures any panics will post to deferpanic website for
// tracking
// typically used in non http go-routines
func (c *Client) Persist() {
	if err := recover(); err != nil {
		c.BaseClient.Prep(err, 0)
	}
}

// PersistRepanic ensures any panics will post to deferpanic website for
// tracking, it also reissues the panic afterwards.
// typically used in non http go-routines
func (c *Client) PersistRepanic() {
	if err := recover(); err != nil {
		c.BaseClient.PrepSync(err, 0)
		panic(err)
	}
}
