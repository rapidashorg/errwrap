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

	if data == nil {
		return ctx
	}

	curr := getErrorData(ctx)

	for k, v := range curr {
		if _, exist := data[k]; exist {
			continue
		}
		data[k] = v
	}

	ctx = context.WithValue(ctx, contextKeyErrorData, data)
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

	errData, ok := errDataItf.(ErrorData)
	if !ok {
		return nil
	}

	return errData
}
