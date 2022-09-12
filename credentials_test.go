package main

import "testing"

func TestCredentials_ToDSN(t *testing.T) {
	tests := []struct {
		name string
		c    *Credentials
		want string
	}{
		{"DSN with text host", NewCredentials("root", "root", "demo", "3306", "localhost"), "root:root@tcp(localhost:3306)/demo"},
		{"DSN with ip host", NewCredentials("root", "root", "demo", "3306", "127.0.0.1"), "root:root@tcp(127.0.0.1:3306)/demo"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.ToDSN(); got != tt.want {
				t.Errorf("Credentials.ToDSN() = %v, want %v", got, tt.want)
			}
		})
	}
}
