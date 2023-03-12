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
