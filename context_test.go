package errwrap

import (
	"context"
	"reflect"
	"testing"
)

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
