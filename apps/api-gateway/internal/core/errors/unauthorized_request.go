package errors

import "reflect"

type UnauthorizedRequest struct {
	message string
}

func NewUnauthorizedRequest(message string) *UnauthorizedRequest {
	return &UnauthorizedRequest{
		message: message,
	}
}

func (e *UnauthorizedRequest) Error() string {
	return e.message
}

func IsUnauthorizedRequestError(err error) bool {
	return reflect.TypeOf(err).String() == "*errors.UnauthorizedRequest"
}
