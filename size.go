// Copyright 2014 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package utils

import (
	"math"
	"strconv"
	"strings"
	"unicode"

	"github.com/juju/errors"
)

// ParseSize parses the string as a size, in mebibytes.
//
// The string must be a is a non-negative number with
// an optional multiplier suffix (M, G, T, P, E, Z, or Y).
// If the suffix is not specified, "M" is implied.
func ParseSize(str string) (MB uint64, err error) {
	// Find the first non-digit/period:
	i := strings.IndexFunc(str, func(r rune) bool {
		return r != '.' && !unicode.IsDigit(r)
	})
	var multiplier float64 = 1
	if i > 0 {
		suffix := str[i:]
		multiplier = 0
		for j := 0; j < len(sizeSuffixes); j++ {
			base := string(sizeSuffixes[j])
			// M, MB, or MiB are all valid.
			switch suffix {
			case base, base + "B", base + "iB":
				multiplier = float64(sizeSuffixMultiplier(j))
				break
			}
		}
		if multiplier == 0 {
			return 0, errors.Errorf("invalid multiplier suffix %q, expected one of %s", suffix, []byte(sizeSuffixes))
		}
		str = str[:i]
	}

	val, err := strconv.ParseFloat(str, 64)
	if err != nil || val < 0 {
		return 0, errors.Errorf("expected a non-negative number, got %q", str)
	}
	val *= multiplier
	return uint64(math.Ceil(val)), nil
}

var sizeSuffixes = "MGTPEZY"

func sizeSuffixMultiplier(i int) int {
	return 1 << uint(i*10)
}
