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
}

func newField(field reflect.StructField) *Field {
	f := &Field{
		StructField: field,
	}
	tagVal := field.Tag.Get("json")
	n, param, _ := strings.Cut(tagVal, ",")
	if len(n) > 0 {
		f.JsonKey = n
	} else {
		f.JsonKey = field.Name
	}
	f.Quoted = strings.Contains(param, "string")
	f.OmitEmpty = strings.Contains(param, "omitempty")
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

type TypeVisitor interface {
	SetWalker(w *TypeWalker)
	OpenStruct()
	CloseStruct()
	OpenArray()
	CloseArray()
	VisitInt64(field *Field, v reflect.Value, key bool)
	VisitUint64(field *Field, v reflect.Value, key bool)
	VisitFloat(field *Field, v reflect.Value, key bool)
	VisitString(field *Field, v reflect.Value, key bool)
	VisitBytes(field *Field, v reflect.Value, key bool)
	VisitField(field *Field, v reflect.Value)
	OpenMap()
	CloseMap()
	GetProcessor(valType reflect.Type) ProcessorFunc
}
