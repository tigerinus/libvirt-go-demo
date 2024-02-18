package util_test

import (
	"testing"

	"github.com/tigerinus/libvirt-go-demo/util"
)

func TestReplaceRegex(t *testing.T) {
	old := "[{|}~[\\]^':; <=>?@!\"#$%`()+/.,*&]"
	type args struct {
		str         string
		old         string
		replacement string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test1",
			args: args{str: "My@VM!Name", old: old, replacement: ""},
			want: "MyVMName",
		},
		{
			name: "test2",
			args: args{str: "Example#1", old: old, replacement: ""},
			want: "Example1",
		},
		{
			name: "test3",
			args: args{str: "Hello, World!", old: old, replacement: ""},
			want: "HelloWorld",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := util.ReplaceRegex(tt.args.str, tt.args.old, tt.args.replacement); got != tt.want {
				t.Errorf("ReplaceRegex() = %v, want %v", got, tt.want)
			}
		})
	}
}
