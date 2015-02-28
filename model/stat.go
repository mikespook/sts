package model

import "time"

type Stat struct {
	ETime    time.Time
	Sessions int
	Agents   int
}
