package utils

import "strings"

func Unhexlify(input string) string {
	if strings.HasPrefix(input, "0x") {
		return strings.TrimPrefix(input, "0x")
	}
	return input
}
