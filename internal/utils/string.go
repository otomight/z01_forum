package utils

import "unicode"

func IsStrEmpty(s string) bool {
	var	r	rune

	if len(s) == 0 {
		return true
	}
	for _, r = range s {
		if !unicode.IsPrint(r) || !unicode.IsSpace(r) {
			return false
		}
	}
	return true
}
