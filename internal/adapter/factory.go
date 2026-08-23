package adapter

import (
	"errors"
	"reflect"
)

var ErrNilAdapter = errors.New("adapter dependency is nil")

type Adapter interface{ Normalize([]byte) ([]byte, error) }

type Factory struct{ adapter Adapter }

func NewFactory(candidate Adapter) (*Factory, error) {
	if isNilInterface(candidate) {
		return nil, ErrNilAdapter
	}
	return &Factory{adapter: candidate}, nil
}

func (f *Factory) Build() Adapter { return f.adapter }

func isNilInterface(value any) bool {
	if value == nil {
		return true
	}
	rv := reflect.ValueOf(value)
	switch rv.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Pointer, reflect.Slice:
		return rv.IsNil()
	default:
		return false
	}
}
