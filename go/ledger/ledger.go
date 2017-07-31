package ledger

import (
	"errors"
	"strconv"
	"strings"
)

const testVersion = 4

type Entry struct {
	Date        string // "Y-m-d"
	Description string
	Change      int // in cents
}

type result struct {
	i int
	s string
	e error
}

func FormatLedger(currency string, locale string, entries []Entry) (string, error) {
	if len(entries) == 0 {
		if _, err := FormatLedger(currency, "en-US", []Entry{{Date: "2014-01-01", Description: "", Change: 0}}); err != nil {
			return "", err
		}
	}

	var entriesCopy []Entry
	for _, e := range entries {
		entriesCopy = append(entriesCopy, e)
	}

	entriesCopy = sortEntries(entriesCopy)

	s, err := fmtHeader(locale)
	if err != nil {
		return s, err
	}

	// Parallelism, always a great idea
	co := make(chan result)
	for i, et := range entriesCopy {
		go fmtEntries(i, currency, locale, et, co)
	}

	ss := make([]string, len(entriesCopy))
	for range entriesCopy {
		v := <-co
		if v.e != nil {
			return "", v.e
		}
		ss[v.i] = v.s
	}
	for i := 0; i < len(entriesCopy); i++ {
		s += ss[i]
	}
	return s, nil
}

func sortEntries(entriesCopy []Entry) []Entry {
	m1 := map[bool]int{true: 0, false: 1}
	m2 := map[bool]int{true: -1, false: 1}
	es := entriesCopy
	for len(es) > 1 {
		first, rest := es[0], es[1:]
		success := false
		for !success {
			success = true
			for i, e := range rest {
				if (m1[e.Date == first.Date]*m2[e.Date < first.Date]*4 +
					m1[e.Description == first.Description]*m2[e.Description < first.Description]*2 +
					m1[e.Change == first.Change]*m2[e.Change < first.Change]*1) < 0 {
					es[0], es[i+1] = es[i+1], es[0]
					success = false
				}
			}
		}
		es = es[1:]
	}
	return entriesCopy
}

func fmtHeader(locale string) (string, error) {
	var s string
	switch locale {
	case "nl-NL":
		s = "Datum" +
			strings.Repeat(" ", 10-len("Datum")) +
			" | " +
			"Omschrijving" +
			strings.Repeat(" ", 25-len("Omschrijving")) +
			" | " + "Verandering" + "\n"
	case "en-US":
		s = "Date" +
			strings.Repeat(" ", 10-len("Date")) +
			" | " +
			"Description" +
			strings.Repeat(" ", 25-len("Description")) +
			" | " + "Change" + "\n"
	default:
		return "", errors.New("")
	}
	return s, nil
}

func fmtEntries(i int, currency string, locale string, entry Entry, co chan result) {
	var dt string
	dt = fmtDate(locale, entry.Date, co)

	var desc string
	desc = fmtDesc(entry)

	cents := entry.Change
	var amt string
	amt = fmtAmt(currency, locale, cents, co)

	var al int
	for range amt {
		al++
	}
	co <- result{i: i, s: dt + strings.Repeat(" ", 10-len(dt)) + " | " + desc + " | " +
		strings.Repeat(" ", 13-al) + amt + "\n"}
}

func fmtDate(locale string, dt string, co chan result) string {
	if len(dt) != 10 {
		co <- result{e: errors.New("")}
	}
	d1, d2, d3, d4, d5 := dt[0:4], dt[4], dt[5:7], dt[7], dt[8:10]
	if d2 != '-' {
		co <- result{e: errors.New("")}
	}
	if d4 != '-' {
		co <- result{e: errors.New("")}
	}
	var d string
	if locale == "nl-NL" {
		d = d5 + "-" + d3 + "-" + d1
	} else if locale == "en-US" {
		d = d3 + "/" + d5 + "/" + d1
	}
	return d
}

func fmtDesc(entry Entry) string {
	de := entry.Description
	if len(de) > 25 {
		de = de[:22] + "..."
	} else {
		de = de + strings.Repeat(" ", 25-len(de))
	}
	return de
}

func fmtAmt(currency string, locale string, cents int, co chan result) string {
	negative := false
	if cents < 0 {
		cents = cents * -1
		negative = true
	}

	var a string
	switch locale {
	case "nl-NL":
		a = fmtAmtNl(currency, cents, negative, co)
	case "en-US":
		a = fmtAmtEn(currency, cents, negative, co)
	default:
		co <- result{e: errors.New("")}
	}
	return a
}

func fmtAmtEn(currency string, cents int, negative bool, co chan result) string {
	var a string
	if negative {
		a += "("
	}

	switch currency {
	case "EUR":
		a += "€"
	case "USD":
		a += "$"
	default:
		co <- result{e: errors.New("")}
	}

	centsStr := strconv.Itoa(cents)
	switch len(centsStr) {
	case 1:
		centsStr = "00" + centsStr
	case 2:
		centsStr = "0" + centsStr
	}

	rest := centsStr[:len(centsStr)-2]
	var parts []string
	for len(rest) > 3 {
		parts = append(parts, rest[len(rest)-3:])
		rest = rest[:len(rest)-3]
	}
	if len(rest) > 0 {
		parts = append(parts, rest)
	}
	for i := len(parts) - 1; i >= 0; i-- {
		a += parts[i] + ","
	}
	a = a[:len(a)-1]
	a += "."
	a += centsStr[len(centsStr)-2:]
	if negative {
		a += ")"
	} else {
		a += " "
	}
	return a
}

func fmtAmtNl(currency string, cents int, negative bool, co chan result) string {
	var a string
	switch currency {
	case "EUR":
		a += "€"
	case "USD":
		a += "$"
	default:
		co <- result{e: errors.New("")}
	}
	a += " "
	centsStr := strconv.Itoa(cents)
	switch len(centsStr) {
	case 1:
		centsStr = "00" + centsStr
	case 2:

		centsStr = "0" + centsStr
	}
	rest := centsStr[:len(centsStr)-2]
	var parts []string
	for len(rest) > 3 {
		parts = append(parts, rest[len(rest)-3:])
		rest = rest[:len(rest)-3]
	}
	if len(rest) > 0 {
		parts = append(parts, rest)
	}
	for i := len(parts) - 1; i >= 0; i-- {
		a += parts[i] + "."
	}
	a = a[:len(a)-1]
	a += ","
	a += centsStr[len(centsStr)-2:]
	if negative {
		a += "-"
	} else {
		a += " "
	}
	return a
}
