// Based on Joe Shaw article https://joeshaw.org/net-context-and-http-handler/

package deferstats

import (
	"fmt"
	"golang.org/x/net/context"
	"net/http"
)

type ContextHandler interface {
	ServeHTTPContext(context.Context, http.ResponseWriter, *http.Request)
}

type ContextHandlerFunc func(context.Context, http.ResponseWriter, *http.Request)

func (h ContextHandlerFunc) ServeHTTPContext(ctx context.Context, rw http.ResponseWriter, req *http.Request) {
	h(ctx, rw, req)
}

func middleware(h ContextHandler) ContextHandler {
	return ContextHandlerFunc(func(ctx context.Context, rw http.ResponseWriter, req *http.Request) {
		//		ctx = newContextWithRequestID(ctx, req)
		h.ServeHTTPContext(ctx, rw, req)
	})
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
		startTime, tracer, headers := c.BeforeRequest(w, r)

		defer func() {
			if err := recover(); err != nil {
				c.BaseClient.Prep(err, tracer.SpanId)
				c.AfterRequest(startTime, tracer, r, headers, 500, true)

				errorMsg := fmt.Sprintf("%v", err)
				WritePanicResponse(w, r, errorMsg)
			}
		}()

		f.ServeHTTPContext(ctx, tracer, r)

		c.AfterRequest(startTime, tracer, r, headers, tracer.Status(), false)
	})
}
