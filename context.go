package errwrap

import "context"

// ErrorData contains additional data for debugging purpose
type ErrorData map[string]interface{}
type contextKey string

var (
	contextKeyErrorData contextKey = "ctxk-errordata"
)

// InjectErrorData injects the data into context with error context key. If the
// context has been injected and injected again with same key but different
// value, the old value will be overwritten with the new value.
func InjectErrorData(ctx context.Context, data ErrorData) context.Context {
	if ctx == nil {
		return nil
	}

	curr := getErrorData(ctx)
	if curr == nil {
		curr = make(ErrorData)
	}

	for k, v := range data {
		curr[k] = v
	}

	ctx = context.WithValue(ctx, contextKeyErrorData, curr)
	return ctx
}

// getErrorData returns ErrorData from given context
func getErrorData(ctx context.Context) ErrorData {
	if ctx == nil {
		return nil
	}

	errDataItf := ctx.Value(contextKeyErrorData)
	if errDataItf == nil {
		return nil
	}

	errData, _ := errDataItf.(ErrorData)
	return errData
}
