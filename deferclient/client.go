// Package deferclient implements access to the deferpanic api.
package deferclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"runtime"
	"runtime/debug"
	"strconv"
	"strings"
)

const (
	// ApiVersion is the version of this client
	ApiVersion = "v1.4"

	// ApiBase is the base url that client requests goto
	ApiBase = "https://api.deferpanic.com/" + ApiVersion

	// UserAgent is the User Agent that is used with this client
	UserAgent = "deferclient " + ApiVersion

	// errorsUrl is the url to post urls to
	errorsUrl = ApiBase + "/panics/create"
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
	AgentId     string
}

// Your deferpanic client token
// this is being DEPRECATED
var Token string

// Bool that turns off tracking of errors and panics - useful for
// dev/test environments
// this is being DEPRECATED
var NoPost = false

// PrintPanics controls whether or not the HTTPHandler function prints
// recovered panics. It is disabled by default.
// this is being DEPRECATED
var PrintPanics = false

// Environment sets an environment tag to differentiate between separate
// environments - default is production.
// this is being DEPRECATED
var Environment = "production"

// AppGroup sets an optional tag to differentiate between your various
// services - default is default
// this is being DEPRECATED
var AppGroup = "default"

// struct that holds expected json body for POSTing to deferpanic API v1
type DeferJSON struct {
	Msg       string `json:"ErrorName"`
	BackTrace string `json:"Body"`
	GoVersion string `json:"Version"`
}

// agentID sets a 'unique' ID for this agent
func agentID() string {

	local := "bad"

	host, _ := os.Hostname()
	addrs, _ := net.LookupIP(host)
	for _, addr := range addrs {
		if ipv4 := addr.To4(); ipv4 != nil {
			local = ipv4.String()
		}
	}

	pid := os.Getpid()

	return local + "-" + strconv.Itoa(pid)
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

	if PrintPanics {
		stack := string(debug.Stack())
		fmt.Println(stack)
	}

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

// postIt Posts an API request w/b body to url and sets appropriate
// headers
func (c *DeferPanicClient) Postit(b []byte, url string) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(b))

	req.Header.Set("X-deferid", c.Token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", c.UserAgent)
	req.Header.Set("X-dpenv", c.Environment)
	req.Header.Set("X-dpgroup", c.AppGroup)
	req.Header.Set("X-dpagentid", c.AgentId)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()

}

// PostIt Posts an API request w/b body to url and sets appropriate
// headers
// this is being DEPRECATED
func PostIt(b []byte, url string) {

	dpc := DeferPanicClient{
		Token:       Token,
		UserAgent:   UserAgent,
		Environment: Environment,
		AppGroup:    AppGroup,
		AgentId:     agentID(),
	}

	dpc.Postit(b, url)
}
