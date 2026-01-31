package errors

import "reflect"

type ResourceNotFound struct {
	message string
}

func NewResourceNotFound(message string) *ResourceNotFound {
	return &ResourceNotFound{
		message: message,
	}
}

func (e *ResourceNotFound) Error() string {
	return e.message
}

func IsResourceNotFound(err error) bool {
	return reflect.TypeOf(err).String() == "*errors.ResourceNotFound"
}
