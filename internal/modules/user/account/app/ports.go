package app

import (
	"context"

	"github.com/google/uuid"

	"github.com/Tokuchi61/Novascans/internal/modules/user/account/domain"
)

type Repository interface {
	CreateProfile(ctx context.Context, profile domain.Profile) (domain.Profile, error)
	GetProfileByUserID(ctx context.Context, userID uuid.UUID) (domain.Profile, error)
	GetProfileByUsername(ctx context.Context, username string) (domain.Profile, error)
	UpdateProfile(ctx context.Context, profile domain.Profile) (domain.Profile, error)
	CreateSettings(ctx context.Context, settings domain.Settings) (domain.Settings, error)
	GetSettingsByUserID(ctx context.Context, userID uuid.UUID) (domain.Settings, error)
	UpdateSettings(ctx context.Context, settings domain.Settings) (domain.Settings, error)
	CreatePrivacySettings(ctx context.Context, settings domain.PrivacySettings) (domain.PrivacySettings, error)
	GetPrivacySettingsByUserID(ctx context.Context, userID uuid.UUID) (domain.PrivacySettings, error)
	UpdatePrivacySettings(ctx context.Context, settings domain.PrivacySettings) (domain.PrivacySettings, error)
}

type UnitOfWork interface {
	WithinTransaction(ctx context.Context, fn func(ctx context.Context, repo Repository) error) error
}
