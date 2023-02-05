package typewalker

import (
	"reflect"
	"sync"
)

type processorCache sync.Map
type ProcessorFunc func(walker *TypeWalker, v reflect.Value)

type cacheKey struct {
	VisitorType reflect.Type
	ValType     reflect.Type
}

var sCache processorCache

func (cache *processorCache) Find(visitorType reflect.Type, valType reflect.Type) ProcessorFunc {
	f, ok := (*sync.Map)(cache).Load(cacheKey{
		VisitorType: visitorType,
		ValType:     valType,
	})
	if !ok {
		return nil
	}
	return f.(ProcessorFunc)
}

func (cache *processorCache) Store(visitorType reflect.Type, valType reflect.Type, f ProcessorFunc) {
	(*sync.Map)(cache).Store(cacheKey{
		VisitorType: visitorType,
		ValType:     valType,
	}, f)
}

type TypeWalker struct {
	visitor     TypeVisitor
	visitorType reflect.Type
	depth       int
}

func NewTypeWalker(visitor TypeVisitor) *TypeWalker {
	w := &TypeWalker{
		visitor:     visitor,
		visitorType: reflect.TypeOf(visitor),
	}
	return w
}

func (w *TypeWalker) Visit(v any) {
	if v == nil {
		w.visitor.VisitNil()
		return
	}
	w.getProcessor(reflect.TypeOf(v))(w, reflect.ValueOf(v))
}

func (w *TypeWalker) getProcessor(valType reflect.Type) ProcessorFunc {
	if f := sCache.Find(w.visitorType, valType); f != nil {
		return f
	}
	f := w.getProcessorSlow(valType)
	sCache.Store(w.visitorType, valType, f)
	return f
}

func (w *TypeWalker) getProcessorSlow(valType reflect.Type) ProcessorFunc {
	if implementCustomVisitor(valType) {

	}

	switch valType.Kind() {
	case reflect.Array:
		return w.newArrayProcessor(valType)
	case reflect.Chan:
		return dummyProcessor
	case reflect.Func:
		return dummyProcessor
	case reflect.Interface:
	case reflect.Map:
		return mapProcessor
	case reflect.Pointer:
	case reflect.Slice:
	case reflect.Struct:
		return newStructProcessor(valType)
	}

	return dummyProcessor
}

func dummyProcessor(walker *TypeWalker, v reflect.Value) {
}

func mapProcessor(walker *TypeWalker, v reflect.Value) {

}

type structProcessor struct {
	fields []Field
}

func newStructProcessor(valType reflect.Type) ProcessorFunc {
	return nil
}

func (s structProcessor) process(walker *TypeWalker, v reflect.Value) {
}

type arrayProcessor struct {
	elemProcessor ProcessorFunc
}

func (a arrayProcessor) process(w *TypeWalker, v reflect.Value) {

}

func (w *TypeWalker) newArrayProcessor(valType reflect.Type) ProcessorFunc {
	elemType := valType.Elem()
	f := w.getProcessor(elemType)
	return arrayProcessor{elemProcessor: f}.process
}
