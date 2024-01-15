package shared

import (
	"regexp"
	"strings"
)

var (
	removeStrings = []string{
		"( Your answer)",
		"( Missed)",
	}
)

// SanitizeText removes extra spaces, carriage returns, new lines, etc.
func SanitizeText(in string) (out string) {
	re_leadclose_whtsp := regexp.MustCompile(`^[\s\p{Zs}]+|[\s\p{Zs}]+$`)
	re_inside_whtsp := regexp.MustCompile(`[\s\p{Zs}]{2,}`)
	out = re_leadclose_whtsp.ReplaceAllString(in, "")
	out = re_inside_whtsp.ReplaceAllString(out, " ")

	for _, str := range removeStrings {
		out = strings.ReplaceAll(out, str, "")
	}

	return out
}
