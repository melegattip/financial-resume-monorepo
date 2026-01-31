package errors

import "reflect"

type TooManyRequests struct {
	message string
}

func NewTooManyRequests(message string) *TooManyRequests {
	return &TooManyRequests{
		message: message,
	}
}

func (e *TooManyRequests) Error() string {
	return e.message
}

func IsTooManyRequestError(err error) bool {
	return reflect.TypeOf(err).String() == "*errors.TooManyRequests"
}
