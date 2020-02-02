//go:generate sh -c "go run gen.go > list.go && gofmt -w list.go"

// Package tz contains timezone lists.
package tz

import (
	"database/sql/driver"
	"fmt"
	"strings"
	"time"
)

// Zone represents a time zone.
type Zone struct {
	*time.Location

	CountryCode string // ID
	Zone        string // Asia/Makassar
	Abbr        string // WITA
	CountryName string // Indonesia
	Comments    string // Borneo (east, south); Sulawesi/Celebes, Bali, Nusa Tengarra; Timor (west)

	display string // cached Display()
}

// New timezone from country code and zone name. If the country code is empty it
// will load the first zone found.
func New(ccode, zone string) (*Zone, error) {
	_, err := time.LoadLocation(zone)
	if err != nil {
		return nil, err
	}

	for _, z := range Zones {
		if (ccode == "" || z.CountryCode == ccode) && z.Zone == zone {
			return z, nil
		}
	}

	return nil, fmt.Errorf("unknown timezone: %q %q", ccode, zone)
}

// Loc gets the time.Location, or UTC if it's not set.
func (t *Zone) Loc() *time.Location {
	if t == nil || t.Location == nil {
		return time.UTC
	}
	return t.Location
}

// Display a human-readable description of the timezone, for e.g. <option>.
func (t *Zone) Display() string {
	if t == nil {
		return ""
	}

	if t.display == "" {
		var b strings.Builder
		b.WriteString(t.CountryName)
		b.WriteString(": ")
		b.WriteString(t.Zone)
		if t.Abbr != "" {
			b.WriteString(" (")
			b.WriteString(t.Abbr)
			b.WriteString(")")
		}
		if t.Comments != "" {
			b.WriteString(" – ")
			b.WriteString(t.Comments)
		}
		t.display = b.String()
	}
	return t.display
}

// String is an unique representation for this timezone.
func (t *Zone) String() string {
	if t == nil {
		return ""
	}
	return fmt.Sprintf("%s.%s", t.CountryCode, t.Zone)
}

// MarshalText converts the data to a human readable representation.
func (t Zone) MarshalText() ([]byte, error) { return []byte(t.String()), nil }

// UnmarshalText parses text in to the Go data structure.
func (t *Zone) UnmarshalText(v []byte) error {
	return t.Scan(v)
}

// Offset gets the timezone offset in minutes.
func (t *Zone) Offset() int {
	if t == nil || t.Location == nil {
		return 0
	}
	now := time.Now().In(t.Location)
	_, offset := now.Zone()
	return offset / 60
}

// Value implements the SQL Value function to determine what to store in the DB.
func (t Zone) Value() (driver.Value, error) {
	if t.CountryCode == "" || t.Zone == "" {
		return nil, fmt.Errorf("CountryCode (%q) and Zone (%q) must be set", t.CountryCode, t.Zone)
	}
	return t.String(), nil
}

// Scan converts the data returned from the DB into the struct.
func (t *Zone) Scan(v interface{}) error {
	var vv string
	switch v.(type) {
	case string:
		vv = v.(string)
	case []byte:
		vv = string(v.([]byte))
	}
	s := strings.SplitN(vv, ".", 2)
	if len(s) != 2 {
		return fmt.Errorf("invalid value: %q", vv)
	}

	z, err := New(s[0], s[1])
	*t = *z
	return err
}
