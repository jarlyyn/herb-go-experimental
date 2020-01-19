package logger

import (
	"time"
)

type Prefix interface {
	NewPrefix() []byte
}

type TimePrefix struct {
	Layout string
}

func (p *TimePrefix) NewPrefix() []byte {
	l := p.Layout
	if l == "" {
		l = DefaultTimeLayout
	}
	return []byte(time.Now().Format(l) + " ")
}

var DefaultTimePrefix = &TimePrefix{}

var DefaultTimeLayout = "2006-01-02 03:04:05"

type FixedPrefix string

func (p FixedPrefix) NewPrefix() []byte {
	return []byte(string(p) + " ")
}
