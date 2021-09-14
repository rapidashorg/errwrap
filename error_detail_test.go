package errors

import (
	"reflect"
	"testing"
)

func TestErrorWrapper_Error(t *testing.T) {
	type fields struct {
		code        int
		message     string
		maskMessage string
		args        []interface{}
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "success 1",
			fields: fields{
				code:    1,
				message: "without format",
				args:    []interface{}{},
			},
			want: "without format (1)",
		},
		{
			name: "success 2",
			fields: fields{
				code:    1,
				message: "with %s",
				args:    []interface{}{"format"},
			},
			want: "with format (1)",
		},
		{
			name: "success 3",
			fields: fields{
				code:        1,
				message:     "with %s",
				args:        []interface{}{"format"},
				maskMessage: DefaultMaskMessage,
			},
			want: "Sorry, there are internal server error occured, please try again later. (1)",
		},
		{
			name: "success 4",
			fields: fields{
				code:        1,
				message:     "with %s",
				args:        []interface{}{"format"},
				maskMessage: "a-mask-message",
			},
			want: "a-mask-message",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &ErrorWrapper{
				Code:    tt.fields.code,
				Message: tt.fields.message,
				Args:    tt.fields.args,
			}
			if tt.fields.maskMessage != "" {
				e.MaskMessage = tt.fields.maskMessage
				e.IsMasked = true
			}

			if got := e.Error(); got != tt.want {
				t.Errorf("ErrorWrapper.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestErrorWrapper_Is(t *testing.T) {
	type fields struct {
		Code           int
		Message        string
		IsMasked       bool
		MaskMessage    string
		HTTPErrorCode  int
		AppendToLog    bool
		Args           []interface{}
		StackTrace     []string
		AdditionalData []*ErrorWrapperAdditionalData
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
			name: "success",
			fields: fields{
				Code: 101,
			},
			args: args{
				ed: &ErrorDefinition{
					code: 101,
				},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &ErrorWrapper{
				Code:           tt.fields.Code,
				Message:        tt.fields.Message,
				IsMasked:       tt.fields.IsMasked,
				MaskMessage:    tt.fields.MaskMessage,
				HTTPErrorCode:  tt.fields.HTTPErrorCode,
				Args:           tt.fields.Args,
				StackTrace:     tt.fields.StackTrace,
				AdditionalData: tt.fields.AdditionalData,
			}
			if got := e.Is(tt.args.ed); got != tt.want {
				t.Errorf("ErrorWrapper.Is() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestErrorWrapper_WithData(t *testing.T) {
	type fields struct {
		Code           int
		Message        string
		IsMasked       bool
		MaskMessage    string
		HTTPErrorCode  int
		AppendToLog    bool
		Args           []interface{}
		StackTrace     []string
		AdditionalData []*ErrorWrapperAdditionalData
	}
	type args struct {
		data map[string]interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *ErrorWrapper
	}{
		{
			name: "success",
			args: args{
				data: map[string]interface{}{
					"key": "value",
				},
			},
			want: &ErrorWrapper{
				AdditionalData: []*ErrorWrapperAdditionalData{
					{
						Data: map[string]interface{}{
							"key": "value",
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &ErrorWrapper{
				Code:           tt.fields.Code,
				Message:        tt.fields.Message,
				IsMasked:       tt.fields.IsMasked,
				MaskMessage:    tt.fields.MaskMessage,
				HTTPErrorCode:  tt.fields.HTTPErrorCode,
				Args:           tt.fields.Args,
				StackTrace:     tt.fields.StackTrace,
				AdditionalData: tt.fields.AdditionalData,
			}

			_ = e.WithData(tt.args.data)
			for _, v := range e.AdditionalData {
				v.Line = ""
			}

			if got := e; !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ErrorWrapper.WithData() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestErrorWrapper_getCallLine(t *testing.T) {
	type fields struct {
		Code           int
		Message        string
		IsMasked       bool
		MaskMessage    string
		HTTPErrorCode  int
		AppendToLog    bool
		Args           []interface{}
		StackTrace     []string
		AdditionalData []*ErrorWrapperAdditionalData
	}
	type args struct {
		offset int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "error call offset is large",
			args: args{
				offset: 10000,
			},
			want: "",
		},
		{
			name: "success",
			args: args{
				offset: 0,
			},
			want: "a-line",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &ErrorWrapper{
				Code:           tt.fields.Code,
				Message:        tt.fields.Message,
				IsMasked:       tt.fields.IsMasked,
				MaskMessage:    tt.fields.MaskMessage,
				HTTPErrorCode:  tt.fields.HTTPErrorCode,
				Args:           tt.fields.Args,
				StackTrace:     tt.fields.StackTrace,
				AdditionalData: tt.fields.AdditionalData,
			}

			got := e.getCallLine(tt.args.offset)
			if got != "" {
				got = "a-line"
			}

			if got != tt.want {
				t.Errorf("ErrorWrapper.getCallLine() = %v, want %v", got, tt.want)
			}
		})
	}
}
