# ErrWrap

> **DISCLAIMER**: This package is at its early development, and is not stable yet before reaching v1. Always expect breaking changes on every changes without further announcements, as we still adjust the API to where this can comfortably fit to most of usecases. If you know the risks to use this package in this current state, you may proceed to use this.

*Yet another error wrapper* for your Golang *web* application. This error wrapper is focused on:

- error codes to make error reporting easier,
- separation between user and developer error messages through message masking,
- easy to compare any errors, including the dynamic errors,
- add arbitrary debugging data to the error for debugging, and
- stack trace.

## Quickstart

```go
package main

import (
    "context"
    "fmt"

    "github.com/rapidashorg/errwrap"
)

// define error categories, this is helpful to categorize the error, and might
// used when creating responses
const (
    ErrCategoryBadRequest errwrap.ErrorCategory = iota
    ErrCategoryInternalServerError
)

// define the errors first
var ErrBadRequest     = errwrap.NewError(100, "ErrBadRequest", "Error bad request: %v", ErrCategoryBadRequest)
var ErrInternalServer = errwrap.NewError(101, "ErrInternalServer", "Error internal server error: %v", ErrCategoryInternalServerError).Masked()

func main() {
    // initialize the defined error
    err1 := ErrBadRequest.NewWithoutContext("body not defined")
    // get the error message
    fmt.Println(err1.Error()) // Error bad request: body not defined (100)

    err2 := ErrInternalServer.NewWithoutContext("file not exists")
    fmt.Println(err2.Error())       // Sorry, there are internal server error occured, please try again later. (101)
    fmt.Println(err2.ActualError()) // Error internal server error: file not exists (101)

    fmt.Println(err1.Is(ErrBadRequest)) // true
    fmt.Println(err2.Is(ErrBadRequest)) // false

    data := "an arbitrary data"

    // injects error data to context, this is helpful for debugging
    ctx := errwrap.InjectErrorData(context.Background(), errwrap.ErrorData{
        "data": data,
    })

    err3 := ErrInternalServer.New(ctx, "a forbidden data")
    fmt.Println(err3.ActualError())  // Error internal server error: a forbidden data (101)
    fmt.Println(err3.Data()) // { "data": "an arbitrary data" } }
    fmt.Println(err3.CodeString()) // "ErrInternalServer"
    fmt.Println(err3.StackTrace())     // ["....", "....", "...."]
    fmt.Println(err3.Category())  // errwrap.ErrorCategory(1)

    if err3.Category() == ErrCategoryBadRequest {
        // do something if the error is categorized as bad request, e.g. set
        // http error code to 400 when building response
    }
}
```

## Documentation

There will be 2 main structs in this package, which is `errors.ErrorDefinition` struct and `errors.ErrorWrapper` interface.

**`errors.ErrorDefinition` struct**

This struct is used to define an error. A single error will have a error code, error code in string, formatted error message, and error category. There is a single functions to define an error, which is `errors.NewError()`:

- `func NewError(code int, codeString string, category ErrorCategory) *ErrorDefinition`
    - This will create a new error definition.
    - Difference between `code` and `codeString` is how it's used. In our case, `code` is used to construct user error message as the numerical error code is anonymized form of error, and `codeString` is used by the developer for metrics tags, to give meaningful error message in metrics dashboard instead of using numeric error code.

In `errors.ErrorDefinition` struct, there will be several functions:

- `func (ed *ErrorDefinition) Masked() *ErrorDefinition`
    - This will set the error definition to use mask state, which will create masked error wrappers.
    - Use `errwrap.DefaultMaskMessage` and `errwrap.DefaultMaskFormatter` as the default message and formatter value.
- `func (ed *ErrorDefinition) MaskedMessage(maskMessage string) *ErrorDefinition`
    - Same as `errors.ErrorDefinition.Masked()`, but we customize the mask message.
    - Use `errwrap.DefaultMaskFormatter` as the default formatter value.
- `func (ed *ErrorDefinition) MaskedFunction(fn MessageFunction) *ErrorDefinition`
    - Same as `errors.ErrorDefinition.Masked()`, but we customize the mask formatter function.
    - This sets the mask message to empty string, so we expect the mask message to be created within given mask formatter function.
- `func (ed *ErrorDefinition) MessageFormatter(fn MessageFunction) *ErrorDefinition`
    - Sets the message formatter function used to format the message
- `func (ed *ErrorDefinition) NewWithoutContext(rawMessage string, args ...interface{}) ErrorWrapper`
    - This will create `errors.ErrorWrapper` object based on the error definition.
    - `args` is arguments that will be passed to `fmt.Sprintf` function to build formatted message. The `rawMessage` parameter will be used as the format string.
- `func (ed *ErrorDefinition) New(ctx context.Context, rawMessage string, args ...interface{}) ErrorWrapper`
    - Same as `errors.ErrorDefinition.NewWithoutContext()`, but we can pass context to the error. This context is used to inject error data for debugging purpose.

**`errors.ErrorWrapper` interface**

This interface is used to wrap an error. There will be several functions defined by the interface, which is:

- `func (ErrorWrapper) Code() int`
    - The error code.
- `func (ErrorWrapper) CodeString() string`
    - The error code string.
- `func (ErrorWrapper) Category() ErrorCategory`
    - The error category.
- `func (ErrorWrapper) Masked() boolean`
    - Is the error is masked? Will be `true` if the error definition is chain-called with mask functions.
- `func (ErrorWrapper) RawMessage() string`
    - The raw error message (haven’t passed to `fmt.Sprintf()`).
- `func (ErrorWrapper) RawMaskMessage() string`
    - The raw error mask message (haven’t passed to `fmt.Sprintf()`).
    - If `true`, then the error mask message will be returned instead from `func (ErrorWrapper) Error()` function.
- `func (ErrorWrapper) Args()`
    - The arguments that will be passed to `fmt.Sprintf()` function when building the error message and/or optionally error mask message too.
- `func (ErrorWrapper) StackTrace()`
    - The stack trace when `errors.ErrorDefinition.New()` or `errors.ErrorDefinition.NewWithoutContext()` is called.
- `func (ErrorWrapper) Data()`
    - The related error data for the error, usually for debugging purpose.
    - The value will be filled from passed context, that has been injected by `errwrap.ErrorData` using `errwrap.InjectErrorData` function.
- `func (e *ErrorWrapper) Error() string`
    - This will return error message, either the formatted message that has been passed to `fmt.Sprintf()`, or the mask message is the `IsMasked` variable is `true`.
    - This function exists because `errwrap.ErrorWrapper` interface extends `error` interface.
- `func (e *ErrorWrapper) Is(ed *ErrorDefinition) bool`
    - This will compare the error wrapper with an error definition, and will returned `true` if the error wrapper is created using the provided error definition.
- `func (e *ErrorWrapper) ActualError() string`
    - This will return the actual error message that has been passed to `fmt.Sprintf()`, completely ignores whether the error is masked or not.
