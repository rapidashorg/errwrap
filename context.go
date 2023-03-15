package errwrap

import (
	"context"
	"sync"
)

// ErrorData contains additional data for debugging purpose
type ErrorData map[string]interface{}
type errorDataWrapper struct {
	errorData ErrorData
	lock      sync.RWMutex
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

	curr := getErrorDataWrapper(ctx)
	if curr == nil {
		curr = &errorDataWrapper{
			errorData: make(ErrorData),
		}
	}

	curr.lock.Lock()
	for k, v := range data {
		curr.errorData[k] = v
	}
	curr.lock.Unlock()

	ctx = context.WithValue(ctx, contextKeyErrorData, curr)
	return ctx
}

// getErrorDataWrapper returns ErrorData from given context
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

func getErrorData(ctx context.Context) ErrorData {
	errDataWrapper := getErrorDataWrapper(ctx)
	if errDataWrapper == nil {
		return nil
	}

	errDataWrapper.lock.RLock()
	defer errDataWrapper.lock.RUnlock()

	errData := make(ErrorData)

	for k, v := range errDataWrapper.errorData {
		errData[k] = v
	}
	return errData
}
