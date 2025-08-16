// Package utils provides utility functions that can be used by external applications
package utils

import "strings"

// StringHelper provides string manipulation utilities
type StringHelper struct{}

// NewStringHelper creates a new StringHelper instance
func NewStringHelper() *StringHelper {
	return &StringHelper{}
}

// Capitalize returns a string with the first letter capitalized
func (s *StringHelper) Capitalize(str string) string {
	if len(str) == 0 {
		return str
	}
	return strings.ToUpper(string(str[0])) + strings.ToLower(str[1:])
}