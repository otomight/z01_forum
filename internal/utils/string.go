package utils

import "unicode"

func IsOnlyPrintable(s string) bool {
	var	r	rune

	if len(s) == 0 {
		return false
	}
	for _, r = range s {
		if !unicode.IsPrint(r) {
			return false
		}
	}
	return true
}
