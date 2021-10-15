package errwrap

import "context"

// MessageFormatter is formatter used to format the message
type MessageFormatter func(msg string, erw ErrorWrapper) string

// MaskFormatter is formatter used to format the mask message
type MaskFormatter func(erw ErrorWrapper) string

// ErrorCategory defines the category of the error. This is the replacement of
// HTTPErrorCode, so applications can fine tune what should they respond when
// they get the error, for example when encounters error with BadRequest
// category, the application can return response with the 400 HTTP code.
type ErrorCategory int

// ErrorDefinition defines an error definition
type ErrorDefinition struct {
	code          int              // error code
	codeString    string           // error code in string
	message       string           // error message
	isMasked      bool             // is error message masked?
	formatter     MessageFormatter // message formatter function
	maskMessage   string           // error message mask
	maskFormatter MaskFormatter    // mask formatter function
	category      ErrorCategory    // error category
}

// NewError creates simple error definition
func NewError(code int, codeString string, message string, category ErrorCategory) *ErrorDefinition {
	return &ErrorDefinition{
		code:       code,
		codeString: codeString,
		message:    message,
		category:   category,
		formatter:  DefaultMessageFormatter,
	}
}

// Masked masks this error definition, makes produced errorWrapper message
// masked with maskMessage. The mask message used by this function is the
// default one.
func (ed *ErrorDefinition) Masked() *ErrorDefinition {
	ed.isMasked = true
	ed.maskMessage = DefaultMaskMessage
	ed.maskFormatter = DefaultMaskFormatter
	return ed
}

// MaskedMessage masks this error definition, makes produced errorWrapper
// message masked with maskMessage. The mask message used by this function is
// passed as arguments, and the mask formatter function used is the default one.
func (ed *ErrorDefinition) MaskedMessage(maskMessage string) *ErrorDefinition {
	ed.isMasked = true
	ed.maskMessage = maskMessage
	ed.maskFormatter = DefaultMaskFormatter
	return ed
}

// MaskedFunction masks this error definition, makes produced ErrorWrapper
// message masked with maskMessage. No mask message is used by this function,
// the mask message is provided from running passed function.
func (ed *ErrorDefinition) MaskedFunction(fn MaskFormatter) *ErrorDefinition {
	ed.isMasked = true
	ed.maskFormatter = fn
	return ed
}

// MessageFormatter sets the message formatter
func (ed *ErrorDefinition) MessageFormatter(fn MessageFormatter) *ErrorDefinition {
	ed.formatter = fn
	return ed
}

// NewWithoutContext creates new ErrorWrapper based on error definition without
// passed context
func (ed *ErrorDefinition) NewWithoutContext(args ...interface{}) ErrorWrapper {
	erw := newErrorWrapper(context.Background(), ed, args...)
	erw.fillStackTrace(1)
	return erw
}

// New creates new ErrorWrapper based on error definition
func (ed *ErrorDefinition) New(ctx context.Context, args ...interface{}) ErrorWrapper {
	erw := newErrorWrapper(ctx, ed, args...)
	erw.fillStackTrace(1)
	return erw
}
