/*
 * Copyright (c) 2024 Mohammadi El Youzghi and contributors.
 */

package util

import (
	"regexp"
	"strings"
)

func WildcardMatch(pattern string, value string) bool {
	pattern = regexp.QuoteMeta(pattern)

	pattern = strings.ReplaceAll(pattern, `\*`, `[^.]+`)

	pattern = strings.ReplaceAll(pattern, `#`, `.*`)

	pattern = "^" + pattern + "$"

	regex, err := regexp.Compile(pattern)
	if err != nil {
		return false
	}

	return regex.MatchString(value)
}
