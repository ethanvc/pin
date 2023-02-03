package typewalker

type TypeVisitor interface {
	VisitNil()
}

type CustomVisitor interface {
	Visit(w *TypeWalker)
}
