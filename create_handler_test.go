package main

import (
	"testing"
)

func Test_generateUpAndDownFileNames(t *testing.T) {
	type args struct {
		prefix string
		name   string
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 string
	}{
		{"Parse test migration name", args{prefix: "AAA", name: "giraffe"}, "AAA_giraffe.up.sql", "AAA_giraffe.down.sql"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := generateUpAndDownFileNames(tt.args.prefix, tt.args.name)
			if got != tt.want {
				t.Errorf("generateUpAndDownFileNames() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("generateUpAndDownFileNames() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
