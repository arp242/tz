package tz

import (
	"fmt"
	"strings"
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		inC, inZ string
		want     string
		wantErr  string
	}{
		{"ID", "Asia/Makassar", "ID.Asia/Makassar", ""},
		{"", "Asia/Makassar", "ID.Asia/Makassar", ""},
		{"NL", "Asia/Makassar", "ID.Asia/Makassar", ""},

		{"GB", "UTC", ".UTC", ""},
		{"ID", "UTC", ".UTC", ""},
		{"", "UTC", ".UTC", ""},

		{"", "CET", "FR.Europe/Paris", ""},
		{"", "Asia/Saigon", "VN.Asia/Ho_Chi_Minh", ""},

		{"ID", "Asia/Denpasar", "", "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.inC+tt.inZ, func(t *testing.T) {
			z, err := New(tt.inC, tt.inZ)
			if !errorContains(err, tt.wantErr) {
				t.Fatalf("\nout:  %#v\nwant: %#v\n", err, tt.wantErr)
			}

			out := fmt.Sprintf("%s", z)
			if out != tt.want {
				t.Errorf("\nout:  %s\nwant: %s", out, tt.want)
			}
		})
	}
}

func errorContains(out error, want string) bool {
	if out == nil {
		return want == ""
	}
	if want == "" {
		return false
	}
	return strings.Contains(out.Error(), want)
}
