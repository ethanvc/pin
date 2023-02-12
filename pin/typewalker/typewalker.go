package typewalker

import (
	"reflect"
	"sync"
)

type ProcessorFunc func(walker *TypeWalker, field *Field, v reflect.Value, key bool)

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
	w.getProcessorByValue(vv)(w, nil, vv, true)
}

func (w *TypeWalker) getProcessorByValue(v reflect.Value) ProcessorFunc {
	if v.Kind() == reflect.Invalid {
		return nilProcess
	}
	valType := v.Type()
	return w.getProcessor(valType)
}

func (w *TypeWalker) getProcessor(valType reflect.Type) ProcessorFunc {
	if f := sCache.Find(w.visitorType, valType); f != nil {
		return f
	}

	var wg sync.WaitGroup
	var f ProcessorFunc

	wg.Add(1)
	f, ok := sCache.LoadOrStore(w.visitorType, valType, func(w *TypeWalker, field *Field, v reflect.Value, key bool) {
		wg.Wait()
		f(w, field, v, key)
	})
	if ok {
		return f
	}

	f = w.getProcessorSlow(valType)
	wg.Done()
	sCache.Store(w.visitorType, valType, f)
	return f
}

func nilProcess(w *TypeWalker, field *Field, v reflect.Value, key bool) {
}

func (w *TypeWalker) getProcessorSlow(valType reflect.Type) ProcessorFunc {
	if f := w.Visitor().GetProcessor(valType); f != nil {
		return f
	}

	switch valType.Kind() {
	case reflect.Chan:
		return dummyProcessor
	case reflect.Func:
		return dummyProcessor
	case reflect.Interface:
	case reflect.Map:
		return w.newMapProcessor(valType)
	case reflect.Pointer:
	case reflect.Slice:
		return w.newSliceProcessor(valType, false)
	case reflect.Array:
		return w.newSliceProcessor(valType, true)
	case reflect.String:
		return stringProcessor{}.process
	case reflect.Struct:
		return w.newStructProcessor(valType)
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
		return intProcessor{}.process
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
		return uintProcessor{}.process
	case reflect.Float32, reflect.Float64:
		return floatProcess
	}

	return dummyProcessor
}

func floatProcess(w *TypeWalker, field *Field, v reflect.Value, key bool) {
	w.Visitor().VisitFloat(field, v, key)
}

func (w *TypeWalker) newMapProcessor(valType reflect.Type) ProcessorFunc {
	p := mapProcessor{}
	p.keyProcessor = w.getProcessor(valType.Key())
	p.valProcessor = w.getProcessor(valType.Elem())
	return p.process
}

type mapProcessor struct {
	keyProcessor ProcessorFunc
	valProcessor ProcessorFunc
}

func (p mapProcessor) process(w *TypeWalker, field *Field, v reflect.Value, key bool) {
	w.Visitor().OpenMap()
	iter := v.MapRange()
	for iter.Next() {
		key := iter.Key()
		p.keyProcessor(w, field, key, true)
		val := iter.Value()
		p.valProcessor(w, field, val, false)
	}
	w.Visitor().CloseMap()
}

type stringProcessor struct {
}

func (p stringProcessor) process(w *TypeWalker, field *Field, v reflect.Value, key bool) {
	w.Visitor().VisitString(field, v, key)
}

type intProcessor struct {
}

func (p intProcessor) process(w *TypeWalker, field *Field, v reflect.Value, key bool) {
	w.Visitor().VisitInt64(field, v, key)
}

type uintProcessor struct {
}

func (p uintProcessor) process(w *TypeWalker, field *Field, v reflect.Value, key bool) {
	w.Visitor().VisitUint64(field, v, key)
}

func dummyProcessor(walker *TypeWalker, field *Field, v reflect.Value, key bool) {
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
		newF.Processor = w.getProcessor(field.Type)
		p.fields = append(p.fields, newF)
		rootField = field
	}
	return p.process
}

func (s structProcessor) process(w *TypeWalker, _ *Field, v reflect.Value, key bool) {
	w.visitor.OpenStruct()
	for _, field := range s.fields {
		w.Visitor().VisitField(field, v)
	}
	w.visitor.CloseStruct()
}

type sliceProcessor struct {
	elemProcessor ProcessorFunc
}

func (a sliceProcessor) process(w *TypeWalker, field *Field, v reflect.Value, key bool) {
	w.Visitor().OpenArray()
	for i := 0; i < v.Len(); i++ {
		a.elemProcessor(w, field, v.Index(i), key)
	}
	w.Visitor().CloseArray()
}

func (w *TypeWalker) newSliceProcessor(valType reflect.Type, array bool) ProcessorFunc {
	elemType := valType.Elem()
	if !array && elemType.Kind() == reflect.Uint8 {
		return bytesProcessor
	}
	f := w.getProcessor(elemType)
	return sliceProcessor{elemProcessor: f}.process
}

func bytesProcessor(w *TypeWalker, field *Field, v reflect.Value, key bool) {
	w.Visitor().VisitBytes(field, v, key)
}
