package deferclient

import (
	"syscall"
	"unsafe"
)

// SetTotal returns the number of bytes from hw.memsize on osx
func (m *Mem) SetTotal() {

	var data interface{}
	data = &m.Total

	val, err := syscall.Sysctl("hw.memsize")
	if err != nil {
		m.Total = 0
	}

	buf := []byte(val)

	v := data.(*uint64)
	*v = *(*uint64)(unsafe.Pointer(&buf[0]))
}
