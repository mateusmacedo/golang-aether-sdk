package err

import (
	"fmt"
	"reflect"
	"testing"
)

func Test_baseError_Error(t *testing.T) {
	type fields struct {
		msg   string
		cause error
		stack []string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Ensure that the error message is returned",
			fields: fields{
				msg: "error message",
			},
			want: "error message",
		},
		{
			name: "Ensure that the error message is returned with the cause",
			fields: fields{
				msg:   "error message",
				cause: fmt.Errorf("cause message"),
			},
			want: "error message: cause message",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &baseError{
				what:   tt.fields.msg,
				cause: tt.fields.cause,
				stack: tt.fields.stack,
			}
			if got := e.Error(); got != tt.want {
				t.Errorf("baseError.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_baseError_Unwrap(t *testing.T) {
	type fields struct {
		msg   string
		cause error
		stack []string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "Ensure that the cause is returned when it is not nil",
			fields: fields{
				cause: fmt.Errorf("cause message"),
			},
			wantErr: true,
		},
		{
			name: "Ensure that the cause is not returned when it is nil",
			fields: fields{
				cause: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &baseError{
				what:   tt.fields.msg,
				cause: tt.fields.cause,
				stack: tt.fields.stack,
			}
			if err := e.Unwrap(); (err != nil) != tt.wantErr {
				t.Errorf("baseError.Unwrap() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_baseError_WithStack(t *testing.T) {
	type args struct {
		s []string
	}
	type fields struct {
		msg   string
		cause error
		stack []string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   CustomError
	}{
		{
			name: "Ensure that the stack is set and cannot be changed",
			fields: fields{
				stack: []string{"stack"},
			},
			args: args{
				s: []string{"new stack"},
			},
			want: &baseError{
				stack: []string{"stack"},
			},
		},
		{
			name: "Ensure that the stack is set when it is nil",
			fields: fields{
				stack: nil,
			},
			args: args{
				s: []string{"new stack"},
			},
			want: &baseError{
				stack: []string{"new stack"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &baseError{
				what:   tt.fields.msg,
				cause: tt.fields.cause,
				stack: tt.fields.stack,
			}
			if got := e.WithStack(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("baseError.WithStack() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_captureStack(t *testing.T) {
	tests := []struct {
		name string
		want []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := captureStack(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("captureStack() = %v, want %v", got, tt.want)
			}
		})
	}
}
