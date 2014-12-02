package stats

import (
	"fmt"
	"runtime"
	"time"
)

// fixme
const (
	apiUrl = "http://127.0.0.1:8080/v1/stats/create"

	// apiUrl = "https://api.deferpanic.com/v1/panics/create"
)

type DeferStats struct {
	Mem        string `json:Mem"`
	GC         string `json:GC"`
	GoRoutines string `json:"GoRoutines"`
}

func ShipStats(stats DeferStats) {
	goVersion := runtime.Version()

	b, err := json.Marshal(stats)

	req, err := http.NewRequest("POST", apiUrl, bytes.NewBuffer(b))
	req.Header.Set("X-deferid", Token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

}

func captureMem() {
	var mem runtime.MemStats

	runtime.ReadMemStats(&mem)

	ds := DeferStats{
		Mem: mem.Alloc,
	}

	ShipStats(ds)

}
