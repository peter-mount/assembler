package util

import (
	"strconv"
	"strings"
)

func Atoi(s string) (int64, error) {
	switch {
	case strings.HasPrefix(s, "0x") && len(s) > 2:
		return strconv.ParseInt(s[2:], 16, 64)

	case (strings.HasPrefix(s, "&") || strings.HasPrefix(s, "$")) && len(s) > 1:
		return strconv.ParseInt(s[1:], 16, 64)

	case strings.HasPrefix(s, "0") && len(s) > 1:
		return strconv.ParseInt(s[1:], 8, 64)

	default:
		return strconv.ParseInt(s, 10, 64)
	}
}

func IsHex(s string) bool {
	_, err := strconv.ParseInt(s, 16, 64)
	return err == nil
}
