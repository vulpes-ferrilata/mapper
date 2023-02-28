package mapper

import "errors"

var (
	ErrGenericParameterFromMustBeAStruct = errors.New("generic parameter From must be a struct")
	ErrGenericParameterToMustBeAStruct   = errors.New("generic parameter To must be a struct")
	ErrMappingFunctionWasNotRegistered   = errors.New("mapping function was not registered")
)
