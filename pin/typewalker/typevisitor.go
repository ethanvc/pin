package typewalker

import (
	"reflect"
	"strings"
)

type Field struct {
	StructField reflect.StructField
	Processor   ProcessorFunc
	JsonKey     string
	OmitEmpty   bool
	Quoted      bool
	Ignore      bool
}

func newField(field reflect.StructField) *Field {
	f := &Field{
		StructField: field,
	}
	tagVal := field.Tag.Get("json")
	if tagVal == "-" {
		f.Ignore = true
		return f
	}
	if tagVal != "-" {
		n, param, _ := strings.Cut(tagVal, ",")
		if len(n) > 0 {
			f.JsonKey = n
		} else {
			f.JsonKey = field.Name
		}
		f.Quoted = strings.Contains(param, "string")
		f.OmitEmpty = strings.Contains(param, "omitempty")
	}
	return f
}

func (f *Field) NeedIgnore() bool {
	return len(f.JsonKey) == 0
}

func (f *Field) AncestorOf(otherF *Field) bool {
	if len(f.StructField.Index) >= len(otherF.StructField.Index) {
		return false
	}
	for i := 0; i < len(f.StructField.Index); i++ {
		if f.StructField.Index[i] != otherF.StructField.Index[i] {
			return false
		}
	}
	return true
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
