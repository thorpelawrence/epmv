package main

import (
	"bytes"
	"strings"
	"unicode"
)

const (
	suffixEPUB = ".epub"
)

func trimSuffixEPUB(s string) string {
	if len(s) >= len(suffixEPUB) && strings.EqualFold(s[len(s)-len(suffixEPUB):], suffixEPUB) {
		return s[:len(s)-len(suffixEPUB)]
	}

	return s
}

func trimNul(s string) string {
	return strings.TrimSpace(string(bytes.Trim([]byte(s), "\x00")))
}

func isASCII(s string) bool {
	for _, c := range s {
		if c > unicode.MaxASCII {
			return false
		}
	}

	return true
}
