package helper

import (
	"strings"
	"unicode"
)

func TruncateString(str string, maxLength int, suffix string) string {
	if len(str) <= maxLength {
		return str
	}

	if maxLength <= len(suffix) {
		return suffix
	}

	return str[:maxLength-len(suffix)] + suffix
}

func ToSnakeCase(s string) string {
	var b strings.Builder
	b.Grow(len(s) + 2)

	for i, r := range s {
		if unicode.IsUpper(r) {
			if i > 0 {
				b.WriteByte('_')
			}

			b.WriteRune(unicode.ToLower(r))
			continue
		}

		b.WriteRune(r)
	}

	return b.String()
}
