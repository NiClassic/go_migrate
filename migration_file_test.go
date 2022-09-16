package main

import (
	"io"
	"reflect"
	"strings"
	"testing"
)

func TestReadStatementsFromReader(t *testing.T) {
	type args struct {
		r io.Reader
	}
	var emptySlice []string
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{"One sql statement", args{strings.NewReader("select * from test;")}, []string{"select * from test"}, false},
		{"Two sql statements", args{strings.NewReader("select * from test;\nselect id from demo;")}, []string{"select * from test", "select id from demo"}, false},
		{"Missing semicolon", args{strings.NewReader("select * from test")}, emptySlice, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReadStatementsFromReader(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadStatementsFromReader() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReadStatementsFromReader() = %v, want %v", got, tt.want)
			}
		})
	}
}
