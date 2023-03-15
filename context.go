package errwrap

import "context"

// ErrorData contains additional data for debugging purpose
type ErrorData map[string]interface{}
type errorDataWrapper struct {
	data   ErrorData
	parent *errorDataWrapper
}
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

	parent := getErrorDataWrapper(ctx)

	curr := &errorDataWrapper{
		data:   data,
		parent: parent,
	}

	ctx = context.WithValue(ctx, contextKeyErrorData, curr)
	return ctx
}

func getErrorDataWrapper(ctx context.Context) *errorDataWrapper {
	if ctx == nil {
		return nil
	}

	errDataWrapperItf := ctx.Value(contextKeyErrorData)
	if errDataWrapperItf == nil {
		return nil
	}

	errDataWrapper, ok := errDataWrapperItf.(*errorDataWrapper)
	if !ok || errDataWrapper == nil {
		return nil
	}

	return errDataWrapper
}

// getErrorData returns ErrorData from given context
func getErrorData(ctx context.Context) ErrorData {
	if ctx == nil {
		return nil
	}

	var errData ErrorData

	errDataWrapper := getErrorDataWrapper(ctx)

	for {
		if errDataWrapper == nil {
			return errData
		}

		if errData == nil {
			errData = make(ErrorData)
		}

		for k, v := range errDataWrapper.data {
			if _, exists := errData[k]; exists {
				continue
			}
			errData[k] = v
		}

		errDataWrapper = errDataWrapper.parent
	}
}
