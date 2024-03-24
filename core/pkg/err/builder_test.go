package err

import (
	"reflect"
	"testing"
)

func TestErrorBuilder(t *testing.T) {
	t.Run("New", func(t *testing.T) {
		type args struct {
			msg string
		}
		tests := []struct {
			name string
			args args
			want *ErrorBuilder
		}{
			{
				name: "Test New",
				args: args{
					msg: "Test New Error",
				},
				want: &ErrorBuilder{
					err: &baseError{what: "Test New Error"},
				},
			},
			{
				name: "Test Another New with different message",
				args: args{
					msg: "Another New Error",
				},
				want: &ErrorBuilder{
					err: &baseError{what: "Another New Error"},
				},
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got := New(tt.args.msg)
				if got.err.what != tt.want.err.what {
					t.Errorf("New() = %v, want %v", got, tt.want)
				}
			})
		}
	})
}

func TestErrorBuilder_Wrap(t *testing.T) {
	type fields struct {
		err *baseError
	}
	type args struct {
		cause error
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *ErrorBuilder
	}{
		{
			name: "Test Wrap with nil cause",
			fields: fields{
				err: &baseError{what: "Test Wrap Error"},
			},
			args: args{
				cause: nil,
			},
			want: &ErrorBuilder{
				err: &baseError{
					what: "Test Wrap Error", cause: nil,
				},
			},
		},
		{
			name: "Test Wrap with cause",
			fields: fields{
				err: &baseError{
					what: "Test Wrap Error And Cause",
				},
			},
			args: args{
				cause: &baseError{what: "Cause Error Message"},
			},
			want: &ErrorBuilder{
				err: &baseError{
					what: "Test Wrap Error And Cause",
					cause: &baseError{what: "Cause Error Message"},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &ErrorBuilder{
				err: tt.fields.err,
			}
			if got := b.Wrap(tt.args.cause); !reflect.DeepEqual(got.Build(), tt.want.Build()){
				t.Errorf("ErrorBuilder.Wrap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestErrorBuilder_Build(t *testing.T) {
	type fields struct {
		err *baseError
	}
	tests := []struct {
		name   string
		fields fields
		want   CustomError
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &ErrorBuilder{
				err: tt.fields.err,
			}
			if got := b.Build(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ErrorBuilder.Build() = %v, want %v", got, tt.want)
			}
		})
	}
}
