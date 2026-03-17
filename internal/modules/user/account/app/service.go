package app

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"

	authdomain "github.com/Tokuchi61/Novascans/internal/modules/identity/auth/domain"
	"github.com/Tokuchi61/Novascans/internal/modules/user/account/domain"
	platformdb "github.com/Tokuchi61/Novascans/internal/platform/db"
)

const (
	defaultLocale          = "en-US"
	defaultTimezone        = "UTC"
	maxUsernameLength      = 32
	minUsernameLength      = 3
	maxDisplayNameLength   = 64
	maxBioLength           = 280
	maxAssetPathLength     = 255
	maxLocaleLength        = 16
	maxTimezoneLength      = 64
	defaultProfileSuffixID = 8
)

var usernamePattern = regexp.MustCompile(`^[a-z0-9._-]{3,32}$`)

type Service struct {
	repo       Repository
	unitOfWork UnitOfWork
}

func NewService(repo Repository, unitOfWork UnitOfWork) *Service {
	return &Service{
		repo:       repo,
		unitOfWork: unitOfWork,
	}
}

func (service *Service) ProvisionDefaults(ctx context.Context, user authdomain.User, now time.Time) error {
	now = now.UTC()

	return service.withWriteRepository(ctx, func(ctx context.Context, repo Repository) error {
		if _, err := repo.GetProfileByUserID(ctx, user.ID); err != nil {
			if !HasCode(err, CodeNotFound) {
				return Internal("failed to load account profile", err)
			}

			username, err := service.uniqueDefaultUsername(ctx, repo, user)
			if err != nil {
				return err
			}

			if _, err := repo.CreateProfile(ctx, domain.Profile{
				UserID:      user.ID,
				Username:    username,
				DisplayName: defaultDisplayName(username),
				Bio:         "",
				AvatarPath:  "",
				BannerPath:  "",
				CreatedAt:   now,
				UpdatedAt:   now,
			}); err != nil {
				return Internal("failed to create account profile", err)
			}
		}

		if _, err := repo.GetSettingsByUserID(ctx, user.ID); err != nil {
			if !HasCode(err, CodeNotFound) {
				return Internal("failed to load account settings", err)
			}

			if _, err := repo.CreateSettings(ctx, domain.Settings{
				UserID:    user.ID,
				Locale:    defaultLocale,
				Timezone:  defaultTimezone,
				CreatedAt: now,
				UpdatedAt: now,
			}); err != nil {
				return Internal("failed to create account settings", err)
			}
		}

		if _, err := repo.GetPrivacySettingsByUserID(ctx, user.ID); err != nil {
			if !HasCode(err, CodeNotFound) {
				return Internal("failed to load account privacy settings", err)
			}

			if _, err := repo.CreatePrivacySettings(ctx, domain.PrivacySettings{
				UserID:            user.ID,
				ProfileVisibility: domain.ProfileVisibilityPublic,
				CreatedAt:         now,
				UpdatedAt:         now,
			}); err != nil {
				return Internal("failed to create account privacy settings", err)
			}
		}

		return nil
	})
}

func (service *Service) GetAccount(ctx context.Context, user authdomain.User) (AccountData, error) {
	profile, err := service.repo.GetProfileByUserID(ctx, user.ID)
	if err != nil {
		return AccountData{}, service.wrapNotFound("account profile not found", "failed to fetch account profile", err)
	}

	settings, err := service.repo.GetSettingsByUserID(ctx, user.ID)
	if err != nil {
		return AccountData{}, service.wrapNotFound("account settings not found", "failed to fetch account settings", err)
	}

	privacy, err := service.repo.GetPrivacySettingsByUserID(ctx, user.ID)
	if err != nil {
		return AccountData{}, service.wrapNotFound("account privacy settings not found", "failed to fetch account privacy settings", err)
	}

	return AccountData{
		User:     user,
		Profile:  profile,
		Settings: settings,
		Privacy:  privacy,
	}, nil
}

func (service *Service) GetProfile(ctx context.Context, userID uuid.UUID) (domain.Profile, error) {
	profile, err := service.repo.GetProfileByUserID(ctx, userID)
	if err != nil {
		return domain.Profile{}, service.wrapNotFound("account profile not found", "failed to fetch account profile", err)
	}

	return profile, nil
}

func (service *Service) UpdateProfile(ctx context.Context, input UpdateProfileInput) (domain.Profile, error) {
	var updated domain.Profile

	err := service.withWriteRepository(ctx, func(ctx context.Context, repo Repository) error {
		profile, err := repo.GetProfileByUserID(ctx, input.UserID)
		if err != nil {
			return service.wrapNotFound("account profile not found", "failed to fetch account profile", err)
		}

		if input.Username == nil && input.DisplayName == nil && input.Bio == nil && input.AvatarPath == nil && input.BannerPath == nil {
			return BadRequest("at least one profile field is required", nil)
		}

		if input.Username != nil {
			username := normalizeUsernameInput(*input.Username)
			if !isValidUsername(username) {
				return BadRequest("username must be 3-32 chars and contain only lowercase letters, numbers, dots, underscores or hyphens", nil)
			}

			if username != profile.Username {
				if _, err := repo.GetProfileByUsername(ctx, username); err == nil {
					return Conflict("username already exists", nil)
				} else if !HasCode(err, CodeNotFound) {
					return Internal("failed to check username uniqueness", err)
				}

				profile.Username = username
			}
		}

		if input.DisplayName != nil {
			displayName := strings.TrimSpace(*input.DisplayName)
			if displayName == "" {
				return BadRequest("display name is required", nil)
			}

			if len(displayName) > maxDisplayNameLength {
				return BadRequest(fmt.Sprintf("display name must be at most %d characters", maxDisplayNameLength), nil)
			}

			profile.DisplayName = displayName
		}

		if input.Bio != nil {
			bio := strings.TrimSpace(*input.Bio)
			if len(bio) > maxBioLength {
				return BadRequest(fmt.Sprintf("bio must be at most %d characters", maxBioLength), nil)
			}

			profile.Bio = bio
		}

		if input.AvatarPath != nil {
			avatarPath := strings.TrimSpace(*input.AvatarPath)
			if len(avatarPath) > maxAssetPathLength {
				return BadRequest(fmt.Sprintf("avatar path must be at most %d characters", maxAssetPathLength), nil)
			}

			profile.AvatarPath = avatarPath
		}

		if input.BannerPath != nil {
			bannerPath := strings.TrimSpace(*input.BannerPath)
			if len(bannerPath) > maxAssetPathLength {
				return BadRequest(fmt.Sprintf("banner path must be at most %d characters", maxAssetPathLength), nil)
			}

			profile.BannerPath = bannerPath
		}

		profile.UpdatedAt = time.Now().UTC()
		updated, err = repo.UpdateProfile(ctx, profile)
		if err != nil {
			return Internal("failed to update account profile", err)
		}

		return nil
	})
	if err != nil {
		return domain.Profile{}, err
	}

	return updated, nil
}

func (service *Service) GetSettings(ctx context.Context, userID uuid.UUID) (domain.Settings, error) {
	settings, err := service.repo.GetSettingsByUserID(ctx, userID)
	if err != nil {
		return domain.Settings{}, service.wrapNotFound("account settings not found", "failed to fetch account settings", err)
	}

	return settings, nil
}

func (service *Service) UpdateSettings(ctx context.Context, input UpdateSettingsInput) (domain.Settings, error) {
	var updated domain.Settings

	err := service.withWriteRepository(ctx, func(ctx context.Context, repo Repository) error {
		settings, err := repo.GetSettingsByUserID(ctx, input.UserID)
		if err != nil {
			return service.wrapNotFound("account settings not found", "failed to fetch account settings", err)
		}

		if input.Locale == nil && input.Timezone == nil {
			return BadRequest("at least one settings field is required", nil)
		}

		if input.Locale != nil {
			locale := strings.TrimSpace(*input.Locale)
			if locale == "" {
				return BadRequest("locale is required", nil)
			}

			if len(locale) > maxLocaleLength {
				return BadRequest(fmt.Sprintf("locale must be at most %d characters", maxLocaleLength), nil)
			}

			settings.Locale = locale
		}

		if input.Timezone != nil {
			timezone := strings.TrimSpace(*input.Timezone)
			if timezone == "" {
				return BadRequest("timezone is required", nil)
			}

			if len(timezone) > maxTimezoneLength {
				return BadRequest(fmt.Sprintf("timezone must be at most %d characters", maxTimezoneLength), nil)
			}

			settings.Timezone = timezone
		}

		settings.UpdatedAt = time.Now().UTC()
		updated, err = repo.UpdateSettings(ctx, settings)
		if err != nil {
			return Internal("failed to update account settings", err)
		}

		return nil
	})
	if err != nil {
		return domain.Settings{}, err
	}

	return updated, nil
}

func (service *Service) GetPrivacy(ctx context.Context, userID uuid.UUID) (domain.PrivacySettings, error) {
	privacy, err := service.repo.GetPrivacySettingsByUserID(ctx, userID)
	if err != nil {
		return domain.PrivacySettings{}, service.wrapNotFound("account privacy settings not found", "failed to fetch account privacy settings", err)
	}

	return privacy, nil
}

func (service *Service) UpdatePrivacy(ctx context.Context, input UpdatePrivacyInput) (domain.PrivacySettings, error) {
	if !domain.IsValidProfileVisibility(input.ProfileVisibility) {
		return domain.PrivacySettings{}, BadRequest("invalid profile visibility", nil)
	}

	var updated domain.PrivacySettings
	err := service.withWriteRepository(ctx, func(ctx context.Context, repo Repository) error {
		privacy, err := repo.GetPrivacySettingsByUserID(ctx, input.UserID)
		if err != nil {
			return service.wrapNotFound("account privacy settings not found", "failed to fetch account privacy settings", err)
		}

		privacy.ProfileVisibility = input.ProfileVisibility
		privacy.UpdatedAt = time.Now().UTC()

		updated, err = repo.UpdatePrivacySettings(ctx, privacy)
		if err != nil {
			return Internal("failed to update account privacy settings", err)
		}

		return nil
	})
	if err != nil {
		return domain.PrivacySettings{}, err
	}

	return updated, nil
}

func (service *Service) GetPublicProfile(ctx context.Context, username string, viewer Viewer) (domain.Profile, error) {
	username = normalizeUsernameInput(username)
	if !isValidUsername(username) {
		return domain.Profile{}, BadRequest("invalid username", nil)
	}

	profile, err := service.repo.GetProfileByUsername(ctx, username)
	if err != nil {
		return domain.Profile{}, service.wrapNotFound("profile not found", "failed to fetch account profile", err)
	}

	privacy, err := service.repo.GetPrivacySettingsByUserID(ctx, profile.UserID)
	if err != nil {
		return domain.Profile{}, service.wrapNotFound("profile not found", "failed to fetch account privacy settings", err)
	}

	if !canViewProfile(profile.UserID, privacy.ProfileVisibility, viewer) {
		return domain.Profile{}, NotFound("profile not found", nil)
	}

	return profile, nil
}

func (service *Service) withWriteRepository(ctx context.Context, fn func(ctx context.Context, repo Repository) error) error {
	if _, ok := platformdb.TxFromContext(ctx); ok {
		return fn(ctx, service.repo)
	}

	if service.unitOfWork == nil {
		return fn(ctx, service.repo)
	}

	return service.unitOfWork.WithinTransaction(ctx, fn)
}

func (service *Service) uniqueDefaultUsername(ctx context.Context, repo Repository, user authdomain.User) (string, error) {
	base := defaultUsernameBase(user.Email)

	if _, err := repo.GetProfileByUsername(ctx, base); err == nil {
		candidate := suffixUsername(base, user.ID.String()[:defaultProfileSuffixID])
		if _, secondErr := repo.GetProfileByUsername(ctx, candidate); secondErr == nil {
			return "", Conflict("generated username already exists", nil)
		} else if !HasCode(secondErr, CodeNotFound) {
			return "", Internal("failed to check generated username uniqueness", secondErr)
		}

		return candidate, nil
	} else if !HasCode(err, CodeNotFound) {
		return "", Internal("failed to check generated username uniqueness", err)
	}

	return base, nil
}

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
	if len(base) > maxUsernameLength {
		base = base[:maxUsernameLength]
	}

	if len(base) < minUsernameLength {
		base = "user"
	}

	if !isValidUsername(base) {
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
	suffix = strings.ToLower(strings.TrimSpace(suffix))
	if suffix == "" {
		return base
	}

	maxBase := maxUsernameLength - len(suffix) - 1
	if maxBase < minUsernameLength {
		maxBase = minUsernameLength
	}

	if len(base) > maxBase {
		base = base[:maxBase]
	}

	return strings.TrimRight(base, "._-") + "_" + suffix
}

func normalizeUsernameInput(value string) string {
	return strings.ToLower(strings.TrimSpace(value))
}

func isValidUsername(value string) bool {
	return usernamePattern.MatchString(value)
}

func canViewProfile(ownerID uuid.UUID, visibility string, viewer Viewer) bool {
	if viewer.UserID != nil && *viewer.UserID == ownerID {
		return true
	}

	switch visibility {
	case domain.ProfileVisibilityPublic:
		return true
	case domain.ProfileVisibilityAuthenticated:
		return viewer.Authenticated
	case domain.ProfileVisibilityPrivate:
		return false
	default:
		return false
	}
}

func (service *Service) wrapNotFound(notFoundMessage string, internalMessage string, err error) error {
	if HasCode(err, CodeNotFound) {
		return NotFound(notFoundMessage, err)
	}

	return Internal(internalMessage, err)
}
