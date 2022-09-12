package main

import (
	"io"
	"reflect"
	"strings"
	"testing"
)

func Test_scanPairsFromFile(t *testing.T) {
	type args struct {
		reader io.Reader
	}
	tests := []struct {
		name string
		args args
		want map[string]string
	}{
		{"No values in file", args{reader: strings.NewReader("Hello world")}, map[string]string{}},
		{"Two values in file", args{reader: strings.NewReader("HELLO=WORLD\r\nWORLD=HELLO")}, map[string]string{"hello": "WORLD", "world": "HELLO"}},
		{"Value with white space", args{reader: strings.NewReader("\tHELLO  =\t\t\tWORLD")}, map[string]string{"hello": "WORLD"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := scanPairsFromFile(tt.args.reader); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("scanPairsFromFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_credentialsFromMap(t *testing.T) {
	type args struct {
		pairs map[string]string
	}
	tests := []struct {
		name    string
		args    args
		want    *Credentials
		wantErr bool
	}{
		{"Map with good values", args{pairs: map[string]string{"username": "root", "password": "root", "dbname": "db", "port": "3306", "host": "localhost"}}, &Credentials{Username: "root", Password: "root", DatabaseName: "db", Port: "3306", Host: "localhost"}, false},
		{"No username in map", args{pairs: map[string]string{"password": "root", "db": "db", "port": "3306", "host": "localhost"}}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := credentialsFromMap(tt.args.pairs)
			if (err != nil) != tt.wantErr {
				t.Errorf("credentialsFromMap() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("credentialsFromMap() = %v, want %v", got, tt.want)
			}
		})
	}
}
