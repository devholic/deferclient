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
	api_url = "https://api.deferpanic.com/v1/panics/create"
)

var Token string
var NoPost = false

type DeferJSON struct {
	Msg       string `json:"ErrorName"`
	BackTrace string `json:"Body"`
	GoVersion string `json:"Version"`
}

func Persist() {
	if err := recover(); err != nil {
		prep(err)
	}
}

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

func ShipTrace(exception string, errorstr string) {
	if NoPost {
		return
	}

	go_version := runtime.Version()

	body := cleanTrace(exception)

	dj := &DeferJSON{Msg: errorstr, BackTrace: body, GoVersion: go_version}
	b, err := json.Marshal(dj)

	req, err := http.NewRequest("POST", api_url, bytes.NewBuffer(b))
	req.Header.Set("X-deferid", Token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

}
