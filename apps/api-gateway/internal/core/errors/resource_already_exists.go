package errors

import "reflect"

type ResourceAlreadyExists struct {
	message string
}

func NewResourceAlreadyExists(message string) *ResourceAlreadyExists {
	return &ResourceAlreadyExists{
		message: message,
	}
}

func (e *ResourceAlreadyExists) Error() string {
	return e.message
}

func IsResourceAlreadyExists(err error) bool {
	return reflect.TypeOf(err).String() == "*errors.ResourceAlreadyExists"
}
