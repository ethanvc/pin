package typewalker

import (
	"reflect"
	"sync"
)

type ProcessorFunc func(walker *TypeWalker, field *Field, v reflect.Value)

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
	visitor.SetWalker(w)
	return w
}

func (w *TypeWalker) Visitor() TypeVisitor {
	return w.visitor
}

func (w *TypeWalker) Visit(v any) {
	vv := reflect.ValueOf(v)
	w.getProcessorByValue(vv)(w, nil, vv)
}

func (w *TypeWalker) getProcessorByValue(v reflect.Value) ProcessorFunc {
	if v.Kind() == reflect.Invalid {
		return nilProcess
	}
	valType := v.Type()
	return w.getProcessor(valType, false)
}

func (w *TypeWalker) getProcessor(valType reflect.Type, key bool) ProcessorFunc {
	if f := sCache.Find(w.visitorType, valType); f != nil {
		return f
	}

	var wg sync.WaitGroup
	var f ProcessorFunc

	wg.Add(1)
	f, ok := sCache.LoadOrStore(w.visitorType, valType, func(w *TypeWalker, field *Field, v reflect.Value) {
		wg.Wait()
		f(w, field, v)
	})
	if ok {
		return f
	}

	f = w.getProcessorSlow(valType, key)
	wg.Done()
	sCache.Store(w.visitorType, valType, f)
	return f
}

func nilProcess(w *TypeWalker, field *Field, v reflect.Value) {
}

func (w *TypeWalker) getProcessorSlow(valType reflect.Type, key bool) ProcessorFunc {
	if implementCustomVisitor(valType) {

	}

	switch valType.Kind() {
	case reflect.Chan:
		return dummyProcessor
	case reflect.Func:
		return dummyProcessor
	case reflect.Interface:
	case reflect.Map:
		return w.newMapProcessor(valType, key)
	case reflect.Pointer:
	case reflect.Slice:
		return w.newSliceProcessor(valType, false, key)
	case reflect.Array:
		return w.newSliceProcessor(valType, true, key)
	case reflect.String:
		return stringProcessor{key: key}.process
	case reflect.Struct:
		return w.newStructProcessor(valType)
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
		return intProcessor{key: key}.process
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
		return uintProcessor{key: key}.process
	}

	return dummyProcessor
}

func (w *TypeWalker) newMapProcessor(valType reflect.Type, key bool) ProcessorFunc {
	p := mapProcessor{}
	p.keyProcessor = w.getProcessor(valType.Key(), true)
	p.valProcessor = w.getProcessor(valType.Elem(), false)
	p.key = key
	return p.process
}

type mapProcessor struct {
	keyProcessor ProcessorFunc
	valProcessor ProcessorFunc
	key          bool
}

func (p mapProcessor) process(w *TypeWalker, field *Field, v reflect.Value) {
	w.Visitor().OpenMap()
	iter := v.MapRange()
	for iter.Next() {
		key := iter.Key()
		p.keyProcessor(w, field, key)
		val := iter.Value()
		p.valProcessor(w, field, val)
	}
	w.Visitor().CloseMap()
}

type stringProcessor struct {
	key bool
}

func (p stringProcessor) process(w *TypeWalker, field *Field, v reflect.Value) {
	w.Visitor().VisitString(field, v, p.key)
}

type intProcessor struct {
	key bool
}

func (p intProcessor) process(w *TypeWalker, field *Field, v reflect.Value) {
	w.Visitor().VisitInt64(field, v, p.key)
}

type uintProcessor struct {
	key bool
}

func (p uintProcessor) process(w *TypeWalker, field *Field, v reflect.Value) {
	w.Visitor().VisitUint64(field, v, p.key)
}

func dummyProcessor(walker *TypeWalker, field *Field, v reflect.Value) {
}

type structProcessor struct {
	fields []*Field
}

func ancestorOf(ancestor, t reflect.StructField) bool {
	if len(ancestor.Index) == 0 {
		return false
	}
	if len(ancestor.Index) >= len(t.Index) {
		return false
	}
	for i := 0; i < len(ancestor.Index); i++ {
		if ancestor.Index[i] != t.Index[i] {
			return false
		}
	}
	return true
}

func (w *TypeWalker) newStructProcessor(valType reflect.Type) ProcessorFunc {
	var p structProcessor
	var rootField reflect.StructField
	fields := reflect.VisibleFields(valType)
	for _, field := range fields {
		if ancestorOf(rootField, field) {
			continue
		}
		jsonPart := field.Tag.Get("json")
		if jsonPart == "-" {
			rootField = field
			continue
		}
		if field.Anonymous && jsonPart == "" {
			continue
		}
		newF := newField(field)
		newF.Processor = w.getProcessor(field.Type, false)
		p.fields = append(p.fields, newF)
		rootField = field
	}
	return p.process
}

func (s structProcessor) process(w *TypeWalker, _ *Field, v reflect.Value) {
	w.visitor.OpenStruct()
	for _, field := range s.fields {
		w.Visitor().VisitField(field, v)
	}
	w.visitor.CloseStruct()
}

type sliceProcessor struct {
	elemProcessor ProcessorFunc
}

func (a sliceProcessor) process(w *TypeWalker, field *Field, v reflect.Value) {
	w.Visitor().OpenArray()
	for i := 0; i < v.Len(); i++ {
		a.elemProcessor(w, field, v.Index(i))
	}
	w.Visitor().CloseArray()
}

func (w *TypeWalker) newSliceProcessor(valType reflect.Type, array bool, key bool) ProcessorFunc {
	elemType := valType.Elem()
	if !array && elemType.Kind() == reflect.Uint8 {
		return bytesProcessor
	}
	f := w.getProcessor(elemType, key)
	return sliceProcessor{elemProcessor: f}.process
}

func bytesProcessor(w *TypeWalker, field *Field, v reflect.Value) {
	w.Visitor().VisitBytes(field, v, false)
}
