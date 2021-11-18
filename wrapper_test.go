package errwrap

import (
	"context"
	"errors"
	"reflect"
	"testing"
)

func TestCast(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name string
		args args
		want ErrorWrapper
	}{
		{
			name: "success",
			args: args{
				err: &errorWrapper{
					code:       100,
					codeString: "ErrTest",
					message:    "Test error message: %s",
					category:   ErrorCategory(1),
				},
			},
			want: &errorWrapper{
				code:       100,
				codeString: "ErrTest",
				message:    "Test error message: %s",
				category:   ErrorCategory(1),
			},
		},
		{
			name: "success error not implemented ErrorWrapper",
			args: args{
				err: errors.New("an error"),
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Cast(tt.args.err); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Cast() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConvert(t *testing.T) {
	type args struct {
		ctx context.Context
		err ErrorWrapper
		ed  *ErrorDefinition
	}
	tests := []struct {
		name string
		args args
		want ErrorWrapper
	}{
		{
			name: "success",
			args: args{
				ctx: InjectErrorData(context.Background(), ErrorData{
					"bar": "baz",
				}),
				err: &errorWrapper{
					code:       100,
					codeString: "ErrTest",
					message:    "Test error message",
					category:   ErrorCategory(1),
					data: ErrorData{
						"foo": "bar",
					},
				},
				ed: &ErrorDefinition{
					code:       101,
					codeString: "ErrTestNew",
					category:   ErrorCategory(2),
				},
			},
			want: &errorWrapper{
				code:       101,
				codeString: "ErrTestNew",
				message:    "Test error message",
				category:   ErrorCategory(2),
				data: ErrorData{
					"foo": "bar",
					"bar": "baz",
				},

				formatter:     nil,
				maskMessage:   DefaultMaskMessage,
				maskFormatter: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Convert(tt.args.ctx, tt.args.err, tt.args.ed)
			if g, ok := got.(*errorWrapper); ok {
				g.stackTrace = nil
				g.formatter = nil
				g.maskFormatter = nil
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Convert() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_newErrorWrapper(t *testing.T) {
	type args struct {
		ctx        context.Context
		ed         *ErrorDefinition
		rawMessage string
		args       []interface{}
	}
	tests := []struct {
		name string
		args args
		want *errorWrapper
	}{
		{
			name: "success",
			args: args{
				ctx: InjectErrorData(context.Background(), ErrorData{
					"foo": "bar",
				}),
				ed: &ErrorDefinition{
					code:       100,
					codeString: "ErrTest",
					category:   ErrorCategory(1),
				},
				rawMessage: "Test error message: %s",
				args:       []interface{}{"Foo"},
			},
			want: &errorWrapper{
				code:       100,
				codeString: "ErrTest",
				message:    "Test error message: %s",
				category:   ErrorCategory(1),
				args:       []interface{}{"Foo"},
				data: ErrorData{
					"foo": "bar",
				},

				formatter:     nil,
				maskMessage:   DefaultMaskMessage,
				maskFormatter: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := newErrorWrapper(tt.args.ctx, tt.args.ed, tt.args.rawMessage, tt.args.args...)

			got.formatter = nil
			got.maskFormatter = nil

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newErrorWrapper() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_errorWrapper_Code(t *testing.T) {
	type fields struct {
		code          int
		codeString    string
		message       string
		category      ErrorCategory
		formatter     MessageFormatter
		isMasked      bool
		maskMessage   string
		maskFormatter MaskFormatter
		args          []interface{}
		stackTrace    []string
		data          ErrorData
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			name: "success",
			fields: fields{
				code:       100,
				codeString: "ErrTest",
				message:    "Test error message: %s",
				category:   ErrorCategory(1),
				args:       []interface{}{"Foo"},
				data: ErrorData{
					"foo": "bar",
				},
			},
			want: 100,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &errorWrapper{
				code:          tt.fields.code,
				codeString:    tt.fields.codeString,
				message:       tt.fields.message,
				category:      tt.fields.category,
				formatter:     tt.fields.formatter,
				isMasked:      tt.fields.isMasked,
				maskMessage:   tt.fields.maskMessage,
				maskFormatter: tt.fields.maskFormatter,
				args:          tt.fields.args,
				stackTrace:    tt.fields.stackTrace,
				data:          tt.fields.data,
			}
			if got := e.Code(); got != tt.want {
				t.Errorf("errorWrapper.Code() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_errorWrapper_CodeString(t *testing.T) {
	type fields struct {
		code          int
		codeString    string
		message       string
		category      ErrorCategory
		formatter     MessageFormatter
		isMasked      bool
		maskMessage   string
		maskFormatter MaskFormatter
		args          []interface{}
		stackTrace    []string
		data          ErrorData
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "success",
			fields: fields{
				code:       100,
				codeString: "ErrTest",
				message:    "Test error message: %s",
				category:   ErrorCategory(1),
				args:       []interface{}{"Foo"},
				data: ErrorData{
					"foo": "bar",
				},
			},
			want: "ErrTest",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &errorWrapper{
				code:          tt.fields.code,
				codeString:    tt.fields.codeString,
				message:       tt.fields.message,
				category:      tt.fields.category,
				formatter:     tt.fields.formatter,
				isMasked:      tt.fields.isMasked,
				maskMessage:   tt.fields.maskMessage,
				maskFormatter: tt.fields.maskFormatter,
				args:          tt.fields.args,
				stackTrace:    tt.fields.stackTrace,
				data:          tt.fields.data,
			}
			if got := e.CodeString(); got != tt.want {
				t.Errorf("errorWrapper.CodeString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_errorWrapper_Category(t *testing.T) {
	type fields struct {
		code          int
		codeString    string
		message       string
		category      ErrorCategory
		formatter     MessageFormatter
		isMasked      bool
		maskMessage   string
		maskFormatter MaskFormatter
		args          []interface{}
		stackTrace    []string
		data          ErrorData
	}
	tests := []struct {
		name   string
		fields fields
		want   ErrorCategory
	}{
		{
			name: "success",
			fields: fields{
				code:       100,
				codeString: "ErrTest",
				message:    "Test error message: %s",
				category:   ErrorCategory(1),
				args:       []interface{}{"Foo"},
				data: ErrorData{
					"foo": "bar",
				},
			},
			want: ErrorCategory(1),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &errorWrapper{
				code:          tt.fields.code,
				codeString:    tt.fields.codeString,
				message:       tt.fields.message,
				category:      tt.fields.category,
				formatter:     tt.fields.formatter,
				isMasked:      tt.fields.isMasked,
				maskMessage:   tt.fields.maskMessage,
				maskFormatter: tt.fields.maskFormatter,
				args:          tt.fields.args,
				stackTrace:    tt.fields.stackTrace,
				data:          tt.fields.data,
			}
			if got := e.Category(); got != tt.want {
				t.Errorf("errorWrapper.Category() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_errorWrapper_Masked(t *testing.T) {
	type fields struct {
		code          int
		codeString    string
		message       string
		category      ErrorCategory
		formatter     MessageFormatter
		isMasked      bool
		maskMessage   string
		maskFormatter MaskFormatter
		args          []interface{}
		stackTrace    []string
		data          ErrorData
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "success",
			fields: fields{
				code:       100,
				codeString: "ErrTest",
				message:    "Test error message: %s",
				category:   ErrorCategory(1),
				args:       []interface{}{"Foo"},
				data: ErrorData{
					"foo": "bar",
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &errorWrapper{
				code:          tt.fields.code,
				codeString:    tt.fields.codeString,
				message:       tt.fields.message,
				category:      tt.fields.category,
				formatter:     tt.fields.formatter,
				isMasked:      tt.fields.isMasked,
				maskMessage:   tt.fields.maskMessage,
				maskFormatter: tt.fields.maskFormatter,
				args:          tt.fields.args,
				stackTrace:    tt.fields.stackTrace,
				data:          tt.fields.data,
			}
			if got := e.Masked(); got != tt.want {
				t.Errorf("errorWrapper.Masked() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_errorWrapper_RawMessage(t *testing.T) {
	type fields struct {
		code          int
		codeString    string
		message       string
		category      ErrorCategory
		formatter     MessageFormatter
		isMasked      bool
		maskMessage   string
		maskFormatter MaskFormatter
		args          []interface{}
		stackTrace    []string
		data          ErrorData
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "success",
			fields: fields{
				code:       100,
				codeString: "ErrTest",
				message:    "Test error message: %s",
				category:   ErrorCategory(1),
				args:       []interface{}{"Foo"},
				data: ErrorData{
					"foo": "bar",
				},
			},
			want: "Test error message: %s",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &errorWrapper{
				code:          tt.fields.code,
				codeString:    tt.fields.codeString,
				message:       tt.fields.message,
				category:      tt.fields.category,
				formatter:     tt.fields.formatter,
				isMasked:      tt.fields.isMasked,
				maskMessage:   tt.fields.maskMessage,
				maskFormatter: tt.fields.maskFormatter,
				args:          tt.fields.args,
				stackTrace:    tt.fields.stackTrace,
				data:          tt.fields.data,
			}
			if got := e.RawMessage(); got != tt.want {
				t.Errorf("errorWrapper.RawMessage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_errorWrapper_RawMaskMessage(t *testing.T) {
	type fields struct {
		code          int
		codeString    string
		message       string
		category      ErrorCategory
		formatter     MessageFormatter
		isMasked      bool
		maskMessage   string
		maskFormatter MaskFormatter
		args          []interface{}
		stackTrace    []string
		data          ErrorData
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "success",
			fields: fields{
				code:       100,
				codeString: "ErrTest",
				message:    "Test error message: %s",
				category:   ErrorCategory(1),
				args:       []interface{}{"Foo"},
				data: ErrorData{
					"foo": "bar",
				},
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &errorWrapper{
				code:          tt.fields.code,
				codeString:    tt.fields.codeString,
				message:       tt.fields.message,
				category:      tt.fields.category,
				formatter:     tt.fields.formatter,
				isMasked:      tt.fields.isMasked,
				maskMessage:   tt.fields.maskMessage,
				maskFormatter: tt.fields.maskFormatter,
				args:          tt.fields.args,
				stackTrace:    tt.fields.stackTrace,
				data:          tt.fields.data,
			}
			if got := e.RawMaskMessage(); got != tt.want {
				t.Errorf("errorWrapper.RawMaskMessage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_errorWrapper_Args(t *testing.T) {
	type fields struct {
		code          int
		codeString    string
		message       string
		category      ErrorCategory
		formatter     MessageFormatter
		isMasked      bool
		maskMessage   string
		maskFormatter MaskFormatter
		args          []interface{}
		stackTrace    []string
		data          ErrorData
	}
	tests := []struct {
		name   string
		fields fields
		want   []interface{}
	}{
		{
			name: "success",
			fields: fields{
				code:       100,
				codeString: "ErrTest",
				message:    "Test error message: %s",
				category:   ErrorCategory(1),
				args:       []interface{}{"Foo"},
				data: ErrorData{
					"foo": "bar",
				},
			},
			want: []interface{}{"Foo"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &errorWrapper{
				code:          tt.fields.code,
				codeString:    tt.fields.codeString,
				message:       tt.fields.message,
				category:      tt.fields.category,
				formatter:     tt.fields.formatter,
				isMasked:      tt.fields.isMasked,
				maskMessage:   tt.fields.maskMessage,
				maskFormatter: tt.fields.maskFormatter,
				args:          tt.fields.args,
				stackTrace:    tt.fields.stackTrace,
				data:          tt.fields.data,
			}
			if got := e.Args(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("errorWrapper.Args() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_errorWrapper_StackTrace(t *testing.T) {
	type fields struct {
		code          int
		codeString    string
		message       string
		category      ErrorCategory
		formatter     MessageFormatter
		isMasked      bool
		maskMessage   string
		maskFormatter MaskFormatter
		args          []interface{}
		stackTrace    []string
		data          ErrorData
	}
	tests := []struct {
		name   string
		fields fields
		want   []string
	}{
		{
			name: "success",
			fields: fields{
				code:       100,
				codeString: "ErrTest",
				message:    "Test error message: %s",
				category:   ErrorCategory(1),
				args:       []interface{}{"Foo"},
				data: ErrorData{
					"foo": "bar",
				},
				stackTrace: []string{
					"github.com/rapidashorg/errwrap/wrapper.go:1",
				},
			},
			want: []string{
				"github.com/rapidashorg/errwrap/wrapper.go:1",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &errorWrapper{
				code:          tt.fields.code,
				codeString:    tt.fields.codeString,
				message:       tt.fields.message,
				category:      tt.fields.category,
				formatter:     tt.fields.formatter,
				isMasked:      tt.fields.isMasked,
				maskMessage:   tt.fields.maskMessage,
				maskFormatter: tt.fields.maskFormatter,
				args:          tt.fields.args,
				stackTrace:    tt.fields.stackTrace,
				data:          tt.fields.data,
			}
			if got := e.StackTrace(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("errorWrapper.StackTrace() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_errorWrapper_Data(t *testing.T) {
	type fields struct {
		code          int
		codeString    string
		message       string
		category      ErrorCategory
		formatter     MessageFormatter
		isMasked      bool
		maskMessage   string
		maskFormatter MaskFormatter
		args          []interface{}
		stackTrace    []string
		data          ErrorData
	}
	tests := []struct {
		name   string
		fields fields
		want   ErrorData
	}{
		{
			name: "success",
			fields: fields{
				code:       100,
				codeString: "ErrTest",
				message:    "Test error message: %s",
				category:   ErrorCategory(1),
				args:       []interface{}{"Foo"},
				data: ErrorData{
					"foo": "bar",
				},
			},
			want: ErrorData{
				"foo": "bar",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &errorWrapper{
				code:          tt.fields.code,
				codeString:    tt.fields.codeString,
				message:       tt.fields.message,
				category:      tt.fields.category,
				formatter:     tt.fields.formatter,
				isMasked:      tt.fields.isMasked,
				maskMessage:   tt.fields.maskMessage,
				maskFormatter: tt.fields.maskFormatter,
				args:          tt.fields.args,
				stackTrace:    tt.fields.stackTrace,
				data:          tt.fields.data,
			}
			if got := e.Data(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("errorWrapper.Data() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_errorWrapper_Error(t *testing.T) {
	type fields struct {
		code          int
		codeString    string
		message       string
		category      ErrorCategory
		formatter     MessageFormatter
		isMasked      bool
		maskMessage   string
		maskFormatter MaskFormatter
		args          []interface{}
		stackTrace    []string
		data          ErrorData
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "success",
			fields: fields{
				code:       100,
				codeString: "ErrTest",
				message:    "Test error message: %s",
				category:   ErrorCategory(1),
				args:       []interface{}{"Foo"},
				data: ErrorData{
					"foo": "bar",
				},
				formatter: DefaultMessageFormatter,
			},
			want: "Test error message: Foo (100)",
		},
		{
			name: "success masked",
			fields: fields{
				code:       100,
				codeString: "ErrTest",
				message:    "Test error message: %s",
				category:   ErrorCategory(1),
				args:       []interface{}{"Foo"},
				data: ErrorData{
					"foo": "bar",
				},
				isMasked:      true,
				maskMessage:   "Test masked error message",
				maskFormatter: DefaultMaskFormatter,
			},
			want: "Test masked error message (100)",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &errorWrapper{
				code:          tt.fields.code,
				codeString:    tt.fields.codeString,
				message:       tt.fields.message,
				category:      tt.fields.category,
				formatter:     tt.fields.formatter,
				isMasked:      tt.fields.isMasked,
				maskMessage:   tt.fields.maskMessage,
				maskFormatter: tt.fields.maskFormatter,
				args:          tt.fields.args,
				stackTrace:    tt.fields.stackTrace,
				data:          tt.fields.data,
			}
			if got := e.Error(); got != tt.want {
				t.Errorf("errorWrapper.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_errorWrapper_Is(t *testing.T) {
	type fields struct {
		code          int
		codeString    string
		message       string
		category      ErrorCategory
		formatter     MessageFormatter
		isMasked      bool
		maskMessage   string
		maskFormatter MaskFormatter
		args          []interface{}
		stackTrace    []string
		data          ErrorData
	}
	type args struct {
		ed *ErrorDefinition
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "success equal",
			fields: fields{
				code:       100,
				codeString: "ErrTest",
				message:    "Test error message: %s",
				category:   ErrorCategory(1),
				args:       []interface{}{"Foo"},
				data: ErrorData{
					"foo": "bar",
				},
			},
			args: args{
				ed: &ErrorDefinition{
					code:       100,
					codeString: "ErrTest",
					category:   ErrorCategory(1),
				},
			},
			want: true,
		},
		{
			name: "success not equal",
			fields: fields{
				code:       100,
				codeString: "ErrTest",
				message:    "Test error message: %s",
				category:   ErrorCategory(1),
				args:       []interface{}{"Foo"},
				data: ErrorData{
					"foo": "bar",
				},
			},
			args: args{
				ed: &ErrorDefinition{
					code:       101,
					codeString: "ErrTestNew",
					category:   ErrorCategory(2),
				},
			},
			want: false,
		},
		{
			name: "success definition is nil",
			fields: fields{
				code:       100,
				codeString: "ErrTest",
				message:    "Test error message: %s",
				category:   ErrorCategory(1),
				args:       []interface{}{"Foo"},
				data: ErrorData{
					"foo": "bar",
				},
			},
			args: args{
				ed: nil,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &errorWrapper{
				code:          tt.fields.code,
				codeString:    tt.fields.codeString,
				message:       tt.fields.message,
				category:      tt.fields.category,
				formatter:     tt.fields.formatter,
				isMasked:      tt.fields.isMasked,
				maskMessage:   tt.fields.maskMessage,
				maskFormatter: tt.fields.maskFormatter,
				args:          tt.fields.args,
				stackTrace:    tt.fields.stackTrace,
				data:          tt.fields.data,
			}
			if got := e.Is(tt.args.ed); got != tt.want {
				t.Errorf("errorWrapper.Is() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_errorWrapper_ActualError(t *testing.T) {
	type fields struct {
		code          int
		codeString    string
		message       string
		category      ErrorCategory
		formatter     MessageFormatter
		isMasked      bool
		maskMessage   string
		maskFormatter MaskFormatter
		args          []interface{}
		stackTrace    []string
		data          ErrorData
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "success not masked",
			fields: fields{
				code:       100,
				codeString: "ErrTest",
				message:    "Test error message: %s",
				category:   ErrorCategory(1),
				args:       []interface{}{"Foo"},
				data: ErrorData{
					"foo": "bar",
				},
				formatter: DefaultMessageFormatter,
			},
			want: "Test error message: Foo (100)",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &errorWrapper{
				code:          tt.fields.code,
				codeString:    tt.fields.codeString,
				message:       tt.fields.message,
				category:      tt.fields.category,
				formatter:     tt.fields.formatter,
				isMasked:      tt.fields.isMasked,
				maskMessage:   tt.fields.maskMessage,
				maskFormatter: tt.fields.maskFormatter,
				args:          tt.fields.args,
				stackTrace:    tt.fields.stackTrace,
				data:          tt.fields.data,
			}
			if got := e.ActualError(); got != tt.want {
				t.Errorf("errorWrapper.ActualError() = %v, want %v", got, tt.want)
			}
		})
	}
}
