package deferstats

import (
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
	statsFrequency = 250 * time.Millisecond
	apiUrl         = "http://127.0.0.1:8080/v1/stats/create"

	// apiUrl = "https://api.deferpanic.com/v1/panics/create"
)

var Token string

type DeferHTTP struct {
	Uri  string `json:Uri"`
	Time int    `json:Time"`
}

type DeferStats struct {
	Mem        string      `json:Mem"`
	GC         string      `json:GC"`
	GoRoutines string      `json:"GoRoutines"`
	HTTPs      []DeferHTTP `json:"HTTPs"`
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

	log.Println("post")

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
	log.Println("posting")

	ds := DeferStats{
		Mem:        strconv.FormatUint(mem.Alloc, 10),
		GoRoutines: strconv.Itoa(runtime.NumGoroutine()),
	}

	ShipStats(ds)

}
