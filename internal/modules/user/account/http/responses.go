package http

import (
	"time"

	authdomain "github.com/Tokuchi61/Novascans/internal/modules/identity/auth/domain"
	accountapp "github.com/Tokuchi61/Novascans/internal/modules/user/account/app"
	"github.com/Tokuchi61/Novascans/internal/modules/user/account/domain"
)

type accountResponse struct {
	User     accountUserResponse `json:"user"`
	Profile  profileResponse     `json:"profile"`
	Settings settingsResponse    `json:"settings"`
	Privacy  privacyResponse     `json:"privacy"`
}

type accountUserResponse struct {
	ID              string `json:"id"`
	Email           string `json:"email"`
	BaseRole        string `json:"base_role"`
	Status          string `json:"status"`
	EmailVerifiedAt string `json:"email_verified_at,omitempty"`
}

type profileResponse struct {
	UserID      string `json:"user_id"`
	Username    string `json:"username"`
	DisplayName string `json:"display_name"`
	Bio         string `json:"bio"`
	AvatarPath  string `json:"avatar_path,omitempty"`
	BannerPath  string `json:"banner_path,omitempty"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type settingsResponse struct {
	UserID    string `json:"user_id"`
	Locale    string `json:"locale"`
	Timezone  string `json:"timezone"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type privacyResponse struct {
	UserID            string `json:"user_id"`
	ProfileVisibility string `json:"profile_visibility"`
	CreatedAt         string `json:"created_at"`
	UpdatedAt         string `json:"updated_at"`
}

func mapAccountResponse(account accountapp.AccountData) accountResponse {
	return accountResponse{
		User:     mapAccountUserResponse(account.User),
		Profile:  mapProfileResponse(account.Profile),
		Settings: mapSettingsResponse(account.Settings),
		Privacy:  mapPrivacyResponse(account.Privacy),
	}
}

func mapAccountUserResponse(user authdomain.User) accountUserResponse {
	response := accountUserResponse{
		ID:       user.ID.String(),
		Email:    user.Email,
		BaseRole: user.BaseRole,
		Status:   user.Status,
	}

	if user.EmailVerifiedAt != nil {
		response.EmailVerifiedAt = user.EmailVerifiedAt.UTC().Format(time.RFC3339)
	}

	return response
}

func mapProfileResponse(profile domain.Profile) profileResponse {
	return profileResponse{
		UserID:      profile.UserID.String(),
		Username:    profile.Username,
		DisplayName: profile.DisplayName,
		Bio:         profile.Bio,
		AvatarPath:  profile.AvatarPath,
		BannerPath:  profile.BannerPath,
		CreatedAt:   profile.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt:   profile.UpdatedAt.UTC().Format(time.RFC3339),
	}
}

func mapSettingsResponse(settings domain.Settings) settingsResponse {
	return settingsResponse{
		UserID:    settings.UserID.String(),
		Locale:    settings.Locale,
		Timezone:  settings.Timezone,
		CreatedAt: settings.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt: settings.UpdatedAt.UTC().Format(time.RFC3339),
	}
}

func mapPrivacyResponse(privacy domain.PrivacySettings) privacyResponse {
	return privacyResponse{
		UserID:            privacy.UserID.String(),
		ProfileVisibility: privacy.ProfileVisibility,
		CreatedAt:         privacy.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt:         privacy.UpdatedAt.UTC().Format(time.RFC3339),
	}
}
