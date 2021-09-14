# ErrWrap

*Yet another error wrapper* for your Golang *web* application. This error wrapper is focused on:

- error codes to make error reporting easier,
- separation between user and developer error messages,
- easy to compare any errors, including the dynamic errors,
- add arbitrary debugging data to the error for debugging, and
- stack trace.

## Quickstart

```go
package main

import "fmt"
import "github.com/rapidashorg/errwrap"

// define the errors first
var ErrBadRequest     = errors.NewError(100, "Error bad request: %v", nil)
var ErrInternalServer = errors.NewMaskedError(101, "Error internal server error: %v", nil, nil)

func main() {
    // initialize the defined error
    err1 := ErrBadRequest.New("body not defined")
    // get the error message
    fmt.Println(err1.Error()) // Error bad request: body not defined (100)

    err2 := ErrInternalServer.New("file not exists")
    fmt.Println(err2.Error())       // Sorry, there are internal server error occured, please try again later. (101)
    fmt.Println(err2.ActualError()) // Error internal server error: file not exists (101)

    fmt.Println(err1.Is(ErrBadRequest)) // true
    fmt.Println(err2.Is(ErrBadRequest)) // false

    data := "an arbitrary data"
    err3 := ErrInternalServer.New("a forbidden data").WithData(map[string]string{
        "data": data,
    })
    fmt.Println(err3.ActualError())  // Error internal server error: a forbidden data (101)
    fmt.Println(err3.AdditionalData) // [{ "line": "....", "data": { "data": "an arbitrary data" } }]
    fmt.Println(err3.StackTrace)     // ["....", "....", "...."]
    fmt.Println(err3.HTTPErrorCode)  // 500
}
```

## Documentation

There will be 2 main structs in this package, which is `errors.ErrorDefinition` and `errors.ErrorWrapper` structs.

**`errors.ErrorDefinition` struct**

This struct is used to define an error. A single error will have a error code, formatted error message, and the HTTP code. There are 2 main functions to define an error, which is `errors.NewError()` and `errors.NewMaskedError()`. The difference between those two is what `errors.ErrorWrapper.Error()` function returns.

- `func NewError(code int, message string, httpCode *int) *ErrorDefinition`
    - This will create a new error definition.
    - If `httpCode` is `nil`, then default HTTP code will be used, which is `http.StatusBadRequest`.
- `func NewMaskedError(code int, message string, httpCode *int, maskMessage *string) *ErrorDefinition`
    - This will create a new masked error definition. The difference with `NewError` is when `errors.ErrorWrapper.Error()` called, which will return the `maskMessage` instead of the formatted `message`. This can be used to mask the error message that will given to users so they won’t be confused by our internal errors.
    - If `maskMessage` is `nil`, then default mask message will be used, which is `Sorry, there are internal server error occured, please try again later. (%d)`, with the `%d` as the error code.
    - If `httpCode` is `nil`, then default HTTP code will be used, which is `http.StatusInternalServerError`.

In `errors.ErrorDefinition` struct, there will be a single function `errors.ErrorDefinition.New()`, which is to create `errors.ErrorWrapper` object.

- `func (ed *ErrorDefinition) New(args ...interface{}) *ErrorWrapper`
    - This will create `errors.ErrorWrapper` object based on the error definition.
    - `args` is arguments that will be passed to `fmt.Sprintf` function to build formatted message. The `message` parameter when defining the error will be used as the format string.

**`errors.ErrorWrapper` struct**

This struct is used to wrap an error. There will be several variables in the struct, which is:

- `errors.ErrorWrapper.Code`: the error code
- `errors.ErrorWrapper.Message`: the formatted error message, but still raw (haven’t passed to `fmt.Sprintf`).
- `errors.ErrorWrapper.IsMasked`: is the error is masked? Will be `true` if generated using errors.NewMaskedError() function.
- `errors.ErrorWrapper.MaskMessage`: the masked error message, and will be used as the returned value from `errors.ErrorWrapper.Error()` function.
- `errors.ErrorWrapper.HTTPErrorCode`: the HTTP error code.
- `errors.ErrorWrapper.Args`: the arguments that will be passed to `fmt.Sprintf` function when building the error message.
- `errors.ErrorWrapper.StackTrace`: the stack trace when `errors.ErrorDefinition.New()` is called.
- `errors.ErrorWrapper.AdditionalData`: the additional data for the error. The entries will be appended when `errors.ErrorWrapper.WithData()` function is called.

There are 4 main functions on the `errors.ErrorWrapper` struct, which is `errors.ErrorWrapper.Error()`, `errors.ErrorWrapper.Is()`, `errors.ErrorWrapper.ActualError()`, and `errors.ErrorWrapper.WithData()`.

- `func (e *ErrorWrapper) Error() string`
    - This will return error message, either the formatted message that has been passed to `fmt.Sprintf`, or the mask message is the `IsMasked` variable is `true`. This function makes `*errors.ErrorWrapper` struct implements `error` interface.
- `func (e *ErrorWrapper) Is(ed *ErrorDefinition) bool`
    - This will compare the error wrapper with an error definition, and will returned `true` if the error wrapper is created using the provided error definition.
- `func (e *ErrorWrapper) ActualError() string`
    - This will return the actual error message that has been passed to `fmt.Sprintf`, completely ignores `IsMasked` variable.
- `func (e *ErrorWrapper) WithData(data map[string]interface{}) *ErrorWrapper`
    - This will append an entry to `AdditionalData` variable, and with line where the function is called.
