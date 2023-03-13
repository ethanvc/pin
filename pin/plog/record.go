package plog

import (
	"github.com/ethanvc/pin/pin/attrrepo"
	"time"
)

type Record struct {
	Time  time.Time
	Event string
	Level Level
	Pc    uintptr
	repo  attrrepo.AttrRepo
	l     *Logger
}

func (r *Record) Done() {
	for _, h := range r.l.handlers {
		h(r.l, *r)
	}
}

func (r *Record) Int64(k string, v int64) *Record {
	if r == nil {
		return nil
	}
	return r
}

func (r *Record) Str(k string, v string) *Record {
	if r == nil {
		return nil
	}
	return r
}

func (r *Record) Any(k string, v any) *Record {
	if r == nil {
		return nil
	}
	return r
}
