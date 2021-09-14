package errors

import (
	"fmt"
	"runtime"
)

// ErrorWrapper defines an error with the error code. the messages are in format
// string
type ErrorWrapper struct {
	Code          int    // error code
	Message       string // error message
	IsMasked      bool   // is error message masked?
	MaskMessage   string // error message mask
	HTTPErrorCode int    // http error code

	Args           []interface{}
	StackTrace     []string
	AdditionalData []*ErrorWrapperAdditionalData
}

type ErrorWrapperAdditionalData struct {
	Line string
	Data map[string]interface{}
}

// NewErrorWrapper creates ErrorWrapper based on error definition
func newErrorWrapper(ed *ErrorDefinition, args ...interface{}) *ErrorWrapper {
	return &ErrorWrapper{
		Code:          ed.code,
		Message:       ed.message,
		IsMasked:      ed.isMasked,
		MaskMessage:   ed.maskMessage,
		HTTPErrorCode: ed.httpErrorCode,
		Args:          args,
	}
}

func (e *ErrorWrapper) Error() string {
	if e.IsMasked {
		if e.MaskMessage == DefaultMaskMessage {
			return fmt.Sprintf(e.MaskMessage, e.Code)
		}
		return e.MaskMessage
	}
	return e.ActualError()
}

// Is checks if ErrorWrapper is equals to ErrorDefinition
func (e *ErrorWrapper) Is(ed *ErrorDefinition) bool {
	return e.Code == ed.code
}

// ActualError returns error message but bypassing mask message
func (e *ErrorWrapper) ActualError() string {
	return fmt.Sprintf("%s (%d)", fmt.Sprintf(e.Message, e.Args...), e.Code)
}

// WithData append current caller line and additinal data to ErrorWrapper
func (e *ErrorWrapper) WithData(data map[string]interface{}) *ErrorWrapper {
	if e.AdditionalData == nil {
		e.AdditionalData = make([]*ErrorWrapperAdditionalData, 0)
	}

	e.AdditionalData = append(e.AdditionalData, &ErrorWrapperAdditionalData{
		Line: e.getCallLine(1),
		Data: data,
	})

	return e
}

// fillStackTrace fills ErrorWrapper stack trace
func (e *ErrorWrapper) fillStackTrace(offset int) {
	lines := make([]string, 0)

	for i := 1 + offset; ; i++ {
		fnptr, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}

		lines = append(lines, fmt.Sprintf("%s:%d (%s)", file, line, runtime.FuncForPC(fnptr).Name()))
	}

	e.StackTrace = lines
}

// getCallLine get current caller line
func (e *ErrorWrapper) getCallLine(offset int) string {
	// https://lawlessguy.wordpress.com/2016/04/17/display-file-function-and-line-number-in-go-golang/
	if fnptr, file, line, ok := runtime.Caller(1 + offset); ok {
		return fmt.Sprintf("%s:%d (%s)", file, line, runtime.FuncForPC(fnptr).Name())
	}

	return ""
}
