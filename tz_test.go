package tz

import (
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
		{"NL", "", "NL.Europe/Brussels", ""},            // Just country
		{"ID", "", "ID.Asia/Jakarta", ""},               // Just country

		{"GB", "UTC", ".UTC", ""}, // Various way of sending UTC
		{"ID", "UTC", ".UTC", ""},
		{"", "UTC", ".UTC", ""},

		{"", "CET", "BE.Europe/Brussels", ""},          // Alias
		{"", "Asia/Saigon", "VN.Asia/Ho_Chi_Minh", ""}, // Alias

		{"ID", "Asia/Denpasar", "", "unknown"}, // Doesn't exist

		{"", "Etc/GMT", ".UTC", ""},
		{"", "Etc/UTC", ".UTC", ""},
		{"", "Etc/Unknown", ".UTC", ""},
		{"", "Etc/XXX", "", "invalid Etc/"},

		{"SG", "Etc/GMT-8", "SG.Asia/Singapore", ""},
		{"IT", "Etc/GMT-1", "IT.Europe/Rome", ""},
		{"MX", "Etc/GMT+7", "MX.America/Ciudad_Juarez", ""},
		{"JO", "Etc/GMT-3", "JO.Asia/Amman", ""},
		{"SG", "Etc/GMT-10", "SG.Asia/Singapore", ""}, // Invalid country/offset
		{"", "Etc/GMT-10", "", "unknown timezone"},
	}

	for _, tt := range tests {
		t.Run(tt.inC+tt.inZ, func(t *testing.T) {
			z, err := New(tt.inC, tt.inZ)
			if !errorContains(err, tt.wantErr) {
				t.Fatalf("\nout:  %#v\nwant: %#v\n", err, tt.wantErr)
			}

			out := z.String()
			if out != tt.want {
				t.Errorf("\nout:  %s\nwant: %s", out, tt.want)
			}

			if tt.wantErr != "" {
				t.Run("MustNew", func(t *testing.T) {
					defer func() {
						if recover() == nil {
							t.Error("recover() is nil")
						}
					}()
					z := MustNew(tt.inC, tt.inZ)
					out := z.String()
					if out != tt.want {
						t.Errorf("\nout:  %s\nwant: %s", out, tt.want)
					}
				})
			}
		})
	}
}

func TestOffsetRFC3999(t *testing.T) {
	tests := []struct {
		in   *Zone
		want string
	}{
		{nil, "UTC"},
		{MustNew("", "UTC"), "UTC"},
		{MustNew("", "America/Sao_Paulo"), "-03:00"},
		{MustNew("", "Australia/Darwin"), "+09:30"},
	}

	for _, tt := range tests {
		t.Run(tt.in.String(), func(t *testing.T) {
			out := tt.in.OffsetRFC3339()
			if out != tt.want {
				t.Errorf("\nout:  %s\nwant: %s", out, tt.want)
			}
		})
	}
}

func TestOffsetDisplay(t *testing.T) {
	tests := []struct {
		in   *Zone
		want string
	}{
		{nil, "UTC"},
		{MustNew("", "UTC"), "UTC"},
		{MustNew("", "America/Sao_Paulo"), "UTC -3:00"},
		{MustNew("", "Australia/Darwin"), "UTC +9:30"},
	}

	for _, tt := range tests {
		t.Run(tt.in.String(), func(t *testing.T) {
			out := tt.in.OffsetDisplay()
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
