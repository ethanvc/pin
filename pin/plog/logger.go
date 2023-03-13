package plog

import (
	"github.com/ethanvc/pin/pin/attrrepo"
	"sync/atomic"
	"time"
)

type Logger struct {
	attrRepo attrrepo.AttrRepo
	handlers []Handler
	level    Level
}

var defaultLogger atomic.Pointer[Logger]

func init() {
	l := &Logger{
		level: LevelInfo,
	}
	defaultLogger.Store(l)
}

func Default() *Logger {
	return defaultLogger.Load()
}

func (this *Logger) Info(event string) *Record {
	return this.Log(LevelInfo, event)
}

func (this *Logger) Log(lvl Level, event string) *Record {
	if !this.Enabled(lvl) {
		return nil
	}
	r := &Record{
		Time:  time.Now(),
		Event: event,
		Level: lvl,
		l:     this,
	}
	return r
}

func (this *Logger) Enabled(lvl Level) bool {
	return lvl <= this.level
}
