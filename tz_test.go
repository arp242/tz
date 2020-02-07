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
		{"ID", "Asia/Makassar", "ID.Asia/Makassar", ""}, // Country+Zone
		{"", "Asia/Makassar", "ID.Asia/Makassar", ""},   // Just zone
		{"NL", "Asia/Makassar", "ID.Asia/Makassar", ""}, // Zone with wrong country
		{"NL", "", "NL.Europe/Amsterdam", ""},           // Just country
		{"ID", "", "ID.Asia/Jayapura", ""},              // Just country

		{"GB", "UTC", ".UTC", ""}, // Various way of sending UTC
		{"ID", "UTC", ".UTC", ""},
		{"", "UTC", ".UTC", ""},

		{"", "CET", "FR.Europe/Paris", ""},             // Alias
		{"", "Asia/Saigon", "VN.Asia/Ho_Chi_Minh", ""}, // Alias

		{"ID", "Asia/Denpasar", "", "unknown"}, // Doesn't exist
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
