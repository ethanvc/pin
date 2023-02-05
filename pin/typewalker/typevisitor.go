package typewalker

import "reflect"

type Field struct {
	StructField *reflect.StructField
	VisitCache  any
}

type TypeVisitor interface {
	OpenStruct()
	CloseStruct()
	VisitNil()
	VisitField(w *TypeWalker, field Field, v reflect.Value)
}

type CustomVisitor interface {
	Visit(w *TypeWalker) bool
}

var customVisitorType = reflect.TypeOf((*CustomVisitor)(nil)).Elem()

func implementCustomVisitor(t reflect.Type) bool {
	return t.Implements(customVisitorType)
}
