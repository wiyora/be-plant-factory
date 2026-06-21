package config

import (
	"strings"
)

type StringList []string

func (s *StringList) UnmarshalText(text []byte) error {
	str := string(text)
	if str == "" {
		*s = []string{}
		return nil
	}

	parts := strings.Split(str, ",")
	result := make([]string, 0, len(parts))

	for _, part := range parts {
		item := strings.TrimSpace(part)
		if item != "" {
			result = append(result, item)
		}
	}

	*s = result
	return nil
}
