package typewalker

import (
	"reflect"
	"sync"
)

type processorCache sync.Map

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

func (cache *processorCache) LoadOrStore(visitorType reflect.Type, valType reflect.Type, f ProcessorFunc) (ProcessorFunc, bool) {
	fany, ok := (*sync.Map)(cache).LoadOrStore(cacheKey{
		VisitorType: visitorType,
		ValType:     valType,
	}, f)

	return fany.(ProcessorFunc), ok
}
