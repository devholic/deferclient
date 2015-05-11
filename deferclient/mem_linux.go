package deferclient

import (
	"io/ioutil"
	"strconv"
	"strings"
)

// SetTotal returns the number of bytes from /proc/meminfo on linux
func (m *Mem) SetTotal() {

	body, err := ioutil.ReadFile("/proc/meminfo")
	if err != nil {
		m.Total = 0
		return
	}

	s := string(body)

	// sir hack-a-lot
	fields := strings.Split(s, "\n")
	fz := strings.Split(fields[0], ":")

	num := strings.TrimLeft(fz[1], " ")
	onum := strings.Split(num, " ")

	i, err := strconv.ParseInt(onum[0], 10, 64)
	if err != nil {
		m.Total = 0
		return
	}

	m.Total = uint64(i * 1024)
}
