package isogram

import "strings"

const testVersion = 1

func IsIsogram(word string) bool {
	chmap := make(map[int]bool)
	lowWord := strings.ToLower(word)

	for _, ch := range lowWord {
		ascii := int(ch)
		if ch != '-' && ch != ' ' {
			if chmap[ascii] {
				return false
			}
			chmap[ascii] = true
		}
	}
	return true
}
