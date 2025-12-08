package views

import "strings"

// Splits string by a specified delimiter (used mainly for parsing post media paths)
func SplitStringBy(s string, delim rune) []string {
	return strings.FieldsFunc(s, func(c rune) bool {
		return c == delim
	})
}
