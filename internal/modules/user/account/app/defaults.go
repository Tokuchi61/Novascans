package app

import (
	"strings"

	authdomain "github.com/Tokuchi61/Novascans/internal/modules/identity/auth/domain"
	"github.com/Tokuchi61/Novascans/internal/modules/user/account/domain"
)

const (
	defaultLocale          = "en-US"
	defaultTimezone        = "UTC"
	defaultProfileSuffixID = 8
)

func defaultUsernameBase(email string) string {
	local := email
	if parts := strings.SplitN(email, "@", 2); len(parts) > 0 {
		local = parts[0]
	}

	local = strings.ToLower(strings.TrimSpace(local))
	builder := strings.Builder{}
	for _, r := range local {
		switch {
		case r >= 'a' && r <= 'z':
			builder.WriteRune(r)
		case r >= '0' && r <= '9':
			builder.WriteRune(r)
		case r == '.' || r == '_' || r == '-':
			builder.WriteRune(r)
		default:
			builder.WriteRune('_')
		}
	}

	base := strings.Trim(builder.String(), "._-")
	if len(base) > domain.UsernameMaxLength {
		base = base[:domain.UsernameMaxLength]
	}

	if len(base) < domain.UsernameMinLength {
		base = "user"
	}

	if !domain.IsValidUsername(base) {
		base = "user"
	}

	return base
}

func defaultDisplayName(username string) string {
	if username == "" {
		return "User"
	}

	return username
}

func suffixUsername(base string, suffix string) string {
	suffix = domain.NormalizeUsername(suffix)
	if suffix == "" {
		return base
	}

	maxBase := domain.UsernameMaxLength - len(suffix) - 1
	if maxBase < domain.UsernameMinLength {
		maxBase = domain.UsernameMinLength
	}

	if len(base) > maxBase {
		base = base[:maxBase]
	}

	return strings.TrimRight(base, "._-") + "_" + suffix
}

func defaultUsernameForUser(user authdomain.User) string {
	return defaultUsernameBase(user.Email)
}
