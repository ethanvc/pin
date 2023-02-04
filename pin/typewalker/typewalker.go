package typewalker

import (
	"reflect"
	"sync"
)

type processorCache sync.Map
type processorFunc func(walker *TypeWalker, v reflect.Value)

type cacheKey struct {
	visitorType reflect.Type
	valType     reflect.Type
}

var sCache processorCache

func (cache *processorCache) Find(visitorType reflect.Type, valType reflect.Type) processorFunc {
	f, ok := (*sync.Map)(cache).Load(cacheKey{
		visitorType: visitorType,
		valType:     valType,
	})
	if !ok {
		return nil
	}
	return f.(processorFunc)
}

func (cache *processorCache) Store(visitorType reflect.Type, valType reflect.Type, f processorFunc) {
	(*sync.Map)(cache).Store(cacheKey{
		visitorType: visitorType,
		valType:     valType,
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

func (w *TypeWalker) getProcessor(valType reflect.Type) processorFunc {
	if f := sCache.Find(w.visitorType, valType); f != nil {
		return f
	}
	f := w.getProcessorSlow(valType)
	sCache.Store(w.visitorType, valType, f)
	return f
}

func (w *TypeWalker) getProcessorSlow(valType reflect.Type) processorFunc {
	if implementCustomVisitor(valType) {

	}
	return dummyProcessor
}

func dummyProcessor(walker *TypeWalker, v reflect.Value) {
}
