package errors

// ErrorDefinition defines an error definition
type ErrorDefinition struct {
	code          int    // error code
	message       string // error message
	isMasked      bool   // is error message masked?
	maskMessage   string // error message mask
	httpErrorCode int    // http error code
}

// NewError creates simple error definition
func NewError(code int, message string, httpCode *int) *ErrorDefinition {
	if httpCode == nil {
		httpCode = &DefaultHTTPErrorCode
	}

	return &ErrorDefinition{
		code:          code,
		message:       message,
		httpErrorCode: *httpCode,
	}
}

// NewMaskedError creates masked error definition. the difference between
// newError() is how error string handled. masked error definition will use
// masked error message as standard error message that will be given to user,
// but use actual error message as log error message
func NewMaskedError(code int, message string, httpCode *int, maskMessage *string) *ErrorDefinition {
	if httpCode == nil {
		httpCode = &DefaultMaskedHTTPErrorCode
	}

	if maskMessage == nil {
		maskMessage = &DefaultMaskMessage
	}

	errDef := NewError(code, message, httpCode)
	errDef.isMasked = true
	errDef.maskMessage = *maskMessage

	return errDef
}

// New creates new ErrorWrapper based on error definition
func (ed *ErrorDefinition) New(args ...interface{}) *ErrorWrapper {
	ew := newErrorWrapper(ed, args...)
	ew.fillStackTrace(1)

	return ew
}
