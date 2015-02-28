package iface

import "time"

type Stat interface {
	ETime() time.Time
	Aggregate(key string) int
}
