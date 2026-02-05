// Package logs provides types and constants for structured logging throughout the application
package logs

// Tags represents a map of key-value pairs for additional log context
type Tags map[string]interface{}

// LogMessage represents a structured log message with a predefined text
type LogMessage struct {
	Message string
}

// GetMessage returns the message text of the log message
func (l LogMessage) GetMessage() string {
	return l.Message
}

var (
	// ErrorLoadingConfiguration represents an error when loading configuration fails
	ErrorLoadingConfiguration = LogMessage{
		Message: "Error loading configuration",
	}
	// ErrorCreatingTransaction represents an error when creating a transaction fails
	ErrorCreatingTransaction = LogMessage{
		Message: "Error creating transaction",
	}
	// ErrorListingTransactions represents an error when listing transactions fails
	ErrorListingTransactions = LogMessage{
		Message: "Error listing transactions",
	}
	// ErrorGettingTransaction represents an error when retrieving a transaction fails
	ErrorGettingTransaction = LogMessage{
		Message: "Error getting transaction",
	}
	// ErrorUpdatingTransaction represents an error when updating a transaction fails
	ErrorUpdatingTransaction = LogMessage{
		Message: "Error updating transaction",
	}
	// ErrorCreatingCategory represents an error when creating a category fails
	ErrorCreatingCategory = LogMessage{
		Message: "Error creating category",
	}
	// ErrorListingCategories represents an error when listing categories fails
	ErrorListingCategories = LogMessage{
		Message: "Error listing categories",
	}
	// ErrorUpdatingCategory represents an error when updating a category fails
	ErrorUpdatingCategory = LogMessage{
		Message: "Error updating category",
	}
	// ErrorDeletingCategory represents an error when deleting a category fails
	ErrorDeletingCategory = LogMessage{
		Message: "Error deleting category",
	}
	// ErrorGeneratingReport represents an error when generating a financial report fails
	ErrorGeneratingReport = LogMessage{
		Message: "Error generating financial report",
	}
)
