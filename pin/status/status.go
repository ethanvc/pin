package status

import "github.com/ethanvc/pin/pin/status/code"

type Status struct {
	codeVal     code.Code
	eventVal    string
	subEventVal string
	msgVal      string
	pcVal       uintptr
}
