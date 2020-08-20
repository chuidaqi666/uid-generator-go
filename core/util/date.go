package util

import (
	"time"
)

func DateToSecond(date string) uint64 {
	c := date + " 00:00:00"
	loc, _ := time.LoadLocation("Local")
	theTime, _ := time.ParseInLocation("2006-01-02 15:04:05", c, loc)
	return uint64(theTime.Unix())
}
