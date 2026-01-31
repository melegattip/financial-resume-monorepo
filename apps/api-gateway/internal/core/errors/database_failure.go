package errors

type DatabaseFailure struct {
	message string
}

func NewDatabaseFailure(message string) DatabaseFailure {
	return DatabaseFailure{
		message: message,
	}
}

func (e DatabaseFailure) Error() string {
	return e.message
}
