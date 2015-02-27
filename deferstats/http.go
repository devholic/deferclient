package deferstats

import (
	"github.com/deferpanic/deferclient/deferclient"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

// curlist holds an array of DeferHTTPs (uri && latency)
var curlist []DeferHTTP

var latencyThreshold = 200

// appendHTTP adds a new http request to the list
func appendHTTP(startTime time.Time, path string, status_code int, span_id int64, parent_span_id int64) {
	endTime := time.Now()

	t := int(((endTime.Sub(startTime)).Nanoseconds() / 1000000))

	// only log over t
	if t > latencyThreshold {

		dh := DeferHTTP{
			Path:         path,
			Time:         t,
			StatusCode:   status_code,
			SpanId:       span_id,
			ParentSpanId: parent_span_id,
		}

		curlist = append(curlist, dh)

	}
}

func GetSpanId(r http.ResponseWriter) int64 {
	mPtr := (r).(*responseTracer)
	return mPtr.SpanId
}

// tracingResponseWriter implements a responsewriter with status
// exportable
type tracingResponseWriter interface {
	http.ResponseWriter
	Status() int
	Size() int
}

type responseTracer struct {
	w            http.ResponseWriter
	status       int
	size         int
	SpanId       int64
	ParentSpanId int64
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
func HTTPHandler(f func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()

		defer func() {
			if err := recover(); err != nil {
				// hack
				deferclient.Token = Token
				deferclient.Prep(err)
				// FIXME
				appendHTTP(startTime, r.URL.Path, 500, 0, 0)
			}
		}()

		var tracer *responseTracer

		tracer = &responseTracer{
			w: w,
		}

		tracer.SpanId = tracer.newId()

		r.ParseForm()
		deferParentSpanId := r.FormValue("defer_parent_span_id")
		if deferParentSpanId != "" {
			log.Println("deferParentSpanId: [" + deferParentSpanId + "]")
			tracer.ParentSpanId, _ = strconv.ParseInt(deferParentSpanId, 10, 64)
		}

		f(tracer, r)

		appendHTTP(startTime, r.URL.Path, tracer.Status(), tracer.SpanId, tracer.ParentSpanId)
	}
}

// AddRequest allows external libraries to add a http request
func AddRequest(start_time time.Time, path string, status_code int, span_id int64, parent_span_id int64) {

	if Verbose {
		log.Printf("Added manual request: %v\n", path)
	}

	// It's just an easier way to create third-party middlewares
	appendHTTP(start_time, path, status_code, span_id, parent_span_id)
}
