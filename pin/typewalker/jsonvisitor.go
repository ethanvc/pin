package typewalker

import (
	"github.com/ethanvc/pin/pin/base"
	"reflect"
)

type JsonVisitor struct {
	b base.JsonBuilder
}

func (j *JsonVisitor) OpenStruct() {
	j.b.OpenObject()
}

func (j *JsonVisitor) CloseStruct() {
	j.b.CloseObject()
}

func (j *JsonVisitor) OpenArray() {
	j.b.OpenArray()
}

func (j *JsonVisitor) CloseArray() {
	j.b.CloseArray()
}

func (j *JsonVisitor) VisitNil() {
	j.b.WriteValueNull()
}

func (j *JsonVisitor) VisitField(w *TypeWalker, field reflect.StructField, v reflect.Value) {
	j.b.WriteKey(j.getKey(field))
	w.Visit(v)
}

func (j *JsonVisitor) getKey(field reflect.StructField) string {
	return field.Name
}
