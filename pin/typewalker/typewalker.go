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
	case reflect.Array:
		return w.newArrayProcessor(valType)
	case reflect.Chan:
		return dummyProcessor
	case reflect.Func:
		return dummyProcessor
	case reflect.Interface:
	case reflect.Map:
	case reflect.Pointer:
	case reflect.Slice:
	case reflect.Struct:
		return w.newStructProcessor(valType)
	}

	return dummyProcessor
}

func dummyProcessor(walker *TypeWalker, field *Field, v reflect.Value) {
}

func mapProcessor(walker *TypeWalker, v reflect.Value) {

}

type structProcessor struct {
	fields []*Field
}

func (w *TypeWalker) newStructProcessor(valType reflect.Type) ProcessorFunc {
	fields := reflect.VisibleFields(valType)
	var p structProcessor
	for _, field := range fields {
		newField := Field{
			StructField: field,
			Processor:   w.getProcessor(field.Type),
			JsonKey:     field.Name,
		}
		p.fields = append(p.fields, &newField)
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

type arrayProcessor struct {
	elemProcessor ProcessorFunc
}

func (a arrayProcessor) process(w *TypeWalker, field *Field, v reflect.Value) {

}

func (w *TypeWalker) newArrayProcessor(valType reflect.Type) ProcessorFunc {
	elemType := valType.Elem()
	f := w.getProcessor(elemType)
	return arrayProcessor{elemProcessor: f}.process
}
