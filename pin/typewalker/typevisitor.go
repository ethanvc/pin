package typewalker

import "reflect"

type TypeVisitor interface {
	OpenStruct()
	CloseStruct()
	VisitNil()
	VisitField(w *TypeWalker, field reflect.StructField, v reflect.Value)
}

type CustomVisitor interface {
	Visit(w *TypeWalker) bool
}

var customVisitorType = reflect.TypeOf((*CustomVisitor)(nil)).Elem()

func implementCustomVisitor(t reflect.Type) bool {
	return t.Implements(customVisitorType)
}
