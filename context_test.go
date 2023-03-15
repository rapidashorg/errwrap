package errwrap

import (
	"context"
	"fmt"
	"reflect"
	"sync"
	"testing"
	"time"
)

func TestConccurrentInjectErrorData(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ctx = InjectErrorData(ctx, ErrorData{"test1": "test1"})

	ctx1 := ctx
	ctx2 := ctx

	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-ctx1.Done():
				return
			default:
				ctx1 = InjectErrorData(ctx1, nil)
				ctx1 = InjectErrorData(ctx1, ErrorData{"test1": "goroutine1"})
			}
		}
	}()
	go func() {
		defer wg.Done()
		for {
			select {
			case <-ctx2.Done():
				return
			default:
				ctx2 = InjectErrorData(ctx2, nil)
				ctx2 = InjectErrorData(ctx2, ErrorData{"test1": "goroutine2"})
			}
		}
	}()

	time.Sleep(5 * time.Second)
	cancel()
	wg.Wait()

	ctx = InjectErrorData(ctx, ErrorData{"test1": "ctx"})
	ctx1 = InjectErrorData(ctx1, ErrorData{"test1": "ctx1"})
	ctx2 = InjectErrorData(ctx2, ErrorData{"test1": "ctx2"})

	if !reflect.DeepEqual(getErrorData(ctx)["test1"], "ctx2") {
		t.Error("not equal", getErrorData(ctx)["test1"])
	}
	if !reflect.DeepEqual(getErrorData(ctx1)["test1"], "ctx2") {
		t.Error("not equal", getErrorData(ctx1)["test1"])
	}
	if !reflect.DeepEqual(getErrorData(ctx2)["test1"], "ctx2") {
		t.Error("not equal", getErrorData(ctx2)["test1"])
	}
}

func BenchmarkInjectErrorData(b *testing.B) {
	oldGetErrData := func(ctx context.Context) ErrorData {
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

	oldMethod := func(ctx context.Context, data ErrorData) context.Context {
		if ctx == nil {
			return nil
		}

		curr := oldGetErrData(ctx)
		if curr == nil {
			curr = make(ErrorData)
		}

		for k, v := range data {
			curr[k] = v
		}

		ctx = context.WithValue(ctx, contextKeyErrorData, curr)
		return ctx
	}

	type dataStruct struct {
		data string
	}

	b.Run("new method", func(b *testing.B) {
		ctx := context.Background()
		for i := 0; i < b.N; i++ {
			ctx = InjectErrorData(ctx, nil)
			ctx = InjectErrorData(ctx, ErrorData{fmt.Sprintf("data_struct_%v", i): dataStruct{data: fmt.Sprintf("value_%v", i)}})
			ctx = InjectErrorData(ctx, ErrorData{fmt.Sprintf("data_string_%v", i): fmt.Sprintf("value_%v", i)})
			ctx = InjectErrorData(ctx, ErrorData{fmt.Sprintf("data_int_%v", i): i})
		}
		b.StopTimer()
		errData := getErrorData(ctx)
		if len(errData) != b.N*3 {
			b.Error("result is not equal")
		}
	})
	b.Run("old method", func(b *testing.B) {
		ctx := context.Background()
		for i := 0; i < b.N; i++ {
			ctx = oldMethod(ctx, nil)
			ctx = oldMethod(ctx, ErrorData{fmt.Sprintf("data_struct_%v", i): dataStruct{data: fmt.Sprintf("value_%v", i)}})
			ctx = oldMethod(ctx, ErrorData{fmt.Sprintf("data_string_%v", i): fmt.Sprintf("value_%v", i)})
			ctx = oldMethod(ctx, ErrorData{fmt.Sprintf("data_int_%v", i): i})
		}
		b.StopTimer()
		errData := oldGetErrData(ctx)
		if len(errData) != b.N*3 {
			b.Error("result is not equal")
		}
	})

	//Results
	//BenchmarkInjectErrorData/new_method-8            1000000              1598 ns/op             830 B/op         25 allocs/op
	//BenchmarkInjectErrorData/old_method-8             699990              2224 ns/op            1968 B/op         31 allocs/op
}

func TestInjectErrorData(t *testing.T) {
	type args struct {
		ctx  context.Context
		data ErrorData
	}
	tests := []struct {
		name string
		args args
		want context.Context
	}{
		{
			name: "success context is nil",
			args: args{
				ctx: nil,
				data: ErrorData{
					"foo": "bar",
				},
			},
			want: nil,
		},
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				data: ErrorData{
					"foo": "bar",
				},
			},
			want: context.WithValue(context.Background(), contextKeyErrorData, &errorDataWrapper{
				errorData: map[string]interface{}{
					"foo": "bar",
				},
				lock: sync.RWMutex{},
			}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := InjectErrorData(tt.args.ctx, tt.args.data); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InjectErrorData() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getErrorData(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name string
		args args
		want ErrorData
	}{
		{
			name: "success context is nil",
			args: args{
				ctx: nil,
			},
			want: nil,
		},
		{
			name: "success appended",
			args: args{
				ctx: InjectErrorData(
					InjectErrorData(context.Background(), ErrorData{
						"foo": "bar",
					}), ErrorData{
						"bar": "baz",
					},
				),
			},
			want: ErrorData{
				"foo": "bar",
				"bar": "baz",
			},
		},
		{
			name: "success overwritten",
			args: args{
				ctx: InjectErrorData(
					InjectErrorData(context.Background(), ErrorData{
						"foo": "bar",
					}), ErrorData{
						"foo": "baz",
					},
				),
			},
			want: ErrorData{
				"foo": "baz",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getErrorData(tt.args.ctx); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getErrorData() = %v, want %v", got, tt.want)
			}
		})
	}
}
