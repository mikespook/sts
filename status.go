package sts

import (
	"sync"
	"time"
)

const (
	StatusConn    = "CONN"
	StatusRead    = "READ"
	StatusWritten = "WRITTEN"
)

type Status struct {
	sync.RWMutex
	start time.Time
	data  map[string]uint64
}

func NewStatus() *Status {
	return &Status{
		start: time.Now(),
		data:  make(map[string]uint64),
	}
}

func (status *Status) Established() time.Duration {
	return time.Now().Sub(status.start)
}

func (status *Status) Inc(t string, delta uint64) uint64 {
	status.Lock()
	defer status.Unlock()
	status.data[t] += delta
	return status.data[t]
}

func (status *Status) Value(t string) uint64 {
	status.RLock()
	defer status.RUnlock()
	return status.data[t]
}
