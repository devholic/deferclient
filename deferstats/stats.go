package deferstats

import (
	"encoding/json"
	"github.com/deferpanic/deferclient/deferclient"
	"log"
	"runtime"
	"runtime/debug"
	"strconv"
	"time"
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
)

// Token is your deferpanic token available in settings
var Token string
var Verbose bool = false

// DeferHTTP holds the path uri and latency for each request
type DeferHTTP struct {
	Path       string `json:"Uri"`
	StatusCode int    `json:"StatusCode"`
	Time       int    `json:"Time"`
}

// DeferDB holds the query and latency for each sql query whose
// threshold was overran
type DeferDB struct {
	Query string `json:"Query"`
	Time  int    `json:"Time"`
}

// DeferStats captures {mem, gc, goroutines and http calls}
type DeferStats struct {
	Mem        string      `json:"Mem"`
	GC         string      `json:"GC"`
	GoRoutines string      `json:"GoRoutines"`
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

	ds := DeferStats{
		Mem:        mems,
		GoRoutines: grs,
		HTTPs:      curlist,
		DBs:        querylist,
		GC:         gcs,
	}

	// empty our https/dbs
	curlist = []DeferHTTP{}
	querylist = []DeferDB{}

	go func() {
		b, err := json.Marshal(ds)
		if err != nil {
			log.Println(err)
		}

		// hack
		deferclient.Token = Token

		deferclient.PostIt(b, statsUrl)
	}()
}
