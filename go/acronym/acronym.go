package acronym

import "strings"

const testVersion = 2

func Abbreviate(fullname string) string {
	var acronym string

	var addNextChar bool = true
	for _, ch := range fullname {
		if ch >= 65 && ch <= 90 {
			acronym = acronym + string(ch)
			addNextChar = false
		} else if addNextChar {
			acronym = acronym + strings.ToUpper(string(ch))
			addNextChar = false
		} else if ch == ' ' || ch == '-' {
			addNextChar = true
		} else if ch == ':' {
			break
		}
	}
	return acronym
}
