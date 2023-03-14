package plog

import (
	"fmt"
	"github.com/ethanvc/pin/pin/base"
	"time"
)

func ConsoleHandler(l *Logger, r Record) {
	bc := BasicLoggerContextFromCtx(l.C)
	sourceInfo := GetSourceFileInfo(r.Pc)
	var builder base.JsonBuilder
	builder.OpenObject()
	builder.WriteKey("t")
	builder.WriteValueString(r.Time.Format(time.RFC3339Nano))
	builder.WriteKey("lvl")
	builder.WriteValueString(r.Level.String())
	builder.WriteKey("tid")
	builder.WriteValueString(bc.TraceId)
	builder.WriteKey("event")
	builder.WriteValueString(r.Event)
	builder.WriteKey("f")
	builder.WriteValueString(sourceInfo.Info)
	builder.CloseObject().Finish()
	fmt.Printf("%s\n", builder.Buf.String())
}
