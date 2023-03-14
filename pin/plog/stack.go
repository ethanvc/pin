package plog

import "runtime"

func GetPc(skip int) uintptr {
	var pcs [1]uintptr
	runtime.Callers(skip+2, pcs[:])
	return pcs[0]
}

type SourceFileInfo struct {
	Info string
}

func GetSourceFileInfo(pc uintptr) SourceFileInfo {
	var info SourceFileInfo
	fs := runtime.CallersFrames([]uintptr{pc})
	f, _ := fs.Next()
	if f.Line == 0 {
		return info
	}
	return info
}
