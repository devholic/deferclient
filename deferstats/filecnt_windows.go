package deferstats

import (
	"log"
	"syscall"
	"unsafe"
)

var (
	kernel32              = syscall.MustLoadDLL("kernel32.dll")
	getProcessHandleCount = kernel32.MustFindProc("GetProcessHandleCount")
)

// openFileCnt returns the number of open files in this process
func openFileCnt() int {
	cp, err := syscall.GetCurrentProcess()
	if err != nil {
		log.Println("getcurrentProcess: %v\n", err)
	}

	var c uint32
	r, _, err := getProcessHandleCount.Call(uintptr(cp), uintptr(unsafe.Pointer(&c)))
	if r == 0 {
		log.Println(err)
	}
	return int(c)
}
