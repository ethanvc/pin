package plog

import "runtime"

func GetPc(skip int) uintptr {
	var pcs [1]uintptr
	runtime.Callers(skip+2, pcs[:])
	return pcs[0]
}
