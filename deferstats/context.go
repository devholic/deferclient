// Based on Joe Shaw article https://joeshaw.org/net-context-and-http-handler/

package deferstats

import (
	"fmt"
	"golang.org/x/net/context"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// ContextHandler is based standard http Handler inteface
type ContextHandler interface {
	ServeHTTPContext(context.Context, http.ResponseWriter, *http.Request)
}

// ContextHandlerFunc is based standard http Handler func
type ContextHandlerFunc func(context.Context, http.ResponseWriter, *http.Request)

// ServeHTTPContext is implementation of context based standard http Handler
func (h ContextHandlerFunc) ServeHTTPContext(ctx context.Context, w http.ResponseWriter, req *http.Request) {
	h(ctx, w, req)
}

// ContextTracer passes spanId/parentSpanId parameters using context
type ContextTracer struct {
	SpanId       int64
	ParentSpanId int64
}

// GetContextSpanIdString is a convenience method to get the string equivalent
// of a span id from context
func GetContextSpanIdString(ctx context.Context) string {
	return strconv.FormatInt(GetContextSpanId(ctx), 10)
}

// GetContextSpanId returns the span id from this context
func GetContextSpanId(ctx context.Context) int64 {
	mPtr := (ctx.Value("ContextTracer")).(*ContextTracer)
	return mPtr.SpanId
}

// GetContextParentSpanIdString is a convenience method to get the string equivalent
// of a parent span id from context
func GetContextParentSpanIdString(ctx context.Context) string {
	return strconv.FormatInt(GetContextParentSpanId(ctx), 10)
}

// GetContextParentSpanId returns the parent span id from this context
func GetContextParentSpanId(ctx context.Context) int64 {
	mPtr := (ctx.Value("ContextTracer")).(*ContextTracer)
	return mPtr.ParentSpanId
}

func (t *ContextTracer) newId() int64 {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Int63()
}

// ResponseWriterExt implements http.ResponseWriter with extended methods
type ResponseWriterExt struct {
	w      http.ResponseWriter
	status int
	size   int
}

// Header is implementaion of standard http ResponseWriter Header method
func (e *ResponseWriterExt) Header() http.Header {
	return e.w.Header()
}

// Write is implementaion of standard http ResponseWriter Header method and setting the size and status
func (e *ResponseWriterExt) Write(b []byte) (int, error) {
	if e.status == 0 {
		// The status will be StatusOK if WriteHeader has not been
		// called yet
		e.status = http.StatusOK
	}
	size, err := e.w.Write(b)
	e.size += size
	return size, err
}

// WriteHeader is implementaion of standard http ResponseWriter WriteHeader method and setting the status
func (e *ResponseWriterExt) WriteHeader(s int) {
	e.w.WriteHeader(s)
	e.status = s
}

// Status returns the HTTP status code
func (e *ResponseWriterExt) Status() int {
	return e.status
}

// Size returns the HTTP size
func (e *ResponseWriterExt) Size() int {
	return e.size
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
		startTime, ext, tracer, headers := c.ContextBeforeRequest(w, r)
		ctx = context.WithValue(ctx, "ContextTracer", tracer)

		defer func() {
			if err := recover(); err != nil {
				c.BaseClient.Prep(err, tracer.SpanId)
				c.ContextAfterRequest(startTime, tracer, r, headers, 500, true)

				errorMsg := fmt.Sprintf("%v", err)
				WritePanicResponse(w, r, errorMsg)
			}
		}()

		f.ServeHTTPContext(ctx, ext, r)

		c.ContextAfterRequest(startTime, tracer, r, headers, ext.Status(), false)
	})
}

// ContextBeforeRequest is called before request processing in context handler
func (c *Client) ContextBeforeRequest(w http.ResponseWriter, r *http.Request) (
	startTime time.Time, ext *ResponseWriterExt, tracer *ContextTracer, headers map[string]string) {
	startTime = time.Now()

	ext = &ResponseWriterExt{
		w: w,
	}

	tracer = new(ContextTracer)
	tracer.SpanId = tracer.newId()

	// add headers
	headers = make(map[string]string, len(r.Header))
	for k, v := range r.Header {
		headers[k] = strings.Join(v, ",")

		// grab SOA tracing header if present
		if k == "X-Dpparentspanid" {
			tracer.ParentSpanId, _ = strconv.ParseInt(v[0], 10, 64)
		}
	}

	return startTime, ext, tracer, headers
}

// ContextAfterRequest is called after request processing in context handler
func (c *Client) ContextAfterRequest(startTime time.Time, tracer *ContextTracer, r *http.Request,
	headers map[string]string, status_code int, isproblem bool) {
	appendHTTP(startTime, r.URL.Path, r.Method, status_code, tracer.SpanId,
		tracer.ParentSpanId, isproblem, headers)
}
