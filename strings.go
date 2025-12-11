package main

import (
	"regexp"
	"strings"
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

var (
	reASCII       = regexp.MustCompile(`^[ -~]*$`)
	reIllegalName = regexp.MustCompile(`[^[:alnum:]-. ]`)
)

func isASCII(s string) bool {
	return reASCII.MatchString(s)
}

func cleanBasename(s string) string {
	return reIllegalName.ReplaceAllString(s, "")
}
