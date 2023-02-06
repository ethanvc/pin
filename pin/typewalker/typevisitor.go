package typewalker

import "reflect"

type Field struct {
	StructField  reflect.StructField
	Processor    ProcessorFunc
	CustomConfig any
}

type TypeVisitor interface {
	OpenStruct(w *TypeWalker)
	CloseStruct(w *TypeWalker)
	OpenArray(w *TypeWalker)
	CloseArray(w *TypeWalker)
	OpenMap(w *TypeWalker)
	CloseMap(w *TypeWalker)
	VisitNil(w *TypeWalker)
	VisitBool(w *TypeWalker, v reflect.Value)
	VisitInt64(w *TypeWalker, v reflect.Value)
	VisitUint64(w *TypeWalker, v reflect.Value)
	VisitFloat64(w *TypeWalker, v reflect.Value)
	VisitString(w *TypeWalker, v reflect.Value)
	VisitBytes(w *TypeWalker, v reflect.Value)
	VisitField(w *TypeWalker, field Field, v reflect.Value)
}

type CustomVisitor interface {
	Visit(w *TypeWalker) bool
}

var customVisitorType = reflect.TypeOf((*CustomVisitor)(nil)).Elem()

func implementCustomVisitor(t reflect.Type) bool {
	return t.Implements(customVisitorType)
}
