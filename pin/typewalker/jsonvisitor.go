package typewalker

import (
	"encoding/base64"
	"github.com/ethanvc/pin/pin/base"
	"reflect"
)

type JsonVisitor struct {
	w *TypeWalker
	B base.JsonBuilder
}

func (j *JsonVisitor) SetWalker(w *TypeWalker) {
	j.w = w
}

func (j *JsonVisitor) OpenStruct() {
	j.B.OpenObject()
}

func (j *JsonVisitor) CloseStruct() {
	j.B.CloseObject()
}

func (j *JsonVisitor) OpenArray() {
	j.B.OpenArray()
}

func (j *JsonVisitor) CloseArray() {
	j.B.CloseArray()
}

func (j *JsonVisitor) VisitInt64(field *Field, v reflect.Value, key bool) {
	j.B.WriteValueInt64(v.Int())
}

func (j *JsonVisitor) VisitUint64(field *Field, v reflect.Value, key bool) {
	j.B.WriteValueUint64(v.Uint())
}

func (j *JsonVisitor) VisitString(field *Field, v reflect.Value, key bool) {
	if key {
		j.B.WriteKey(v.String())
	} else {
		j.B.WriteValueString(v.String())
	}
}

func (j *JsonVisitor) VisitBytes(field *Field, v reflect.Value, key bool) {
	j.B.WriteValueString(base64.StdEncoding.EncodeToString(v.Bytes()))
}

func (j *JsonVisitor) VisitField(field *Field, v reflect.Value) {
	j.B.WriteKey(field.JsonKey)
	field.Processor(j.w, field, v.FieldByIndex(field.StructField.Index))
}

func (j *JsonVisitor) OpenMap() {
	j.OpenStruct()
}

func (j *JsonVisitor) CloseMap() {
	j.CloseStruct()

}
