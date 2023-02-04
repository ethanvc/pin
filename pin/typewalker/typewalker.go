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
	Tag         string
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

func (w *TypeWalker) Visit(v any, visitor TypeVisitor) {
	if v == nil {
		visitor.VisitNil()
		return
	}

	w.visitor = visitor
	w.visitorType = reflect.TypeOf(visitor)
	w.depth = 0
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

	f := w.visitor.GetProcessor(valType, "")
	if f != nil {
		return f
	}
	return dummyProcessor
}

func dummyProcessor(walker *TypeWalker, v reflect.Value) {
}
