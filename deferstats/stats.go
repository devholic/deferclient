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

// fixme
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

// Token is your deferpanic token available in settings
var Token string

// Verbose logs any stats output - for testing/dev
var Verbose bool = false

// Environment sets an environment tag to differentiate between separate
// environments - default is production.
var Environment = "production"

// AppGroup sets an optional tag to differentiate between your various
// services - default is default
var AppGroup = "default"

// DeferStats captures {mem, gc, goroutines and http calls}
type DeferStats struct {
	Mem        string      `json:"Mem"`
	GC         string      `json:"GC"`
	GoRoutines string      `json:"GoRoutines"`
	Cgos       string      `json:"Cgos"`
	HTTPs      []DeferHTTP `json:"HTTPs"`
	DBs        []DeferDB   `json:"DBs"`
}

// CaptureStats POSTs DeferStats every
func CaptureStats() {

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
