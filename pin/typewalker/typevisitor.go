package typewalker

import "reflect"

type Field struct {
	StructField reflect.StructField
	Processor   ProcessorFunc
	JsonKey     string
}

type MapKey struct {
}

type TypeVisitor interface {
	SetWalker(w *TypeWalker)
	OpenStruct()
	CloseStruct()
	OpenArray()
	CloseArray()
	VisitNil()
	VisitBool(field *Field, v reflect.Value)
	VisitInt64(field *Field, v reflect.Value)
	VisitUint64(field *Field, v reflect.Value)
	VisitFloat64(field *Field, v reflect.Value)
	VisitString(field *Field, v reflect.Value)
	VisitBytes(field *Field, v reflect.Value)
	VisitField(field *Field, v reflect.Value)
}

type CustomVisitor interface {
	Visit(w *TypeWalker) bool
}

var customVisitorType = reflect.TypeOf((*CustomVisitor)(nil)).Elem()

func implementCustomVisitor(t reflect.Type) bool {
	return t.Implements(customVisitorType)
}
