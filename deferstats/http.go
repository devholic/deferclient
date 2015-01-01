package deferstats

import (
	"github.com/deferpanic/deferclient/deferclient"
	"net/http"
	"time"
)

// curlist holds an array of DeferHTTPs (uri && latency)
var curlist []DeferHTTP

func appendHTTP(startTime time.Time, path string) {
	endTime := time.Now()

	t := int(((endTime.Sub(startTime)).Nanoseconds() / 1000000))

	dh := DeferHTTP{
		Path: path,
		Time: t,
	}

	curlist = append(curlist, dh)
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

				appendHTTP(startTime, r.URL.Path)
			}
		}()

		f(w, r)

		appendHTTP(startTime, r.URL.Path)
	}
}
