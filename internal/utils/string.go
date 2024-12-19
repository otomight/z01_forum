package utils

import (
	"fmt"
	"strings"
	"unicode"
)

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

func CutStrAtPattern(str string, pattern string) (string, error) {
	var	index	int

	index = strings.Index(str, pattern)
	if index == -1 {
		return "", fmt.Errorf(
			"pattern '%s' not found in path '%s'", pattern, str,
		)
	}
	return str[index:], nil
}
