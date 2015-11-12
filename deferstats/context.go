// Based on Joe Shaw article https://joeshaw.org/net-context-and-http-handler/

package deferstats

import (
	"fmt"
	"golang.org/x/net/context"
	"net/http"
)

var spanId int64

type ContextHandler interface {
	ServeHTTPContext(context.Context, http.ResponseWriter, *http.Request)
}

type ContextHandlerFunc func(context.Context, http.ResponseWriter, *http.Request)

func (h ContextHandlerFunc) ServeHTTPContext(ctx context.Context, w http.ResponseWriter, req *http.Request) {
	h(ctx, w, req)
}

func newContextWithSpanID(ctx context.Context, req *http.Request) context.Context {
	spanId++
	return context.WithValue(ctx, "spanId", spanId)
}

// HTTPContextHandlerFunc wraps a http handler func and captures the latency of each
// request
// this currently happens in a global list :( - TBFS
func (c *Client) HTTPContextHandlerFunc(f ContextHandlerFunc) ContextHandlerFunc {
	return c.HTTPContextHandler(f).(ContextHandlerFunc)
}

// HTTPContextHandler wraps a http handler and captures the latency of each
// request
// this currently happens in a global list :( - TBFS
func (c *Client) HTTPContextHandler(f ContextHandler) ContextHandler {
	return ContextHandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		//		startTime, tracer, headers := c.BeforeRequest(w, r)
		ctx = newContextWithSpanID(ctx, r)
		defer func() {
			if err := recover(); err != nil {
				c.BaseClient.Prep(err, spanId /*tracer.SpanId*/)
				//				c.AfterRequest(startTime, tracer, r, headers, 500, true)

				errorMsg := fmt.Sprintf("%v", err)
				WritePanicResponse(w, r, errorMsg)
			}
		}()

		f.ServeHTTPContext(ctx, w, r)

		//		c.AfterRequest(startTime, tracer, r, headers, tracer.Status(), false)
	})
}
