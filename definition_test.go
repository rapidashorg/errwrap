package errwrap

import (
	"context"
	"fmt"
	"reflect"
	"testing"
)

func TestNewError(t *testing.T) {
	type args struct {
		code       int
		codeString string
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
				category:   ErrorCategory(1),
			},
			want: &ErrorDefinition{
				code:       100,
				codeString: "ErrTest",
				category:   ErrorCategory(1),

				formatter:     &DefaultMessageFormatter,
				maskMessage:   &DefaultMaskMessage,
				maskFormatter: &DefaultMaskFormatter,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewError(tt.args.code, tt.args.codeString, tt.args.category)

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
				category:   ErrorCategory(1),
			},
			want: &ErrorDefinition{
				code:       100,
				codeString: "ErrTest",
				category:   ErrorCategory(1),

				isMasked:      true,
				maskMessage:   &DefaultMaskMessage,
				formatter:     &DefaultMessageFormatter,
				maskFormatter: &DefaultMaskFormatter,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ed := NewError(tt.fields.code, tt.fields.codeString, tt.fields.category)

			got := ed.Masked()

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
				category:   ErrorCategory(1),
			},
			args: args{
				maskMessage: "Test masked error message",
			},
			want: &ErrorDefinition{
				code:       100,
				codeString: "ErrTest",
				category:   ErrorCategory(1),

				isMasked:      true,
				maskMessage:   stringPtr("Test masked error message"),
				formatter:     &DefaultMessageFormatter,
				maskFormatter: &DefaultMaskFormatter,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ed := NewError(tt.fields.code, tt.fields.codeString, tt.fields.category)

			got := ed.MaskedMessage(tt.args.maskMessage)

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
				category:   ErrorCategory(1),
			},
			args: args{
				fn: func(erw ErrorWrapper) string {
					return "test " + erw.RawMaskMessage()
				},
			},
			want: &ErrorDefinition{
				code:       100,
				codeString: "ErrTest",
				category:   ErrorCategory(1),

				isMasked:    true,
				formatter:   &DefaultMessageFormatter,
				maskMessage: &DefaultMaskMessage,
				maskFormatter: maskFormatterPtr(func(erw ErrorWrapper) string {
					return "test " + erw.RawMaskMessage()
				}),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ed := NewError(tt.fields.code, tt.fields.codeString, tt.fields.category)

			got := ed.MaskedFunction(tt.args.fn)

			gotErr := got.NewWithoutContext("error message")
			wantErr := tt.want.NewWithoutContext("error message")

			if !reflect.DeepEqual(gotErr.Error(), wantErr.Error()) {
				t.Errorf("ErrorDefinition.MaskedFunction() = %v, want %v", gotErr.Error(), wantErr.Error())
			}
		})
	}
}

func TestErrorDefinition_MessageFormatter(t *testing.T) {
	type fields struct {
		code          int
		codeString    string
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
				category:   ErrorCategory(1),
			},
			args: args{
				fn: func(msg string, erw ErrorWrapper) string {
					return fmt.Sprintf("test %s (%d)", msg, erw.Code())
				},
			},
			want: &ErrorDefinition{
				code:       100,
				codeString: "ErrTest",
				category:   ErrorCategory(1),

				formatter: messageFormatterPtr(func(msg string, erw ErrorWrapper) string {
					return fmt.Sprintf("test %s (%d)", msg, erw.Code())
				}),
				maskMessage:   &DefaultMaskMessage,
				maskFormatter: &DefaultMaskFormatter,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ed := NewError(tt.fields.code, tt.fields.codeString, tt.fields.category)

			got := ed.MessageFormatter(tt.args.fn)

			gotErr := got.NewWithoutContext("error message")
			wantErr := tt.want.NewWithoutContext("error message")

			if !reflect.DeepEqual(gotErr.Error(), wantErr.Error()) {
				t.Errorf("ErrorDefinition.MaskedFunction() = %v, want %v", gotErr.Error(), wantErr.Error())
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
		rawMessage string
		args       []interface{}
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
				args:       []interface{}{"Foo"},
				rawMessage: "Test error message",
			},
			want: &errorWrapper{
				code:       100,
				codeString: "ErrTest",
				message:    "Test error message",
				category:   ErrorCategory(1),
				args:       []interface{}{"Foo"},

				formatter:     nil,
				maskMessage:   DefaultMaskMessage,
				maskFormatter: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ed := NewError(tt.fields.code, tt.fields.codeString, tt.fields.category)

			got := ed.NewWithoutContext(tt.args.rawMessage, tt.args.args...)
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
		isMasked      bool
		maskMessage   string
		maskFormatter MaskFormatter
		category      ErrorCategory
	}
	type args struct {
		ctx        context.Context
		args       []interface{}
		rawMessage string
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
				category:   ErrorCategory(1),
			},
			args: args{
				ctx: InjectErrorData(context.Background(), ErrorData{
					"foo": "bar",
				}),
				args:       []interface{}{"Foo"},
				rawMessage: "Test error message",
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

				formatter:     nil,
				maskMessage:   DefaultMaskMessage,
				maskFormatter: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ed := NewError(tt.fields.code, tt.fields.codeString, tt.fields.category)

			got := ed.New(tt.args.ctx, tt.args.rawMessage, tt.args.args...)
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
