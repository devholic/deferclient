// Package deferstats implements deferpanic stats
package deferstats

import (
	"encoding/json"
	"log"
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

	// GrabGC determines if we should grab gc stats
	GrabGC = true

	// GrabMem determines if we should grab mem stats
	GrabMem = true

	// GrabGR determines if we should grab go routine stats
	GrabGR = true

	// GrabCgo determines if we should grab cgo count
	GrabCgo = true
)

// being DEPRECATED
//
// please use deferstats.NewClient(token)
var (
	// Token is your deferpanic token available in settings
	Token string

	// Verbose logs any stats output - for testing/dev
	Verbose bool = false

	// Environment sets an environment tag to differentiate between separate
	// environments - default is production.
	Environment = "production"

	// AppGroup sets an optional tag to differentiate between your various
	// services - default is default
	AppGroup = "default"
)

// DeferStats captures {mem, gc, goroutines and http calls}
type DeferStats struct {
	Mem        string      `json:"Mem"`
	GC         string      `json:"GC"`
	GoRoutines string      `json:"GoRoutines"`
	Cgos       string      `json:"Cgos"`
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
		Verbose:        false,
		Token:          token,
		environment:    "production",
		appGroup:       "default",
	}

	ds.BaseClient = deferclient.NewDeferPanicClient(token)
	ds.BaseClient.Environment = ds.environment
	ds.BaseClient.AppGroup = ds.appGroup

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

// CaptureStats POSTs DeferStats every statsFrequency
func (c *Client) CaptureStats() {
	tickerChannel := time.Tick(time.Duration(c.statsFrequency) * time.Second)
	for tc := range tickerChannel {

		// Capture the stats every statsFrequency seconds
		go c.capture()

		if c.Verbose {
			log.Printf("Captured at:%v\n", tc)
		}
	}
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

	ds := DeferStats{
		Mem:        mems,
		GoRoutines: grs,
		Cgos:       cgos,
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

// CaptureStats POSTs DeferStats every statsFrequency
// being DEPRECATED
func CaptureStats() {
	log.Println("please consider using deferstats.NewClient(token)")

	tickerChannel := time.Tick(time.Duration(statsFrequency) * time.Second)
	for ts := range tickerChannel {

		// Capture the stats every X seconds
		go capture()

		if Verbose {
			log.Printf("Captured at:%v\n", ts)
		}
	}
}

// capture does a one time collection of DeferStats
// being DEPRECATED
func capture() {

	var mem runtime.MemStats
	var gc debug.GCStats

	mems := ""
	if GrabMem {
		runtime.ReadMemStats(&mem)
		mems = strconv.FormatUint(mem.Alloc, 10)
	}

	gcs := ""
	if GrabGC {
		debug.ReadGCStats(&gc)
		gcs = strconv.FormatInt(gc.NumGC, 10)
	}

	grs := ""
	if GrabGR {
		grs = strconv.Itoa(runtime.NumGoroutine())
	}

	cgos := ""
	if GrabCgo {
		cgos = strconv.FormatInt(runtime.NumCgoCall(), 10)
	}

	ds := DeferStats{
		Mem:        mems,
		GoRoutines: grs,
		Cgos:       cgos,
		HTTPs:      curlist.List(),
		DBs:        Querylist.List(),
		GC:         gcs,
	}

	// empty our https/dbs
	curlist.Reset()
	Querylist.Reset()

	go func() {
		b, err := json.Marshal(ds)
		if err != nil {
			log.Println(err)
		}

		// hack
		deferclient.Token = Token
		deferclient.Environment = Environment
		deferclient.AppGroup = AppGroup

		deferclient.PostIt(b, statsUrl)
	}()
}
