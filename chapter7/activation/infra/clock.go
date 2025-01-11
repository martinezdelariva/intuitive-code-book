package infra

import "time"

type ClockSystem struct{}

func (c ClockSystem) Now() int64 {
	return time.Now().UTC().Unix()
}
