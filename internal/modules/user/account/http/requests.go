package http

import (
	"strings"

	"github.com/google/uuid"

	accountapp "github.com/Tokuchi61/Novascans/internal/modules/user/account/app"
	"github.com/Tokuchi61/Novascans/internal/modules/user/account/domain"
	"github.com/Tokuchi61/Novascans/internal/platform/validation"
)

type updateProfileRequest struct {
	Username    *string `json:"username"`
	DisplayName *string `json:"display_name"`
	Bio         *string `json:"bio"`
	AvatarPath  *string `json:"avatar_path"`
	BannerPath  *string `json:"banner_path"`
}

type updateSettingsRequest struct {
	Locale   *string `json:"locale"`
	Timezone *string `json:"timezone"`
}

type updatePrivacyRequest struct {
	ProfileVisibility string `json:"profile_visibility"`
}

func (request updateProfileRequest) Validate(_ *validation.Validator) validation.FieldErrors {
	errs := validation.FieldErrors{}
	if request.Username == nil && request.DisplayName == nil && request.Bio == nil && request.AvatarPath == nil && request.BannerPath == nil {
		errs.Add("profile", "at least one field is required")
		return errs
	}

	if request.Username != nil {
		username := domain.NormalizeUsername(*request.Username)
		if !domain.IsValidUsername(username) {
			errs.Add("username", "must be 3-32 chars and contain only lowercase letters, numbers, dots, underscores or hyphens")
		}
	}

	if request.DisplayName != nil && strings.TrimSpace(*request.DisplayName) == "" {
		errs.Add("display_name", "required")
	}

	return errs
}

func (request updateProfileRequest) ToInput(userID uuid.UUID) accountapp.UpdateProfileInput {
	return accountapp.UpdateProfileInput{
		UserID:      userID,
		Username:    request.Username,
		DisplayName: request.DisplayName,
		Bio:         request.Bio,
		AvatarPath:  request.AvatarPath,
		BannerPath:  request.BannerPath,
	}
}

func (request updateSettingsRequest) Validate(_ *validation.Validator) validation.FieldErrors {
	errs := validation.FieldErrors{}
	if request.Locale == nil && request.Timezone == nil {
		errs.Add("settings", "at least one field is required")
		return errs
	}

	if request.Locale != nil && strings.TrimSpace(*request.Locale) == "" {
		errs.Add("locale", "required")
	}

	if request.Timezone != nil && strings.TrimSpace(*request.Timezone) == "" {
		errs.Add("timezone", "required")
	}

	return errs
}

func (request updateSettingsRequest) ToInput(userID uuid.UUID) accountapp.UpdateSettingsInput {
	return accountapp.UpdateSettingsInput{
		UserID:   userID,
		Locale:   request.Locale,
		Timezone: request.Timezone,
	}
}

func (request updatePrivacyRequest) Validate(_ *validation.Validator) validation.FieldErrors {
	errs := validation.FieldErrors{}
	if !domain.IsValidProfileVisibility(strings.TrimSpace(request.ProfileVisibility)) {
		errs.Add("profile_visibility", "must be one of public, authenticated, private")
	}

	return errs
}

func (request updatePrivacyRequest) ToInput(userID uuid.UUID) accountapp.UpdatePrivacyInput {
	return accountapp.UpdatePrivacyInput{
		UserID:            userID,
		ProfileVisibility: strings.TrimSpace(request.ProfileVisibility),
	}
}
