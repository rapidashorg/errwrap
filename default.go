package errwrap

import "fmt"

type stackTraceMode int

const (
	// StackTraceModeFull will gather full data of the stack traces (filename,
	// line number, and function name)
	StackTraceModeFull stackTraceMode = iota

	// StackTraceModeLineOnly will gather filename and line number only for data
	// of the stack traces
	StackTraceModeLineOnly

	// StackTraceModeFuncOnly will gather package and function name only for
	// data of the stack traces
	StackTraceModeFuncOnly
)

var (
	// DefaultMaskMessage defines the default mask message used when mask
	// message is not defined
	DefaultMaskMessage string = "Sorry, there are internal server error occured, please try again later."

	// DefaultMaskFormatter defines the mask formatter function used to format
	// mask message when the function is not defined
	DefaultMaskFormatter = func(erw ErrorWrapper) string {
		return erw.RawMaskMessage()
	}

	// DefaultMessageFormatter defines the formatter function used to format
	// formatted plain or mask message when the function is not defined.
	// In default, this function will set error code to plain or mask message.
	DefaultMessageFormatter = func(msg string, e ErrorWrapper) string {
		return fmt.Sprintf("%s (%d)", msg, e.Code())
	}

	// DefaultPackagePrefix defines the project package prefix. This variable is
	// used to trim stack trace to only include related project files. Set this
	// to empty string if you want to disable the trim functionality
	DefaultPackagePrefix string

	// DefaultStackTraceMode defines the mode used to gather stack traces data.
	DefaultStackTraceMode = StackTraceModeFull
)
