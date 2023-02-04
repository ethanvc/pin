package typewalker

import "reflect"

type TypeVisitor interface {
	VisitNil() bool
	GetProcessor(valType reflect.Type, tag string) ProcessorFunc
}

type CustomVisitor interface {
	Visit(w *TypeWalker) bool
}

var customVisitorType = reflect.TypeOf((*CustomVisitor)(nil)).Elem()

func implementCustomVisitor(t reflect.Type) bool {
	return t.Implements(customVisitorType)
}
