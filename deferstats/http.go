package deferstats

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"
)

// FIXME
// soon to be DEPRECATED
var (
	// curlist holds an array of DeferHTTPs (uri && latency)
	curlist = &deferHTTPList{}

	// latencyThreshold is the threshold in milliseconds that if
	// exceeded a request will be added to the curlist
	latencyThreshold = 200
)

// DeferHTTP holds the path uri and latency for each request
type DeferHTTP struct {
	Path         string            `json:"Path"`
	Method       string            `json:"Method"`
	StatusCode   int               `json:"StatusCode"`
	Time         int               `json:"Time"`
	SpanId       int64             `json:"SpanId"`
	ParentSpanId int64             `json:"ParentSpanId"`
	IsProblem    bool              `json:"IsProblem"`
	Headers      map[string]string `json:"Headers"`
}

// deferHTTPList is used to keep a list of DeferHTTP objects
// and interact with them in a thread-safe manner
type deferHTTPList struct {
	lock sync.RWMutex
	list []DeferHTTP
}

// tracingResponseWriter implements a responsewriter with status
// exportable
type tracingResponseWriter interface {
	http.ResponseWriter
	Status() int
	Size() int
}

// responseTracer implements a responsewriter with spanId/parentSpanId
type responseTracer struct {
	w            http.ResponseWriter
	status       int
	size         int
	SpanId       int64
	ParentSpanId int64
}

// Add adds a DeferHTTP object to the list
func (d *deferHTTPList) Add(item DeferHTTP) {
	d.lock.Lock()
	d.list = append(d.list, item)
	d.lock.Unlock()
}

// List returns a copy of the list
func (d *deferHTTPList) List() []DeferHTTP {
	d.lock.RLock()
	list := make([]DeferHTTP, len(d.list))
	for i, v := range d.list {
		list[i] = v
	}
	d.lock.RUnlock()
	return list
}

// Reset removes all entries from the list
func (d *deferHTTPList) Reset() {
	d.lock.Lock()
	d.list = []DeferHTTP{}
	d.lock.Unlock()
}

// WritePanicResponse is an overridable function that, by default, writes the contents of the panic
// error message with a 500 Internal Server Error.
var WritePanicResponse = func(w http.ResponseWriter, r *http.Request, errMsg string) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(errMsg))
}

// appendHTTP adds a new http request to the list
func appendHTTP(startTime time.Time, path string, method string, status_code int, span_id int64,
	parent_span_id int64, isProblem bool, headers map[string]string) {
	endTime := time.Now()

	t := int(((endTime.Sub(startTime)).Nanoseconds() / 1000000))

	// only log if t over latencyThreshold or if a panic/error occurred
	if (t > latencyThreshold) || isProblem {

		dh := DeferHTTP{
			Path:         path,
			Method:       method,
			Time:         t,
			StatusCode:   status_code,
			SpanId:       span_id,
			ParentSpanId: parent_span_id,
			IsProblem:    isProblem,
			Headers:      headers,
		}

		curlist.Add(dh)

	}
}

// GetSpanIdString is a conveinence method to get the string equivalent
// of a span id
func GetSpanIdString(r http.ResponseWriter) string {
	return strconv.FormatInt(GetSpanId(r), 10)
}

// GetSpanId returns the span id for this http request
func GetSpanId(r http.ResponseWriter) int64 {
	mPtr := (r).(*responseTracer)
	return mPtr.SpanId
}

func (l *responseTracer) newId() int64 {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Int63()
}

func (l *responseTracer) Header() http.Header {
	return l.w.Header()
}

func (l *responseTracer) Write(b []byte) (int, error) {
	if l.status == 0 {
		// The status will be StatusOK if WriteHeader has not been
		// called yet
		l.status = http.StatusOK
	}
	size, err := l.w.Write(b)
	l.size += size
	return size, err
}

// WriteHeader sets the header
func (l *responseTracer) WriteHeader(s int) {
	l.w.WriteHeader(s)
	l.status = s
}

// Status returns the HTTP status code
func (l *responseTracer) Status() int {
	return l.status
}

func (l *responseTracer) Size() int {
	return l.size
}

// HTTPHandler wraps a http handler and captures the latency of each
// request
// this currently happens in a global list :( - TBFS
func (c *Client) HTTPHandlerFunc(f http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		// cp body - header read inadvert reads this
		var bodyBytes []byte
		if r.Body != nil {
			bodyBytes, _ = ioutil.ReadAll(r.Body)
		}

		startTime := time.Now()

		var tracer *responseTracer

		tracer = &responseTracer{
			w: w,
		}

		tracer.SpanId = tracer.newId()

		v, e := url.ParseQuery(string(bodyBytes))
		if e != nil {
			log.Println(e)
		}

		deferParentSpanId := ""
		if v["defer_parent_span_id"] != nil {
			deferParentSpanId = v["defer_parent_span_id"][0]
		}

		if deferParentSpanId != "" {
			if c.Verbose {
				log.Println("deferParentSpanId: [" + deferParentSpanId + "]")
			}
			tracer.ParentSpanId, _ = strconv.ParseInt(deferParentSpanId, 10, 64)
		}

		// add headers
		headers := make(map[string]string, len(r.Header))

		for k, v := range r.Header {
			headers[k] = strings.Join(v, ",")
		}

		defer func() {
			if err := recover(); err != nil {
				c.BaseClient.Prep(err, tracer.SpanId)

				appendHTTP(startTime, r.URL.Path, r.Method, 500, tracer.SpanId, tracer.ParentSpanId,
					true, headers)

				errorMsg := fmt.Sprintf("%v", err)
				WritePanicResponse(w, r, errorMsg)
			}
		}()

		// Restore our body to use in the request
		r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

		f(tracer, r)

		appendHTTP(startTime, r.URL.Path, r.Method, tracer.Status(), tracer.SpanId, tracer.ParentSpanId,
			false, headers)
	}
}
