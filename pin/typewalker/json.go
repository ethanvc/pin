package typewalker

import "github.com/ethanvc/pin/pin/base"

func ToJson(v any) []byte {
	var visitor JsonVisitor
	w := NewTypeWalker(&visitor)
	w.Visit(v)
	return nil
}

func ToJsonStr(v any) string {
	return base.BytesToStr(ToJson(v))
}
