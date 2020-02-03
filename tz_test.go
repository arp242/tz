package tz

import (
	"fmt"
	"strings"
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		inC, inZ string
		want     *Zone
		wantErr  string
	}{
		{"ID", "Asia/Makassar", &Zone{
			CountryCode: "ID",
			Zone:        "Asia/Makassar",
			Abbr:        "WITA",
			CountryName: "Indonesia",
			Comments:    "Borneo (east, south); Sulawesi/Celebes, Bali, Nusa Tengarra; Timor (west)",
		}, ""},

		{"", "Asia/Makassar", &Zone{
			CountryCode: "ID",
			Zone:        "Asia/Makassar",
			Abbr:        "WITA",
			CountryName: "Indonesia",
			Comments:    "Borneo (east, south); Sulawesi/Celebes, Bali, Nusa Tengarra; Timor (west)",
		}, ""},

		{"GB", "UTC", UTC, ""},
		{"ID", "UTC", UTC, ""},
		{"", "UTC", UTC, ""},

		{"NL", "Asia/Makassar", nil, "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.inC+tt.inZ, func(t *testing.T) {
			z, err := New(tt.inC, tt.inZ)
			if !errorContains(err, tt.wantErr) {
				t.Fatalf("\nout:  %#v\nwant: %#v\n", err, tt.wantErr)
			}

			out := fmt.Sprintf("%v", z)
			want := fmt.Sprintf("%v", tt.want)
			if out != want {
				t.Errorf("\nout:  %s\nwant: %s", out, want)
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
