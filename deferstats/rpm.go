package deferstats

import (
	"sync"
)

var rpms = rpmSet{}

type rpmSet struct {
	lock sync.RWMutex
	rpm  Rpm
}

// HTTPRpm holds the count of each HTTP status code during a stats
// collection interval
type Rpm struct {
	StatusOk int `json:"200,omitempty"`

	StatusMovedPermanently int `json:"301,omitempty"`
	StatusFound            int `json:"302,omitempty"`

	StatusBadRequest   int `json:"400,omitempty"`
	StatusUnauthorized int `json:"401,omitempty"`
	StatusForbidden    int `json:"403,omitempty"`
	StatusNotFound     int `json:"404,omitempty"`

	StatusInternalServerError int `json:"500,omitempty"`
	StatusServiceUnavailable  int `json:"503,omitempty"`
}

// ResetRPM clobbers old rpmset
func ResetRPM() {
	rpms = rpmSet{}
}

func (r *rpmSet) List() Rpm {
	r.lock.Lock()
	defer r.lock.Unlock()

	rpm := r.rpm
	return rpm
}

func (r *rpmSet) Inc(code int) {
	r.lock.Lock()
	defer r.lock.Unlock()

	switch code {
	case 200:
		r.rpm.StatusOk += 1
	case 301:
		r.rpm.StatusMovedPermanently += 1
	case 302:
		r.rpm.StatusFound += 1
	case 400:
		r.rpm.StatusBadRequest += 1
	case 401:
		r.rpm.StatusUnauthorized += 1
	case 403:
		r.rpm.StatusForbidden += 1
	case 404:
		r.rpm.StatusNotFound += 1
	case 500:
		r.rpm.StatusInternalServerError += 1
	case 503:
		r.rpm.StatusServiceUnavailable += 1
	}

}
