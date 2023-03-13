package plog

import "context"

type contextKeyLogger struct{}

func WithLogger(c context.Context, l *Logger) context.Context {
	return context.WithValue(c, contextKeyLogger{}, l)
}

func FromCtx(c context.Context) *Logger {
	l, _ := c.Value(contextKeyLogger{}).(*Logger)
	if l != nil {
		return l
	}
	return Default()
}
