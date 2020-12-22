// +build go_run_only

package main

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"sort"
	"strings"
	"time"
)

type Zone struct {
	CountryCode string   // ID
	Zone        string   // Asia/Makassar
	Abbr        []string // WITA
	CountryName string   // Indonesia
	Comments    string   // Borneo (east, south); Sulawesi/Celebes, Bali, Nusa Tengarra; Timor (west)
}

func readISO() map[string]string {
	f, err := ioutil.ReadFile("/usr/share/zoneinfo/iso3166.tab")
	if err != nil {
		panic(err)
	}

	r := make(map[string]string)
	for _, line := range strings.Split(string(f), "\n") {
		if p := strings.Index(line, "#"); p > -1 {
			line = line[:p]
		}
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		s := strings.Split(line, "\t")
		r[s[0]] = s[1]
	}

	return r
}

// TODO: this is wrong, as it may return "CET" or "CEST" depending on DST:
//
// [~]% zdump -v -c 2019,2020 Europe/Berlin Pacific/Auckland
// Europe/Berlin     -9223372036854775808 = NULL
// Europe/Berlin     -9223372036854689408 = NULL
// Europe/Berlin     Sun Mar 31 00:59:59 2019 UT = Sun Mar 31 01:59:59 2019 CET isdst=0 gmtoff=3600
// Europe/Berlin     Sun Mar 31 01:00:00 2019 UT = Sun Mar 31 03:00:00 2019 CEST isdst=1 gmtoff=7200
// Europe/Berlin     Sun Oct 27 00:59:59 2019 UT = Sun Oct 27 02:59:59 2019 CEST isdst=1 gmtoff=7200
// Europe/Berlin     Sun Oct 27 01:00:00 2019 UT = Sun Oct 27 02:00:00 2019 CET isdst=0 gmtoff=3600
// Europe/Berlin     9223372036854689407 = NULL
// Europe/Berlin     9223372036854775807 = NULL
// Pacific/Auckland  -9223372036854775808 = NULL
// Pacific/Auckland  -9223372036854689408 = NULL
// Pacific/Auckland  Sat Apr  6 13:59:59 2019 UT = Sun Apr  7 02:59:59 2019 NZDT isdst=1 gmtoff=46800
// Pacific/Auckland  Sat Apr  6 14:00:00 2019 UT = Sun Apr  7 02:00:00 2019 NZST isdst=0 gmtoff=43200
// Pacific/Auckland  Sat Sep 28 13:59:59 2019 UT = Sun Sep 29 01:59:59 2019 NZST isdst=0 gmtoff=43200
// Pacific/Auckland  Sat Sep 28 14:00:00 2019 UT = Sun Sep 29 03:00:00 2019 NZDT isdst=1 gmtoff=46800
// Pacific/Auckland  9223372036854689407 = NULL
// Pacific/Auckland  9223372036854775807 = NULL
//
// [~]% zdump -v -c 2019,2020 Europe/Berlin Pacific/Auckland | grep gmtoff= | awk '{print $14 " " $15}' | sort -u
// CEST isdst=1
// CET isdst=0
// NZDT isdst=1
// NZST isdst=0
func readAbbr(names []string) map[string][]string {
	f := fmt.Sprintf("-Vc%s,%s", time.Now().UTC().Format("2006"), time.Now().Add(365*time.Hour*24).UTC().Format("2006"))
	out, err := exec.Command("zdump", append([]string{f}, names...)...).Output()
	if err != nil {
		panic(err)
	}

	r := make(map[string][]string)
	for _, line := range strings.Split(string(out), "\n") {
		f := strings.Fields(line)
		if len(f) < 13 {
			continue
		}

		abbr := f[13]
		if abbr[0] != '-' && abbr[0] != '+' {
			r[f[0]] = append(r[f[0]], abbr)
		}
	}

	for k := range r {
		r[k] = Uniq(r[k])
	}
	return r
}

// Uniq removes duplicate entries from list; the list will be sorted.
func Uniq(list []string) []string {
	sort.Strings(list)
	var last string
	l := list[:0]
	for _, str := range list {
		if str != last {
			l = append(l, str)
		}
		last = str
	}
	return l
}

func main() {
	iso := readISO()

	f, err := ioutil.ReadFile("/usr/share/zoneinfo/zone1970.tab")
	if err != nil {
		panic(err)
	}

	var (
		r     []Zone
		names []string
	)
	for _, line := range strings.Split(string(f), "\n") {
		if p := strings.Index(line, "#"); p > -1 {
			line = line[:p]
		}
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// #codes	coordinates	TZ	comments
		s := strings.Split(line, "\t")
		countries := strings.Split(s[0], ",")
		desc := ""
		if len(s) > 3 {
			desc = s[3]
		}

		names = append(names, s[2])
		for _, country := range countries {
			r = append(r, Zone{
				CountryCode: country,
				CountryName: iso[country],
				Zone:        s[2],
				Comments:    desc,
			})
		}
	}
	sort.Slice(r, func(i, j int) bool {
		return r[i].CountryName < r[j].CountryName
	})
	// Move Åland Islands to proper location; don't want to bother with unicode
	// collate just for this.
	last := r[len(r)-1]
	r = append([]Zone{r[0]}, append([]Zone{last}, r[1:len(r)-1]...)...)

	abbr := readAbbr(names)
	for i := range r {
		if a := abbr[r[i].Zone]; len(a) > 0 {
			r[i].Abbr = a
		}
	}

	fmt.Print("package tz\n\n")
	fmt.Println("// Zones is a list of all timezones by country.")
	fmt.Println("var Zones = []*Zone{")
	for i := range r {
		l := fmt.Sprintf("%#v,\n", r[i])
		fmt.Print("\t" + l[9:])
	}
	fmt.Println("}")
}
