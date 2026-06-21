package helper

import (
	"strings"
	"unicode"
)

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
