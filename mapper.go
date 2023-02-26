package mapper

import (
	"fmt"
	"reflect"
)

type Key struct {
	From string
	To   string
}

var mapper map[Key]mappingFunc[any, any]

type mappingFunc[From any, To any] func(from From) (To, error)

func init() {
	mapper = make(map[Key]mappingFunc[any, any])
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

func CreateMap[From any, To any](f mappingFunc[From, To]) {
	var emptyFrom From
	var emptyTo To

	fromName := reflect.TypeOf(emptyFrom).String()
	toName := reflect.TypeOf(emptyTo).String()

	mapper[Key{fromName, toName}] = wrapMappingFunc(f)
}

func Map[From any, To any](from From) (To, error) {
	var emptyFrom From
	var emptyTo To

	fromName := reflect.TypeOf(emptyFrom).String()
	toName := reflect.TypeOf(emptyTo).String()

	f, ok := mapper[Key{fromName, toName}]
	if !ok {
		return emptyTo, fmt.Errorf("%w: from %s to %s", ErrMappingFunctionWasNotRegistered, fromName, toName)
	}

	return parseMappingFunc[From, To](f)(from)
}
