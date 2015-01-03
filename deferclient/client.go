// Package deferclient implements access to the deferpanic api.
package deferclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strings"
)

const (
	// ApiVersion is the version of this client
	ApiVersion = "v1"

	// ApiBase is the base url that client requests goto
	ApiBase = "https://api.deferpanic.com/" + ApiVersion

	// UserAgent is the User Agent that is used with this client
	UserAgent = "deferclient " + ApiVersion

	// errorsUrl is the url to post urls to
	errorsUrl = ApiBase + "/panics/create"
)

// Your deferpanic client token
var Token string

// Bool that turns off tracking of errors and panics - useful for
// dev/test environments
var NoPost = false

// struct that holds expected json body for POSTing to deferpanic API v1
type DeferJSON struct {
	Msg       string `json:"ErrorName"`
	BackTrace string `json:"Body"`
	GoVersion string `json:"Version"`
}

// Persists ensures any panics will post to deferpanic website for
// tracking
func Persist() {
	if err := recover(); err != nil {
		Prep(err)
	}
}

// recovers from http handler panics and posts to deferpanic website for
// tracking
func PanicRecover(f func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				Prep(err)
			}
		}()
		f(w, r)
	}
}

// Prep cleans up the trace before posting
func Prep(err interface{}) {
	errorMsg := fmt.Sprintf("%q", err)

	errorMsg = strings.Replace(errorMsg, "\"", "", -1)

	body := ""
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

	go ShipTrace(body, errorMsg)
}

// encoding
func cleanTrace(body string) string {
	body = strings.Replace(body, "\n", "\\n", -1)
	body = strings.Replace(body, "\t", "\\t", -1)
	body = strings.Replace(body, "\x00", " ", -1)
	body = strings.TrimSpace(body)

	return body
}

// ShipTrace POSTs a DeferJSON json body to the deferpanic website
func ShipTrace(exception string, errorstr string) {
	if NoPost {
		return
	}

	goVersion := runtime.Version()

	body := cleanTrace(exception)

	dj := &DeferJSON{Msg: errorstr, BackTrace: body, GoVersion: goVersion}
	b, err := json.Marshal(dj)
	if err != nil {
		log.Println(err)
	}

	PostIt(b, errorsUrl)
}

func PostIt(b []byte, url string) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(b))
	req.Header.Set("X-deferid", Token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", UserAgent)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()
}
