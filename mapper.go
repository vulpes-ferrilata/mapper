package mapper

import (
	"reflect"
)

type key struct {
	From string
	To   string
}

type mappingFunc[From any, To any] func(from From) (To, error)

type Mapper struct {
	m map[key]mappingFunc[any, any]
}

func Register[From any, To any](m *Mapper) func(mappingFunc[*From, *To]) error {
	return func(f mappingFunc[*From, *To]) error {
		var emptyFrom From
		var emptyTo To

		if reflect.ValueOf(emptyFrom).Kind() != reflect.Struct {
			return ErrGenericParameterFromMustBeAStruct
		}

		if reflect.ValueOf(emptyTo).Kind() != reflect.Struct {
			return ErrGenericParameterToMustBeAStruct
		}

		k := key{
			From: getFullPath(reflect.TypeOf(emptyFrom)),
			To:   getFullPath(reflect.TypeOf(emptyTo)),
		}

		if m.m == nil {
			m.m = make(map[key]mappingFunc[any, any])
		}

		m.m[k] = wrapMappingFunc(f)

		return nil
	}
}

func Map[From any, To any](m *Mapper) mappingFunc[*From, *To] {
	return func(from *From) (*To, error) {
		var emptyFrom From
		var emptyTo To

		if reflect.ValueOf(emptyFrom).Kind() != reflect.Struct {
			return nil, ErrGenericParameterFromMustBeAStruct
		}

		if reflect.ValueOf(emptyTo).Kind() != reflect.Struct {
			return nil, ErrGenericParameterToMustBeAStruct
		}

		k := key{
			From: getFullPath(reflect.TypeOf(emptyFrom)),
			To:   getFullPath(reflect.TypeOf(emptyTo)),
		}

		f, ok := m.m[k]
		if !ok {
			return nil, ErrMappingFunctionWasNotRegistered
		}

		return parseMappingFunc[*From, *To](f)(from)
	}
}
