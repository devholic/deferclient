// +build go1.5

package deferclient

import (
	"bytes"
	"encoding/json"
	"hash/crc32"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime/trace"
	"time"
)

// MakeTrace POST Trace binaries to the deferpanic website
func (c *DeferPanicClient) MakeTrace(commandId int, agent *Agent) {
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

	log.Println("trace started")
	err := trace.Start(buffer)
	if err != nil {
		log.Println(err)
		return
	}

	select {
	case <-time.After(30 * time.Second):
		trace.Stop()
		log.Println("trace finished")

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
		t := NewTrace(out, pkg, commandId, false)

		b, err := json.Marshal(t)
		if err != nil {
			log.Println(err)
			return
		}

		c.Postit(b, traceUrl, false)
	}
}
