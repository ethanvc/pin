package typewalker

import "github.com/ethanvc/pin/pin/base"

func ToLogJson(v any) []byte {
	var visitor JsonVisitor
	w := NewTypeWalker(&visitor)
	w.Visit(v)
	visitor.B.Finish()
	return visitor.B.Bytes()
}

func ToLogJsonStr(v any) string {
	return base.BytesToStr(ToLogJson(v))
}
