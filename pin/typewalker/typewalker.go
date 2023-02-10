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
	return w.getProcessor(valType)
}

func (w *TypeWalker) getProcessor(valType reflect.Type) ProcessorFunc {
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

	f = w.getProcessorSlow(valType)
	wg.Done()
	sCache.Store(w.visitorType, valType, f)
	return f
}

func nilProcess(w *TypeWalker, field *Field, v reflect.Value) {
	w.Visitor().VisitNil()
}

func customVisitorProcess(w *TypeWalker, v reflect.Value) {

}

func (w *TypeWalker) getProcessorSlow(valType reflect.Type) ProcessorFunc {
	if implementCustomVisitor(valType) {

	}

	switch valType.Kind() {
	case reflect.Chan:
		return dummyProcessor
	case reflect.Func:
		return dummyProcessor
	case reflect.Interface:
	case reflect.Map:
	case reflect.Pointer:
	case reflect.Slice:
		return w.newSliceProcessor(valType)
	case reflect.Array:
	case reflect.String:
		return stringProcessor
	case reflect.Struct:
		return w.newStructProcessor(valType)
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
		return intProcessor
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
		return uintProcessor
	}

	return dummyProcessor
}

func stringProcessor(w *TypeWalker, field *Field, v reflect.Value) {
	w.Visitor().VisitString(field, v)
}

func uintProcessor(w *TypeWalker, field *Field, v reflect.Value) {
	w.Visitor().VisitUint64(field, v)
}

func intProcessor(w *TypeWalker, field *Field, v reflect.Value) {
	w.Visitor().VisitInt64(field, v)
}

func dummyProcessor(walker *TypeWalker, field *Field, v reflect.Value) {
}

func mapProcessor(walker *TypeWalker, v reflect.Value) {

}

type structProcessor struct {
	fields []*Field
}

func ancestorOf(ancestor, t *reflect.StructField) bool {
	if ancestor == nil {
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
	fields := reflect.VisibleFields(valType)
	var p structProcessor
	var ignoredField *reflect.StructField
	for _, field := range fields {
		if ancestorOf(ignoredField, &field) {
			continue
		}
		if field.Tag.Get("json") == "-" {
			ignoredField = &field
			continue
		}
		newF := newField(field)
		newF.Processor = w.getProcessor(field.Type)
		p.fields = append(p.fields, newF)

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

func (w *TypeWalker) newSliceProcessor(valType reflect.Type) ProcessorFunc {
	elemType := valType.Elem()
	if elemType.Kind() == reflect.Uint8 {
		return bytesProcessor
	}
	f := w.getProcessor(elemType)
	return sliceProcessor{elemProcessor: f}.process
}

func bytesProcessor(w *TypeWalker, field *Field, v reflect.Value) {
	w.Visitor().VisitBytes(field, v)
}
