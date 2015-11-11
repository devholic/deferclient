package deferclient

import (
	"net"
	"os"
	"runtime"
	"strconv"
)

// Mem simply grabs the total memory avail on the system
type Mem struct {
	Total uint64
}

// Agent contains information about this client's agent
type Agent struct {
	Name       string `json:"Name"`
	Cpucores   int    `json:"Cpucores"`
	Goarch     string `json:"goarch"`
	Goos       string `json:"goos"`
	Totalmem   uint64 `json:"totalmem"`
	Govers     string `json:"govers"`
	ApiVersion string `json:"ApiVersion"`
	CRC32      uint32 `json:"CRC32"`
	Size       int64  `json:"Size"`
}

// SetName sets a 'unique'ish id for this agent
func (a *Agent) SetName() {

	local := "bad"

	host, _ := os.Hostname()
	addrs, _ := net.LookupIP(host)
	for _, addr := range addrs {
		if ipv4 := addr.To4(); ipv4 != nil {
			local = ipv4.String()
		}
	}

	pid := os.Getpid()

	a.Name = local + "-" + strconv.Itoa(pid)
}

// NewAgent instantitates and returns a new agent
// it is meant to be called once at the start of a new agent checking in
// (a new process)
func NewAgent() *Agent {

	m := Mem{}
	m.SetTotal()

	a := &Agent{
		Cpucores:   runtime.NumCPU(),
		Goarch:     runtime.GOARCH,
		Goos:       runtime.GOOS,
		Totalmem:   m.Total,
		Govers:     runtime.Version(),
		ApiVersion: ApiVersion,
		CRC32:      0,
		Size:       0,
	}

	a.SetName()

	return a
}
