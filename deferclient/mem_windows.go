package deferclient

import (
	"encoding/binary"
	"log"
	"syscall"
	"unsafe"
)

var (
	kernel32, _             = syscall.LoadLibrary("kernel32.dll")
	globalMemoryStatusEx, _ = syscall.GetProcAddress(kernel32, "GlobalMemoryStatusEx")
)

// SetTotal returns the number of bytes on windows
func (m *Mem) SetTotal() {

	var memoryStatusEx [64]byte
	binary.LittleEndian.PutUint32(memoryStatusEx[:], 64)
	p := uintptr(unsafe.Pointer(&memoryStatusEx[0]))

	ret, _, callErr := syscall.Syscall(uintptr(globalMemoryStatusEx), 1, p, 0, 0)
	if ret == 0 {
		log.Println(callErr)
	}

	m.Total = binary.LittleEndian.Uint64(memoryStatusEx[8:])
}
