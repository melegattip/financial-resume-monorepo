package errors

import "reflect"

type ResourceLockedError struct {
	message string
}

func NewResourceLockedError(message string) ResourceLockedError {
	return ResourceLockedError{
		message: message,
	}
}

func (e ResourceLockedError) Error() string {
	return e.message
}

func IsResourceLocked(err error) bool {
	return reflect.TypeOf(err).String() == "errors.ResourceLockedError"
}
