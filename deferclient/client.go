// Package deferclient implements access to the deferpanic api.
package deferclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"hash/crc32"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"runtime/trace"
	"strings"
	"sync"
	"time"
)

const (
	// ApiVersion is the version of this client
	ApiVersion = "v1.15"

	// ApiBase is the base url that client requests goto
	//	ApiBase = "https://api.deferpanic.com/" + ApiVersion
	ApiBase = "http://localhost:8080/" + ApiVersion

	// UserAgent is the User Agent that is used with this client
	UserAgent = "deferclient " + ApiVersion

	// errorsUrl is the url to post panics && errors to
	errorsUrl = ApiBase + "/panics/create"

	// traceUrl is the url to post traces to
	traceUrl = ApiBase + "/uploads/trace/create"
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

	RunningCommands map[int]bool
	sync.Mutex
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
		Token:           token,
		UserAgent:       "deferclient " + ApiVersion,
		Agent:           a,
		PrintPanics:     false,
		NoPost:          false,
		RunningCommands: make(map[int]bool),
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

// PersistRepanic ensures any panics will post to deferpanic website for
// tracking, it also reissues the panic afterwards.
// typically used in non http go-routines
func (c *DeferPanicClient) PersistRepanic() {
	if err := recover(); err != nil {
		c.PrepSync(err, 0)
		panic(err)
	}
}

// Prep takes an error && a spanId
// it cleans up the error/trace before calling ShipTrace
// if spanId is zero it is ommited
func (c *DeferPanicClient) Prep(err interface{}, spanId int64) {
	c.prep(err, spanId, false)
}

// PrepSync takes an error && a spanId
// it cleans up the error/trace before calling ShipTrace
// waits for ShipTrace, in a go routine, to complete before continuing
// if spanId is zero it is ommited
func (c *DeferPanicClient) PrepSync(err interface{}, spanId int64) {
	c.prep(err, spanId, true)
}

// prep is an internal function that can be called to synchronize after
// shipping the the trace to ensure completion.
func (c *DeferPanicClient) prep(err interface{}, spanId int64, syncShipTrace bool) {
	errorMsg := fmt.Sprintf("%q", err)

	errorMsg = strings.Replace(errorMsg, "\"", "", -1)

	if c.PrintPanics {
		stack := string(debug.Stack())
		fmt.Println(stack)
	}

	body := backTrace()

	if syncShipTrace {
		done := make(chan bool)
		go func() {
			c.ShipTrace(body, errorMsg, spanId)
			done <- true
		}()
		<-done
	} else {
		go c.ShipTrace(body, errorMsg, spanId)
	}
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

	c.Postit(b, errorsUrl, false)
}

// Postit Posts an API request w/b body to url and sets appropriate
// headers
func (c *DeferPanicClient) Postit(b []byte, url string, analyseResponse bool) {
	defer func() {
		if rec := recover(); rec != nil {
			err := fmt.Sprintf("%q", rec)
			log.Println(err)
		}
	}()

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

	if analyseResponse {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println(err)
			return
		}

		var commands []Command
		err = json.Unmarshal(body, &commands)
		if err != nil {
			log.Println(err)
			return
		}

		for _, command := range commands {
			c.Lock()
			running := c.RunningCommands[command.Id]
			c.Unlock()
			if !running {
				if command.GenerateTrace {
					go c.MakeTrace(command.Id)
				}
			}
		}
	}

}

// MakeTrace POST a Trace html to the deferpanic website
func (c *DeferPanicClient) MakeTrace(commandId int) {
	var buf []byte
	buffer := bytes.NewBuffer(buf)

	c.Lock()
	c.RunningCommands[commandId] = true
	c.Unlock()
	defer func() {
		c.Lock()
		delete(c.RunningCommands, commandId)
		c.Unlock()
	}()

	log.Println("trace started")
	err := trace.Start(buffer)
	if err != nil {
		log.Println(err)
		return
	}

	select {
	case <-time.After(30 * time.Second):
		trace.Stop()
		log.Println("trace finished")

		out := make([]byte, len(buffer.Bytes()))
		copy(out, buffer.Bytes())
		pkgpath, err := filepath.Abs(os.Args[0])
		if err != nil {
			log.Println(err)
			return
		}
		pkg, err := ioutil.ReadFile(pkgpath)
		if err != nil {
			log.Println(err)
			return
		}
		crc32 := crc32.ChecksumIEEE(pkg)
		size := int64(len(pkg))
		t := NewTrace(out, pkg, crc32, size, commandId)

		b, err := json.Marshal(t)
		if err != nil {
			log.Println(err)
			return
		}

		c.Postit(b, traceUrl, false)
	}
}
