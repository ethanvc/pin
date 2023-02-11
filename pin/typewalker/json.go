package typewalker

import (
	"github.com/ethanvc/pin/pin/base"
	"sync"
)

var walkerPool sync.Pool

func getWalker() (*TypeWalker, *JsonVisitor) {
	w, _ := walkerPool.Get().(*TypeWalker)
	if w != nil {
		v := w.Visitor().(*JsonVisitor)
		v.B.Reset()
		return w, v
	}

	var visitor JsonVisitor
	return NewTypeWalker(&visitor), &visitor

}

func ToLogJson(v any) []byte {
	w, visitor := getWalker()
	w.Visit(v)
	visitor.B.Finish()
	buf := append([]byte(nil), visitor.B.Bytes()...)
	walkerPool.Put(w)
	return buf
}

func ToLogJsonStr(v any) string {
	return base.BytesToStr(ToLogJson(v))
}
