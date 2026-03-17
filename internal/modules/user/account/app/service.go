package app

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	authdomain "github.com/Tokuchi61/Novascans/internal/modules/identity/auth/domain"
	"github.com/Tokuchi61/Novascans/internal/modules/user/account/domain"
	platformdb "github.com/Tokuchi61/Novascans/internal/platform/db"
)

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

func (service *Service) GetAccount(ctx context.Context, userID uuid.UUID) (AccountData, error) {
	profile, err := service.repo.GetProfileByUserID(ctx, userID)
	if err != nil {
		return AccountData{}, service.wrapNotFound("account profile not found", "failed to fetch account profile", err)
	}

	settings, err := service.repo.GetSettingsByUserID(ctx, userID)
	if err != nil {
		return AccountData{}, service.wrapNotFound("account settings not found", "failed to fetch account settings", err)
	}

	privacy, err := service.repo.GetPrivacySettingsByUserID(ctx, userID)
	if err != nil {
		return AccountData{}, service.wrapNotFound("account privacy settings not found", "failed to fetch account privacy settings", err)
	}

	return AccountData{
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
			username := domain.NormalizeUsername(*input.Username)
			if !domain.IsValidUsername(username) {
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

			if len(displayName) > domain.DisplayNameMaxLength {
				return BadRequest(fmt.Sprintf("display name must be at most %d characters", domain.DisplayNameMaxLength), nil)
			}

			profile.DisplayName = displayName
		}

		if input.Bio != nil {
			bio := strings.TrimSpace(*input.Bio)
			if len(bio) > domain.BioMaxLength {
				return BadRequest(fmt.Sprintf("bio must be at most %d characters", domain.BioMaxLength), nil)
			}

			profile.Bio = bio
		}

		if input.AvatarPath != nil {
			avatarPath := strings.TrimSpace(*input.AvatarPath)
			if len(avatarPath) > domain.AssetPathMaxLength {
				return BadRequest(fmt.Sprintf("avatar path must be at most %d characters", domain.AssetPathMaxLength), nil)
			}

			profile.AvatarPath = avatarPath
		}

		if input.BannerPath != nil {
			bannerPath := strings.TrimSpace(*input.BannerPath)
			if len(bannerPath) > domain.AssetPathMaxLength {
				return BadRequest(fmt.Sprintf("banner path must be at most %d characters", domain.AssetPathMaxLength), nil)
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

			if len(locale) > domain.LocaleMaxLength {
				return BadRequest(fmt.Sprintf("locale must be at most %d characters", domain.LocaleMaxLength), nil)
			}

			settings.Locale = locale
		}

		if input.Timezone != nil {
			timezone := strings.TrimSpace(*input.Timezone)
			if timezone == "" {
				return BadRequest("timezone is required", nil)
			}

			if len(timezone) > domain.TimezoneMaxLength {
				return BadRequest(fmt.Sprintf("timezone must be at most %d characters", domain.TimezoneMaxLength), nil)
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
	username = domain.NormalizeUsername(username)
	if !domain.IsValidUsername(username) {
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

	if !domain.CanViewProfile(profile.UserID, privacy.ProfileVisibility, viewer.UserID, viewer.Authenticated) {
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
	base := defaultUsernameForUser(user)

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

func (service *Service) wrapNotFound(notFoundMessage string, internalMessage string, err error) error {
	if HasCode(err, CodeNotFound) {
		return NotFound(notFoundMessage, err)
	}

	return Internal(internalMessage, err)
}
