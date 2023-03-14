package plog

import (
	"context"
	"sync/atomic"
	"time"
)

type Logger struct {
	Handlers []Handler
	Level    Level
	C        context.Context
}

var defaultLogger atomic.Pointer[Logger]

func init() {
	l := &Logger{
		Handlers: []Handler{ConsoleHandler},
		Level:    LevelInfo,
		C: WithBasicLoggerContext(nil, BasicLoggerContext{
			TraceId: GenerateTraceId(),
		}),
	}
	defaultLogger.Store(l)
}

func Default() *Logger {
	return defaultLogger.Load()
}

func SetDefault(l *Logger) {
	defaultLogger.Store(l)
}

func (this *Logger) Clone() *Logger {
	return &*this
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
	return lvl <= this.Level
}
