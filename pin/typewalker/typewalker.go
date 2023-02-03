package typewalker

import (
	"reflect"
	"sync"
)

type structVisitorCache sync.Map
type encodeFunc func(walker *TypeWalker, v reflect.Value)

type cacheKey struct {
	visitorType reflect.Type
	valType     reflect.Type
}

var sCache structVisitorCache

func (cache *structVisitorCache) Find(visitorType reflect.Type, valType reflect.Type) encodeFunc {
	f, ok := (*sync.Map)(cache).Load(cacheKey{
		visitorType: visitorType,
		valType:     valType,
	})
	if !ok {
		return nil
	}
	return f.(encodeFunc)
}

func (cache *structVisitorCache) Store(visitorType reflect.Type, valType reflect.Type, f encodeFunc) {
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
	w.getEncoder(reflect.TypeOf(v))(w, reflect.ValueOf(v))
}

func (w *TypeWalker) getEncoder(valType reflect.Type) encodeFunc {
	if f := sCache.Find(w.visitorType, valType); f != nil {
		return f
	}
	return nil
}
