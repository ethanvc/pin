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
		return structProcessor
	}

	return dummyProcessor
}

func dummyProcessor(walker *TypeWalker, v reflect.Value) {
}

func mapProcessor(walker *TypeWalker, v reflect.Value) {

}

func structProcessor(walker *TypeWalker, v reflect.Value) {
	valType := v.Type()
	walker.visitor.OpenStruct()
	for i := 0; i < valType.NumField(); i++ {
		field := valType.Field(i)
		fieldVal := v.Field(i)
		walker.visitor.VisitField(walker, field, fieldVal)
	}
	walker.visitor.CloseStruct()
}
