package errwrap

import (
	"context"
	"fmt"
	"runtime"
)

// ErrorWrapper contains functions to define a wrapped error
type ErrorWrapper interface {
	error

	// Code is error code
	Code() int

	// CodeString is string-formatted error code
	CodeString() string

	// Category is error category
	Category() ErrorCategory

	// Masked determines if the error is masked
	Masked() bool

	// RawMessage is error message but haven't passed to fmt.Sprintf
	RawMessage() string

	// RawMaskMessage is error mask message but haven't paassed to mask
	// formatter
	RawMaskMessage() string

	// Args is error arguments to build error message
	Args() []interface{}

	// StackTrace is stack trace where the error is created
	StackTrace() []string

	// Data is additinal error data for further debugging
	Data() ErrorData

	// Is checks if errorWrapper is equals to ErrorDefinition
	Is(ed *ErrorDefinition) bool

	// ActualError returns error message but bypassing mask message
	ActualError() string
}

// Cast asserts error interface type to ErrorWrapper interface. If the error
// doesn't implement ErrorWrapper interface, returns nil.
func Cast(err error) ErrorWrapper {
	if erw, ok := err.(ErrorWrapper); ok {
		return erw
	}
	return nil
}

// Convert converts an ErrorWrapper into new ErrorWrapper based on
// *ErrorDefinition
func Convert(ctx context.Context, err ErrorWrapper, ed *ErrorDefinition) ErrorWrapper {
	ctx = InjectErrorData(ctx, err.Data())
	newErw := newErrorWrapper(ctx, ed, err.RawMessage(), err.Args()...)
	newErw.fillStackTrace(1)
	return newErw
}

// errorWrapper defines an error with the error code. the messages are in format
// string
type errorWrapper struct {
	code       int    // error code
	codeString string // error code in string

	message   string           // error message
	category  ErrorCategory    // error category
	formatter MessageFormatter // message formatter function

	isMasked      bool          // is error message masked?
	maskMessage   string        // error message mask
	maskFormatter MaskFormatter // mask formatter function

	args       []interface{}
	stackTrace []string
	data       ErrorData
}

// newErrorWrapper creates errorWrapper based on error definition
func newErrorWrapper(ctx context.Context, ed *ErrorDefinition, rawMessage string, args ...interface{}) *errorWrapper {
	erw := &errorWrapper{
		code:       ed.code,
		codeString: ed.codeString,

		message:   rawMessage,
		category:  ed.category,
		formatter: ed.formatter,

		isMasked:      ed.isMasked,
		maskMessage:   ed.maskMessage,
		maskFormatter: ed.maskFormatter,

		args: args,
		data: getErrorData(ctx),
	}
	return erw
}

func (e *errorWrapper) Code() int {
	return e.code
}

func (e *errorWrapper) CodeString() string {
	return e.codeString
}

func (e *errorWrapper) Category() ErrorCategory {
	return e.category
}

func (e *errorWrapper) Masked() bool {
	return e.isMasked
}

func (e *errorWrapper) RawMessage() string {
	return e.message
}

func (e *errorWrapper) RawMaskMessage() string {
	return e.maskMessage
}

func (e *errorWrapper) Args() []interface{} {
	return e.args
}

func (e *errorWrapper) StackTrace() []string {
	return e.stackTrace
}

func (e *errorWrapper) Data() ErrorData {
	return e.data
}

func (e *errorWrapper) Error() string {
	if e.isMasked {
		fn := DefaultMaskFormatter
		if e.maskFormatter != nil {
			fn = e.maskFormatter
		}

		return e.formatErrorMessage(fn(e))
	}
	return e.ActualError()
}

func (e *errorWrapper) Is(ed *ErrorDefinition) bool {
	if e == nil || ed == nil {
		return false
	}
	return e.code == ed.code
}

func (e *errorWrapper) ActualError() string {
	return e.formatErrorMessage(fmt.Sprintf(e.message, e.args...))
}

// formatErrorMessage formats message using formatter function
func (e *errorWrapper) formatErrorMessage(msg string) string {
	fn := DefaultMessageFormatter
	if e.formatter != nil {
		fn = e.formatter
	}

	return fn(msg, e)
}

// fillStackTrace fills errorWrapper stack trace
func (e *errorWrapper) fillStackTrace(offset int) {
	lines := make([]string, 0)

	for i := 1 + offset; ; i++ {
		fnptr, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}

		lines = append(lines, fmt.Sprintf("%s:%d (%s)", file, line, runtime.FuncForPC(fnptr).Name()))
	}

	e.stackTrace = lines
}

// getCallLine get current caller line
func (e *errorWrapper) getCallLine(offset int) string {
	// https://lawlessguy.wordpress.com/2016/04/17/display-file-function-and-line-number-in-go-golang/
	if fnptr, file, line, ok := runtime.Caller(1 + offset); ok {
		return fmt.Sprintf("%s:%d (%s)", file, line, runtime.FuncForPC(fnptr).Name())
	}

	return ""
}
