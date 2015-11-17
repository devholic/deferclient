package deferclient

import (
	"bytes"
	"encoding/json"
	"hash/crc32"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime/pprof"
	"time"
)

// CPUProfile contains information about this client's cpu profile and its producing package
type CPUProfile struct {
	Out       []byte `json:"Out,omitempty"`
	Pkg       []byte `json:"Pkg,omitempty"`
	CommandId int    `json:"CommandId"`
	Ignored   bool   `json:"Ignored"`
}

// NewCPUProfile instantitates and returns a new cpu profile
// it is meant to be called once after the completing application cpu profiling
func NewCPUProfile(out []byte, pkg []byte, commandid int, ignored bool) *CPUProfile {
	c := &CPUProfile{
		Out:       out,
		Pkg:       pkg,
		CommandId: commandid,
		Ignored:   ignored,
	}

	return c
}

// MakeCPUProfile POST CPUProfile binaries to the deferpanic website
func (c *DeferPanicClient) MakeCPUProfile(commandId int, agent *Agent) {
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

	log.Println("cpu profile started")
	err := pprof.StartCPUProfile(buffer)
	if err != nil {
		log.Println(err)
		return
	}

	select {
	case <-time.After(30 * time.Second):
		pprof.StopCPUProfile()
		log.Println("cpu profile finished")

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
		if agent.CRC32 == crc32 && agent.Size == size {
			pkg = []byte{}
		}
		t := NewCPUProfile(out, pkg, commandId, false)

		b, err := json.Marshal(t)
		if err != nil {
			log.Println(err)
			return
		}

		c.Postit(b, cpuprofileUrl, false)
	}
}
