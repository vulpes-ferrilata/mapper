package mapper

import (
	"reflect"
)

type key struct {
	From interface{}
	To   interface{}
}

type MappingFunc[From any, To any] func(from From) (To, error)

var mapper map[key]MappingFunc[any, any]

func init() {
	mapper = make(map[key]MappingFunc[any, any])
}

func wrapMappingFunc[From any, To any](f MappingFunc[From, To]) MappingFunc[any, any] {
	return func(from any) (any, error) {
		return f(from.(From))
	}
}

func parseMappingFunc[From any, To any](f MappingFunc[any, any]) MappingFunc[From, To] {
	return func(from From) (To, error) {
		to, err := f(from)
		return to.(To), err
	}
}

func CreateMap[From any, To any](f MappingFunc[*From, *To]) error {
	var emptyFrom From
	var emptyTo To

	if reflect.ValueOf(emptyFrom).Kind() != reflect.Struct {
		return ErrGenericParameterFromMustBeAStruct
	}

	if reflect.ValueOf(emptyTo).Kind() != reflect.Struct {
		return ErrGenericParameterToMustBeAStruct
	}

	mapper[key{emptyFrom, emptyTo}] = wrapMappingFunc(f)

	return nil
}

func Map[From any, To any](from *From) (*To, error) {
	var emptyFrom From
	var emptyTo To

	if reflect.ValueOf(emptyFrom).Kind() != reflect.Struct {
		return nil, ErrGenericParameterFromMustBeAStruct
	}

	if reflect.ValueOf(emptyTo).Kind() != reflect.Struct {
		return nil, ErrGenericParameterToMustBeAStruct
	}

	f, ok := mapper[key{emptyFrom, emptyTo}]
	if !ok {
		return nil, ErrMappingFunctionWasNotRegistered
	}

	return parseMappingFunc[*From, *To](f)(from)
}
