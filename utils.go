package mapper

import "reflect"

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

func getFullPath(t reflect.Type) string {
	return t.PkgPath() + "." + t.Name()
}
