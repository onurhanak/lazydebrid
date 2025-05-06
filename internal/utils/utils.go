package utils

import (
	"strings"
)

func Match(filename, query string) bool {
	return strings.Contains(strings.ToLower(filename), strings.ToLower(query))
}
