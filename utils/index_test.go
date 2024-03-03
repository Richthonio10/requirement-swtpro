package utils

import (
	"errors"
	"testing"
)

func Test_ErrorMessage(t *testing.T) {
	type args struct {
		input error
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "no error",
			args: args{
				input: nil,
			},
			want: "",
		},
		{
			name: "error",
			args: args{
				input: errors.New("expected error"),
			},
			want: "expected error",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ErrorMessage(tt.args.input)
			if got != tt.want {
				t.Errorf("Result When ErrorMessage() %s, detail = %s", got, tt.want)
			}
		})
	}
}
