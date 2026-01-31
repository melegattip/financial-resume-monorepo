package errors

type ResourceParsingError struct {
	message string
}

func NewResourceParsingError(message string) ResourceParsingError {
	return ResourceParsingError{
		message: message,
	}
}

func (e ResourceParsingError) Error() string {
	return e.message
}
