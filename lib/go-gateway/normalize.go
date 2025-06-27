package gateway

import "strings"

// NormalizeText trim space and validate utf8
func NormalizeText(s string) string {
	s = strings.TrimSpace(s)
	return strings.ToValidUTF8(s, "")
}
