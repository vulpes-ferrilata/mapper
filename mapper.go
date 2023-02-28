package mapper

import (
	"reflect"
)

type key struct {
	From interface{}
	To   interface{}
}

var mapper map[key]mappingFunc[any, any]

type mappingFunc[From any, To any] func(from From) (To, error)

func init() {
	mapper = make(map[key]mappingFunc[any, any])
}

func wrapMappingFunc[From any, To any](f mappingFunc[From, To]) mappingFunc[any, any] {
	return func(from any) (any, error) {
		return f(from.(From))
	}
}

func parseMappingFunc[From any, To any](f mappingFunc[any, any]) mappingFunc[From, To] {
	return func(from From) (To, error) {
		to, err := f(from)
		return to.(To), err
	}
}

func CreateMap[From comparable, To comparable](f mappingFunc[*From, *To]) error {
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

func Map[From comparable, To comparable](from *From) (*To, error) {
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
