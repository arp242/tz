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

	// For zones that don't cange -c2020,2021 returns nothing, so append those
	// too.
	out, err = exec.Command("zdump", names...).Output()
	if err != nil {
		panic(err)
	}

	for _, line := range strings.Split(string(out), "\n") {
		f := strings.Fields(line)
		if len(f) > 5 && f[6][0] != '-' && f[6][0] != '+' {
			r[f[0]] = append(r[f[0]], f[6])
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
