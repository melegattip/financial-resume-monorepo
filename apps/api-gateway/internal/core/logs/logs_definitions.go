package logs

import (
	"fmt"
	"sort"
)

type Log string

type Params map[string]interface{}

const (
	//--------------------INFO-----------------------
	InfoExternalAPICalled  Log = "An external API was called"
	InfoSubUseCaseExecuted Log = "Sub use case executed"
	InfoUseCaseExecuted    Log = "Use case executed"

	//--------------------WARNING--------------------
	WarningResourceLocked   Log = "Resource is already locked and processed"
	WarningDoesNotHaveScore Log = "Does not have a created resume"
	WarningPanicRecovered   Log = "Panic recovered"
	//--------------------ERROR----------------------
	ErrorExecuting             Log = "Error executing"
	ErrorUnmarshallingResponse Log = "Error unmarshalling response"
	ErrorBinding               Log = "Error binding object"
	ErrorGettingDataFromAPI    Log = "Error getting data from API"
	ErrorCreatingEndpoint      Log = "Error creating endpoint"
	ErrorUnauthorizedRequest   Log = "Unauthorized request"
	ErrorDeletingTransaction   Log = "Error deleting transaction"
	BindingError               Log = "Error binding request"
	ErrorInvalidParameter      Log = "Invalid parameter"
	ErrorInvalidRestMethod     Log = "Error invalid rest method"
	ErrorRecordNotFound        Log = "Error record not found"
	ErrorParsingDate           Log = "Error parsing date"
	ErrorParsingUserID         Log = "Error parsing user id"
)

func (l Log) GetMessage() string {
	return string(l)
}

func (l Log) GetMessageWithMapParams(params Params) string {
	msg := concatStrWithMap(l.GetMessage(), params)
	return msg
}

func concatStrWithMap(str string, params map[string]interface{}) string {
	keys := make([]string, 0)
	for k := range params {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	for _, k := range keys {
		m := fmt.Sprintf(" %v:%v", k, params[k])
		str += m
	}

	return str
}
