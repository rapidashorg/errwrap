package errwrap

import (
	"context"
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

	if !reflect.DeepEqual(getErrorData(ctx)["test1"], "test1") {
		t.Error("not equal", getErrorData(ctx)["test1"])
	}
	if !reflect.DeepEqual(getErrorData(ctx1)["test1"], "goroutine1") {
		t.Error("not equal", getErrorData(ctx1)["test1"])
	}
	if !reflect.DeepEqual(getErrorData(ctx2)["test1"], "goroutine2") {
		t.Error("not equal", getErrorData(ctx2)["test1"])
	}
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
			want: context.WithValue(context.Background(), contextKeyErrorData, ErrorData{
				"foo": "bar",
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
