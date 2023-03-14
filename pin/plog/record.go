package plog

import (
	"github.com/ethanvc/pin/pin/kvrepo"
	"time"
)

type Record struct {
	Time  time.Time
	Event string
	Level Level
	Pc    uintptr
	repo  kvrepo.KvRepo
	l     *Logger
}

func (r *Record) Done() {
	if r == nil {
		return
	}
	for _, h := range r.l.Handlers {
		h(r.l, *r)
	}
}

func (r *Record) Int64(k string, v int64) *Record {
	if r == nil {
		return nil
	}
	r.repo.AddKvs(kvrepo.Kv{
		Key:   k,
		Value: kvrepo.Int64Value(v),
	})
	return r
}

func (r *Record) Str(k string, v string) *Record {
	if r == nil {
		return nil
	}
	r.repo.AddKvs(kvrepo.Kv{
		Key:   k,
		Value: kvrepo.StringValue(v),
	})
	return r
}

func (r *Record) Any(k string, v any) *Record {
	if r == nil {
		return nil
	}
	return r
}
