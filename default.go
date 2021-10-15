package errwrap

import "fmt"

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
)
