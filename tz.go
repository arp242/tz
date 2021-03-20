//go:generate sh -c "go run gen.go > list.go && gofmt -w list.go"

// Package tz contains timezone lists.
package tz

import (
	"database/sql/driver"
	"fmt"
	"math"
	"os"
	"strings"
	"sync"
	"time"
)

// Add time.Location; also serves as sanity-check on startup.
// func init() {
// 	for _, z := range Zones {
// 		var err error
// 		z.Location, err = time.LoadLocation(z.Zone)
// 		if err != nil {
// 			if strings.Contains(err.Error(), "unknown time zone") {
// 				fmt.Fprintf(os.Stderr, "warning: tz.init: %s; you probably need to update your tzdata/zoneinfo\n", err)
// 			}
// 		}
// 	}
// }

// Zone represents a time zone.
type Zone struct {
	*time.Location

	CountryCode string   // ID
	Zone        string   // Asia/Makassar
	Abbr        []string // WITA – the correct abbreviation may change depending on the time of year (i.e. CET and CEST, depending on DST).
	CountryName string   // Indonesia
	Comments    string   // Borneo (east, south); Sulawesi/Celebes, Bali, Nusa Tengarra; Timor (west)

	display string // cached Display()
}

var loadLocationOnce sync.Once

// New timezone from country code and zone name. The country code is only
// informative, and may be blank or wrong, in which case it will load the first
// zone found.
func New(ccode, zone string) (*Zone, error) {
	// Add time.Location to all the zones. This is about 68k memory without the
	// loaded zones, and 670k with. Not super huge, but kinda large. Also takes
	// about 12ms on my laptop.
	loadLocationOnce.Do(func() {
		for _, z := range Zones {
			var err error
			z.Location, err = time.LoadLocation(z.Zone)
			if err != nil {
				if strings.Contains(err.Error(), "unknown time zone") {
					fmt.Fprintf(os.Stderr, "warning: zgo.at/tz: %s; you probably need to update your tzdata or zoneinfo\n", err)
				}
			}
		}
	})

	if zone == "UTC" {
		return UTC, nil
	}
	if a, ok := aliases[zone]; ok {
		zone = a
	}

	// No zone name but country given; just get the first zone for that country,
	// which is better than nothing.
	if zone == "" && ccode != "" {
		for _, z := range Zones {
			if z.CountryCode == ccode {
				return z, nil
			}
		}
	}

	var match *Zone
	for _, z := range Zones {
		if (ccode == "" || z.CountryCode == ccode) && z.Zone == zone {
			return z, nil
		}
		if match == nil && z.Zone == zone {
			match = z
		}
	}

	if match != nil {
		return match, nil
	}
	return nil, fmt.Errorf("unknown timezone: %q %q", ccode, zone)
}

// MustNew behaves like New(), but will panic on errors.
func MustNew(ccode, zone string) *Zone {
	z, err := New(ccode, zone)
	if err != nil {
		panic(fmt.Errorf("tz.MustNew: %w", err))
	}
	return z
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

	// TODO: this could be aligned better with some spaces.
	if t.display == "" {
		var b strings.Builder
		b.WriteString(t.CountryName)
		b.WriteString(": ")
		b.WriteString(t.Zone)
		if len(t.Abbr) > 0 {
			b.WriteString(" (")
			b.WriteString(strings.Join(t.Abbr, ", "))
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

// OffsetRFC3339 gets the offset as a RFC3339 string: "+08:00", "-07:30", "UTC".
//
// Note that this displays the offset that is currently valid. For example
// Europe/Berlin may be +0100 or +0200, depending on whether DST is in effect.
func (t *Zone) OffsetRFC3339() string {
	o := t.Offset()
	if o == 0 {
		return "UTC"
	}
	m := o % 60
	return fmt.Sprintf("%+03.0f:%02d", math.Floor(float64(o)/60), m)
}

// OffsetDisplay gets the offset as a human readable string: "UTC +8:00", "UTC
// -7:30", "UTC".
//
// Note that this displays the offset that is currently valid. For example
// Europe/Berlin may be +0100 or +0200, depending on whether DST is in effect.
func (t *Zone) OffsetDisplay() string {
	o := t.Offset()
	if o == 0 {
		return "UTC"
	}
	m := o % 60
	return fmt.Sprintf("UTC %+.0f:%02d", math.Floor(float64(o)/60), m)
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
