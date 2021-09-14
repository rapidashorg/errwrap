package errors

import (
	"reflect"
	"testing"
)

func TestErrorDefinition_New(t *testing.T) {
	type fields struct {
		code          int
		message       string
		isMasked      bool
		maskMessage   string
		httpErrorCode int
		args          []interface{}
	}
	tests := []struct {
		name   string
		fields fields
		want   *ErrorWrapper
	}{
		{
			name: "success",
			fields: fields{
				args: []interface{}{},
			},
			want: &ErrorWrapper{
				Args:       []interface{}{},
				StackTrace: []string{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ed := &ErrorDefinition{
				code:          tt.fields.code,
				message:       tt.fields.message,
				isMasked:      tt.fields.isMasked,
				maskMessage:   tt.fields.maskMessage,
				httpErrorCode: tt.fields.httpErrorCode,
			}

			args := tt.fields.args
			if args == nil {
				args = make([]interface{}, 0)
			}

			ew := ed.New(args...)
			ew.StackTrace = []string{}

			if !reflect.DeepEqual(ew, tt.want) {
				t.Errorf("ErrorDefinition.New() = %v, want %v", ew, tt.want)
			}
		})
	}
}

func TestNewError(t *testing.T) {
	type args struct {
		code     int
		message  string
		httpCode *int
	}
	tests := []struct {
		name string
		args args
		want *ErrorDefinition
	}{
		{
			name: "success 1",
			args: args{
				code:     1,
				message:  "test",
				httpCode: nil,
			},
			want: &ErrorDefinition{
				code:          1,
				message:       "test",
				httpErrorCode: DefaultHTTPErrorCode,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewError(tt.args.code, tt.args.message, tt.args.httpCode); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewMaskedError(t *testing.T) {
	type args struct {
		code        int
		message     string
		httpCode    *int
		maskMessage *string
	}
	tests := []struct {
		name string
		args args
		want *ErrorDefinition
	}{
		{
			name: "success 1",
			args: args{
				code:        1,
				message:     "test",
				httpCode:    nil,
				maskMessage: nil,
			},
			want: &ErrorDefinition{
				code:          1,
				message:       "test",
				httpErrorCode: DefaultMaskedHTTPErrorCode,
				isMasked:      true,
				maskMessage:   DefaultMaskMessage,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewMaskedError(tt.args.code, tt.args.message, tt.args.httpCode, tt.args.maskMessage); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewMaskedError() = %v, want %v", got, tt.want)
			}
		})
	}
}
