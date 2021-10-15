package errwrap

import (
	"context"
	"reflect"
	"testing"
)

func TestNewError(t *testing.T) {
	type args struct {
		code       int
		codeString string
		message    string
		category   ErrorCategory
	}
	tests := []struct {
		name string
		args args
		want *ErrorDefinition
	}{
		{
			name: "success",
			args: args{
				code:       100,
				codeString: "ErrTest",
				message:    "Test error message",
				category:   ErrorCategory(1),
			},
			want: &ErrorDefinition{
				code:       100,
				codeString: "ErrTest",
				message:    "Test error message",
				category:   ErrorCategory(1),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewError(tt.args.code, tt.args.codeString, tt.args.message, tt.args.category)
			got.maskFormatter = nil
			got.formatter = nil

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestErrorDefinition_Masked(t *testing.T) {
	type fields struct {
		code          int
		codeString    string
		message       string
		isMasked      bool
		maskMessage   string
		maskFormatter MaskFormatter
		category      ErrorCategory
	}
	tests := []struct {
		name   string
		fields fields
		want   *ErrorDefinition
	}{
		{
			name: "success",
			fields: fields{
				code:       100,
				codeString: "ErrTest",
				message:    "Test error message",
				category:   ErrorCategory(1),
			},
			want: &ErrorDefinition{
				code:       100,
				codeString: "ErrTest",
				message:    "Test error message",
				category:   ErrorCategory(1),

				isMasked:    true,
				maskMessage: DefaultMaskMessage,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ed := &ErrorDefinition{
				code:          tt.fields.code,
				codeString:    tt.fields.codeString,
				message:       tt.fields.message,
				isMasked:      tt.fields.isMasked,
				maskMessage:   tt.fields.maskMessage,
				maskFormatter: tt.fields.maskFormatter,
				category:      tt.fields.category,
			}

			got := ed.Masked()
			got.maskFormatter = nil
			got.formatter = nil

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ErrorDefinition.Masked() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestErrorDefinition_MaskedMessage(t *testing.T) {
	type fields struct {
		code          int
		codeString    string
		message       string
		isMasked      bool
		maskMessage   string
		maskFormatter MaskFormatter
		category      ErrorCategory
	}
	type args struct {
		maskMessage string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *ErrorDefinition
	}{
		{
			name: "success",
			fields: fields{
				code:       100,
				codeString: "ErrTest",
				message:    "Test error message",
				category:   ErrorCategory(1),
			},
			args: args{
				maskMessage: "Test masked error message",
			},
			want: &ErrorDefinition{
				code:       100,
				codeString: "ErrTest",
				message:    "Test error message",
				category:   ErrorCategory(1),

				isMasked:    true,
				maskMessage: "Test masked error message",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ed := &ErrorDefinition{
				code:          tt.fields.code,
				codeString:    tt.fields.codeString,
				message:       tt.fields.message,
				isMasked:      tt.fields.isMasked,
				maskMessage:   tt.fields.maskMessage,
				maskFormatter: tt.fields.maskFormatter,
				category:      tt.fields.category,
			}

			got := ed.MaskedMessage(tt.args.maskMessage)
			got.maskFormatter = nil
			got.formatter = nil

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ErrorDefinition.MaskedMessage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestErrorDefinition_MaskedFunction(t *testing.T) {
	type fields struct {
		code          int
		codeString    string
		message       string
		isMasked      bool
		maskMessage   string
		maskFormatter MaskFormatter
		category      ErrorCategory
	}
	type args struct {
		fn MaskFormatter
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *ErrorDefinition
	}{
		{
			name: "success",
			fields: fields{
				code:       100,
				codeString: "ErrTest",
				message:    "Test error message",
				category:   ErrorCategory(1),
			},
			args: args{
				fn: nil,
			},
			want: &ErrorDefinition{
				code:       100,
				codeString: "ErrTest",
				message:    "Test error message",
				category:   ErrorCategory(1),

				isMasked: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ed := &ErrorDefinition{
				code:          tt.fields.code,
				codeString:    tt.fields.codeString,
				message:       tt.fields.message,
				isMasked:      tt.fields.isMasked,
				maskMessage:   tt.fields.maskMessage,
				maskFormatter: tt.fields.maskFormatter,
				category:      tt.fields.category,
			}

			got := ed.MaskedFunction(tt.args.fn)
			got.maskFormatter = nil
			got.formatter = nil

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ErrorDefinition.MaskedFunction() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestErrorDefinition_MessageFormatter(t *testing.T) {
	type fields struct {
		code          int
		codeString    string
		message       string
		isMasked      bool
		formatter     MessageFormatter
		maskMessage   string
		maskFormatter MaskFormatter
		category      ErrorCategory
	}
	type args struct {
		fn MessageFormatter
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *ErrorDefinition
	}{
		{
			name: "success",
			fields: fields{
				code:       100,
				codeString: "ErrTest",
				message:    "Test error message",
				category:   ErrorCategory(1),
			},
			args: args{
				fn: nil,
			},
			want: &ErrorDefinition{
				code:       100,
				codeString: "ErrTest",
				message:    "Test error message",
				category:   ErrorCategory(1),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ed := &ErrorDefinition{
				code:          tt.fields.code,
				codeString:    tt.fields.codeString,
				message:       tt.fields.message,
				isMasked:      tt.fields.isMasked,
				formatter:     tt.fields.formatter,
				maskMessage:   tt.fields.maskMessage,
				maskFormatter: tt.fields.maskFormatter,
				category:      tt.fields.category,
			}

			got := ed.MessageFormatter(tt.args.fn)
			got.maskFormatter = nil
			got.formatter = nil

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ErrorDefinition.MessageFormatter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestErrorDefinition_NewWithoutContext(t *testing.T) {
	type fields struct {
		code          int
		codeString    string
		message       string
		isMasked      bool
		maskMessage   string
		maskFormatter MaskFormatter
		category      ErrorCategory
	}
	type args struct {
		args []interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   ErrorWrapper
	}{
		{
			name: "success",
			fields: fields{
				code:       100,
				codeString: "ErrTest",
				message:    "Test error message",
				category:   ErrorCategory(1),
			},
			args: args{
				args: []interface{}{"Foo"},
			},
			want: &errorWrapper{
				code:       100,
				codeString: "ErrTest",
				message:    "Test error message",
				category:   ErrorCategory(1),
				args:       []interface{}{"Foo"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ed := &ErrorDefinition{
				code:          tt.fields.code,
				codeString:    tt.fields.codeString,
				message:       tt.fields.message,
				isMasked:      tt.fields.isMasked,
				maskMessage:   tt.fields.maskMessage,
				maskFormatter: tt.fields.maskFormatter,
				category:      tt.fields.category,
			}

			got := ed.NewWithoutContext(tt.args.args...)
			if g, ok := got.(*errorWrapper); ok {
				g.stackTrace = nil
				g.maskFormatter = nil
				g.formatter = nil
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ErrorDefinition.NewWithoutContext() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestErrorDefinition_New(t *testing.T) {
	type fields struct {
		code          int
		codeString    string
		message       string
		isMasked      bool
		maskMessage   string
		maskFormatter MaskFormatter
		category      ErrorCategory
	}
	type args struct {
		ctx  context.Context
		args []interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   ErrorWrapper
	}{
		{
			name: "success",
			fields: fields{
				code:       100,
				codeString: "ErrTest",
				message:    "Test error message",
				category:   ErrorCategory(1),
			},
			args: args{
				ctx: InjectErrorData(context.Background(), ErrorData{
					"foo": "bar",
				}),
				args: []interface{}{"Foo"},
			},
			want: &errorWrapper{
				code:       100,
				codeString: "ErrTest",
				message:    "Test error message",
				category:   ErrorCategory(1),
				args:       []interface{}{"Foo"},
				data: ErrorData{
					"foo": "bar",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ed := &ErrorDefinition{
				code:          tt.fields.code,
				codeString:    tt.fields.codeString,
				message:       tt.fields.message,
				isMasked:      tt.fields.isMasked,
				maskMessage:   tt.fields.maskMessage,
				maskFormatter: tt.fields.maskFormatter,
				category:      tt.fields.category,
			}

			got := ed.New(tt.args.ctx, tt.args.args...)
			if g, ok := got.(*errorWrapper); ok {
				g.stackTrace = nil
				g.maskFormatter = nil
				g.formatter = nil
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ErrorDefinition.New() = %v, want %v", got, tt.want)
			}
		})
	}
}
