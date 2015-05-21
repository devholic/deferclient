// Package deferstats implements deferpanic stats
package deferstats

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"runtime/debug"
	"strconv"
	"time"

	"github.com/deferpanic/deferclient/deferclient"
)

// being DEPRECATED
//
// please use deferstats.NewClient(token)
const (
	// statsFrequency controls how often to report into deferpanic in seconds
	statsFrequency = 60

	// statsUrl is the stats api endpoint
	statsUrl = deferclient.ApiBase + "/stats/create"
)

// being DEPRECATED
//
// please use deferstats.NewClient(token)
var (
	// Token is your deferpanic token available in settings
	Token string
)

// DeferStats captures {mem, gc, goroutines and http calls}
type DeferStats struct {
	Mem        string      `json:"Mem"`
	GC         string      `json:"GC"`
	GoRoutines string      `json:"GoRoutines"`
	Cgos       string      `json:"Cgos"`
	Fds        string      `json:"Fds"`
	HTTPs      []DeferHTTP `json:"HTTPs"`
	DBs        []DeferDB   `json:"DBs"`
}

// Client is the client for making metrics requests to the
// defer panic api
type Client struct {
	// statsFrequency controls how often to report into deferpanic in seconds
	statsFrequency int

	// statsUrl is the stats api endpoint
	statsUrl string

	// GrabGC determines if we should grab gc stats
	GrabGC bool

	// GrabMem determines if we should grab mem stats
	GrabMem bool

	// GrabGR determines if we should grab go routine stats
	GrabGR bool

	// GrabCgo determines if we should grab cgo count
	GrabCgo bool

	// GrabFd determines if we should grab fd count
	GrabFd bool

	// Token is your deferpanic token available in settings
	Token string

	// Verbose logs any stats output - for testing/dev
	Verbose bool

	// Environment sets an environment tag to differentiate between separate
	// environments - default is production.
	environment string

	// AppGroup sets an optional tag to differentiate between your various
	// services - default is default
	appGroup string

	// noPost when set to true disables reporting to deferpanic - useful
	// for dev/test envs
	noPost bool

	// BaseClient is the base deferpanic client that all http requests use
	BaseClient *deferclient.DeferPanicClient
}

// NewClient instantiates and returns a new client
func NewClient(token string) *Client {

	ds := &Client{
		statsFrequency: 60,
		statsUrl:       deferclient.ApiBase + "/stats/create",
		GrabGC:         true,
		GrabMem:        true,
		GrabGR:         true,
		GrabCgo:        true,
		GrabFd:         true,
		Verbose:        false,
		Token:          token,
		environment:    "production",
		appGroup:       "default",
		noPost:         false,
	}

	ds.BaseClient = deferclient.NewDeferPanicClient(token)
	ds.BaseClient.Environment = ds.environment
	ds.BaseClient.AppGroup = ds.appGroup
	ds.BaseClient.NoPost = ds.noPost

	return ds
}

// Setenvironment sets the environment
// default is 'production'
func (c *Client) Setenvironment(environment string) {
	c.environment = environment
	c.BaseClient.Environment = c.environment
}

// SetappGroup sets the app group
// default is 'default'
func (c *Client) SetappGroup(appGroup string) {
	c.appGroup = appGroup
	c.BaseClient.AppGroup = c.appGroup
}

// Setnopost disables reporting to deferpanic
// default is false
func (c *Client) SetnoPost(noPost bool) {
	c.noPost = noPost
	c.BaseClient.NoPost = c.noPost
}

// CaptureStats POSTs DeferStats every statsFrequency
func (c *Client) CaptureStats() {
	if !c.noPost {
		c.updateAgent()
	}

	tickerChannel := time.Tick(time.Duration(c.statsFrequency) * time.Second)
	for tc := range tickerChannel {

		// Capture the stats every statsFrequency seconds
		go c.capture()

		if c.Verbose {
			log.Printf("Captured at:%v\n", tc)
		}
	}
}

// updateAgent sets the agent details
func (c *Client) updateAgent() {
	if c.noPost {
		return
	}

	b, err := json.Marshal(c.BaseClient.Agent)
	if err != nil {
		log.Println(err)
	}

	agentUrl := deferclient.ApiBase + "/agent_ids/create"

	c.BaseClient.Postit(b, agentUrl)
}

// capture does a one time collection of DeferStats
func (c *Client) capture() {

	var mem runtime.MemStats
	var gc debug.GCStats

	mems := ""
	if c.GrabMem {
		runtime.ReadMemStats(&mem)
		mems = strconv.FormatUint(mem.Alloc, 10)
	}

	gcs := ""
	if c.GrabGC {
		debug.ReadGCStats(&gc)
		gcs = strconv.FormatInt(gc.NumGC, 10)
	}

	grs := ""
	if c.GrabGR {
		grs = strconv.Itoa(runtime.NumGoroutine())
	}

	cgos := ""
	if c.GrabCgo {
		cgos = strconv.FormatInt(runtime.NumCgoCall(), 10)
	}

	fds := ""
	if c.GrabFd {
		fds = strconv.Itoa(openFileCnt())
	}

	ds := DeferStats{
		Mem:        mems,
		GoRoutines: grs,
		Cgos:       cgos,
		Fds:        fds,
		HTTPs:      curlist.List(),
		DBs:        Querylist.List(),
		GC:         gcs,
	}

	// FIXME
	// empty our https/dbs
	curlist.Reset()
	Querylist.Reset()

	go func() {
		b, err := json.Marshal(ds)
		if err != nil {
			log.Println(err)
		}

		c.BaseClient.Postit(b, c.statsUrl)
	}()
}

// openFileCnt returns the number of open files in this process
func openFileCnt() int {
	out, err := exec.Command("/bin/sh", "-c", fmt.Sprintf("lsof -p %v", os.Getpid())).Output()
	if err != nil {
		log.Println(err.Error())
		return 0
	}
	return bytes.Count(out, []byte("\n"))
}
