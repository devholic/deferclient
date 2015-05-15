// Package deferclient implements access to the deferpanic api.
package deferclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"runtime/debug"
	"strings"
)

// being DEPRECATED
const (
	// ApiVersion is the version of this client
	ApiVersion = "v1.8"

	// ApiBase is the base url that client requests goto
	ApiBase = "https://api.deferpanic.com/" + ApiVersion

	// UserAgent is the User Agent that is used with this client
	UserAgent = "deferclient " + ApiVersion

	// errorsUrl is the url to post panics && errors to
	errorsUrl = ApiBase + "/panics/create"
)

// being DEPRECATED
var (
	// Your deferpanic client token
	// this is being DEPRECATED
	Token string

	// Bool that turns off tracking of errors and panics - useful for
	// dev/test environments
	// this is being DEPRECATED
	NoPost = false

	// PrintPanics controls whether or not the HTTPHandler function prints
	// recovered panics. It is disabled by default.
	// this is being DEPRECATED
	PrintPanics = false

	// Environment sets an environment tag to differentiate between separate
	// environments - default is production.
	// this is being DEPRECATED
	Environment = "production"

	// AppGroup sets an optional tag to differentiate between your various
	// services - default is default
	// this is being DEPRECATED
	AppGroup = "default"
)

// DeferPanicClient is the base struct for making requests to the defer
// panic api
//
// FIXME: move all globals for future api bump
type DeferPanicClient struct {
	Token       string
	UserAgent   string
	Environment string
	AppGroup    string

	Agent       *Agent
	NoPost      bool
	PrintPanics bool
}

// struct that holds expected json body for POSTing to deferpanic API
type DeferJSON struct {
	Msg       string `json:"ErrorName"`
	BackTrace string `json:"Body"`
	SpanId    int64  `json:"SpanId,omitempty"`
}

// NewDeferPanicClient instantiates and returns a new deferpanic client
func NewDeferPanicClient(token string) *DeferPanicClient {
	a := NewAgent()

	dc := &DeferPanicClient{
		Token:       token,
		UserAgent:   "deferclient " + ApiVersion,
		Agent:       a,
		PrintPanics: false,
		NoPost:      false,
	}

	return dc
}

// Persists ensures any panics will post to deferpanic website for
// tracking
// typically used in non http go-routines
func (c *DeferPanicClient) Persist() {
	if err := recover(); err != nil {
		c.Prep(err, 0)
	}
}

// Prep takes an error && a spanId
// it cleans up the error/trace before calling ShipTrace
// if spanId is zero it is ommited
func (c *DeferPanicClient) Prep(err interface{}, spanId int64) {
	errorMsg := fmt.Sprintf("%q", err)

	errorMsg = strings.Replace(errorMsg, "\"", "", -1)

	if c.PrintPanics {
		stack := string(debug.Stack())
		fmt.Println(stack)
	}

	body := backTrace()

	go c.ShipTrace(body, errorMsg, spanId)
}

// backtrace grabs the backtrace
func backTrace() (body string) {

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

// cleanTrace should be fixed
// encoding
func cleanTrace(body string) string {
	body = strings.Replace(body, "\n", "\\n", -1)
	body = strings.Replace(body, "\t", "\\t", -1)
	body = strings.Replace(body, "\x00", " ", -1)
	body = strings.TrimSpace(body)

	return body
}

// ShipTrace POSTs a DeferJSON json body to the deferpanic website
// if spanId is zero it is ignored
func (c *DeferPanicClient) ShipTrace(exception string, errorstr string, spanId int64) {
	if c.NoPost {
		return
	}

	body := cleanTrace(exception)

	dj := &DeferJSON{
		Msg:       errorstr,
		BackTrace: body,
	}

	if spanId > 0 {
		dj.SpanId = spanId
	}

	b, err := json.Marshal(dj)
	if err != nil {
		log.Println(err)
	}

	c.Postit(b, errorsUrl)
}

// Postit Posts an API request w/b body to url and sets appropriate
// headers
func (c *DeferPanicClient) Postit(b []byte, url string) {
	if c.NoPost {
		return
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(b))

	req.Header.Set("X-deferid", c.Token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", c.UserAgent)
	req.Header.Set("X-dpenv", c.Environment)
	req.Header.Set("X-dpgroup", c.AppGroup)
	req.Header.Set("X-dpagentid", c.Agent.Name)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case 401:
		log.Println("wrong or invalid API token")
	case 429:
		log.Println("too many requests - you are being rate limited")
	case 503:
		log.Println("service not available")
	default:
	}

}
