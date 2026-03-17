package domain

import (
	"regexp"
	"strings"
)

const (
	UsernameMinLength    = 3
	UsernameMaxLength    = 32
	DisplayNameMaxLength = 64
	BioMaxLength         = 280
	AssetPathMaxLength   = 255
)

var usernamePattern = regexp.MustCompile(`^[a-z0-9._-]{3,32}$`)

func NormalizeUsername(value string) string {
	return strings.ToLower(strings.TrimSpace(value))
}

func IsValidUsername(value string) bool {
	return usernamePattern.MatchString(value)
}
