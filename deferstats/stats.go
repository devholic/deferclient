package deferstats

import (
	"runtime/debug"

	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"runtime"
	"strconv"
	"time"
)

// fixme
const (
	// statsFrequency controls how often to report into deferpanic
	statsFrequency = 60 * time.Second
	// apiUrl is the stats api endpoint
	apiUrl = "https://api.deferpanic.com/v1/stats/create"
)

// Token is your deferpanic token available in settings
var Token string

// DeferHTTP holds the path uri and latency for each request
type DeferHTTP struct {
	Path string `json:Uri"`
	Time int    `json:Time"`
}

// DeferDB holds the query and latency for each sql query whose
// threshold was overran
type DeferDB struct {
	Query string `json:Query"`
	Time  int    `json:Time"`
}

// DeferStats captures {mem, gc, goroutines and http calls}
type DeferStats struct {
	Mem        string      `json:Mem"`
	GC         string      `json:GC"`
	GoRoutines string      `json:"GoRoutines"`
	HTTPs      []DeferHTTP `json:"HTTPs"`
	DBs        []DeferDB   `json:"DBs"`
}

// CaptureStats POSTs DeferStats every
func CaptureStats() {
	for {
		go capture()
		// set me to a reasonable timeout
		time.Sleep(statsFrequency)
	}
}

// ShipStats sends DeferStats to deferpanic
func ShipStats(stats DeferStats) {
	b, err := json.Marshal(stats)

	req, err := http.NewRequest("POST", apiUrl, bytes.NewBuffer(b))
	req.Header.Set("X-deferid", Token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()

}

// capture does a one time collection of DeferStats
func capture() {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)

	var gc debug.GCStats
	debug.ReadGCStats(&gc)

	ds := DeferStats{
		Mem:        strconv.FormatUint(mem.Alloc, 10),
		GoRoutines: strconv.Itoa(runtime.NumGoroutine()),
		HTTPs:      curlist,
		DBs:        querylist,
		GC:         strconv.FormatInt(gc.NumGC, 10),
	}

	// empty our https/dbs
	curlist = []DeferHTTP{}
	querylist = []DeferDB{}

	go ShipStats(ds)

}
