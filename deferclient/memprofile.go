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
)

// MemProfile contains information about this client's memory profile and its producing package
type MemProfile struct {
	Out       []byte `json:"Out,omitempty"`
	Pkg       []byte `json:"Pkg,omitempty"`
	CommandId int    `json:"CommandId"`
	Ignored   bool   `json:"Ignored"`
}

// NewMemProfile instantitates and returns a new memory profile
// it is meant to be called once after the completing application memory profiling
func NewMemProfile(out []byte, pkg []byte, commandid int, ignored bool) *MemProfile {
	c := &MemProfile{
		Out:       out,
		Pkg:       pkg,
		CommandId: commandid,
		Ignored:   ignored,
	}

	return c
}

// MakeMemProfile POST MemProfile binaries to the deferpanic website
func (c *DeferPanicClient) MakeMemProfile(commandId int, agent *Agent) {
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

	log.Println("mem profile started")
	pprof.Lookup("heap").WriteTo(buffer, 0)
	log.Println("mem profile finished")

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
	t := NewMemProfile(out, pkg, commandId, false)

	b, err := json.Marshal(t)
	if err != nil {
		log.Println(err)
		return
	}

	c.Postit(b, memprofileUrl, false)
}
