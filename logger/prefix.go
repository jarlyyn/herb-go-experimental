package logger

import (
	"time"
)

type Prefix interface {
	NewPrefix() string
}

type TimePrefix struct {
	Layout string
}

func (p *TimePrefix) NewPrefix() string {
	return time.Now().Format(p.Layout) + " "
}

var DefaultTimePrefix = &TimePrefix{
	Layout: "2006-01-02 03:04:05",
}

type FixedPrefix string

func (p FixedPrefix) NewPrefix() string {
	return string(p) + " "
}
