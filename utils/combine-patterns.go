// Package utils
package utils

import (
	"regexp"
	"strings"
)

func CombinePatterns[T ~[]string](patterns T) regexp.Regexp {
	pattern := strings.Join(patterns, "|")
	return *regexp.MustCompile(pattern)
}
