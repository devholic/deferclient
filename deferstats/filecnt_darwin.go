package deferstats

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
)

// openFileCnt returns the number of open files in this process
func openFileCnt() int {
	out, err := exec.Command("/bin/sh", "-c", fmt.Sprintf("lsof -p %v", os.Getpid())).Output()
	if err != nil {
		log.Println("having trouble getting the open file count" + err.Error())
		return 0
	}
	return bytes.Count(out, []byte("\n"))
}
