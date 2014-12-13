// Package deferclient implements access to the deferpanic api.
package deferclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
	"strings"
)

const (
	apiUrl = "https://api.deferpanic.com/v1/panics/create"
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
		prep(err)
	}
}

// recovers from http handler panics and posts to deferpanic website for
// tracking
func PanicRecover(f func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				prep(err)
			}
		}()
		f(w, r)
	}
}

func prep(err interface{}) {
	errorMsg := fmt.Sprintf("%q", err)

	errorMsg = strings.Replace(errorMsg, "\"", "", -1)

	buf := make([]byte, 1<<16)
	runtime.Stack(buf, false)
	sz := len(buf) - 1
	body := string(buf[:sz])

	ShipTrace(body, errorMsg)
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

	req, err := http.NewRequest("POST", apiUrl, bytes.NewBuffer(b))
	req.Header.Set("X-deferid", Token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

}
