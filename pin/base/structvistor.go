package base

import "sync"

type StructVisitorCache sync.Map

type StructWalker struct {
	visitor StructVisitor
	depth   int
	cache   *StructVisitorCache
}

func (w StructWalker) Visit(v interface{}, visitor StructVisitor) {
	w.visitor = visitor
	w.depth = 0
	w.cache = visitor.GetCache()
}

type StructVisitor interface {
	GetCache() *StructVisitorCache
}
