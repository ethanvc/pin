package plog

import (
	"github.com/ethanvc/pin/pin/attrrepo"
	"sync/atomic"
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
	return nil
}
