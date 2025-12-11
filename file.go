package main

import (
	"errors"
	"os"
)

func fileExists(path string) bool {
	if _, err := os.Stat(path); !errors.Is(err, os.ErrNotExist) {
		return true
	}
	return false
}
