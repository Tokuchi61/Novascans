package store

import (
	"context"
	"database/sql"

	"github.com/google/uuid"

	accountsqlc "github.com/Tokuchi61/Novascans/internal/gen/sqlc/user/account"
	accountapp "github.com/Tokuchi61/Novascans/internal/modules/user/account/app"
	"github.com/Tokuchi61/Novascans/internal/modules/user/account/domain"
	platformdb "github.com/Tokuchi61/Novascans/internal/platform/db"
)

type PostgresRepository struct {
	db *sql.DB
}

type PostgresUnitOfWork struct {
	repo      *PostgresRepository
	txManager platformdb.TxManager
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func NewPostgresUnitOfWork(repo *PostgresRepository, txManager platformdb.TxManager) *PostgresUnitOfWork {
	return &PostgresUnitOfWork{
		repo:      repo,
		txManager: txManager,
	}
}

func (repo *PostgresRepository) CreateProfile(ctx context.Context, profile domain.Profile) (domain.Profile, error) {
	row, err := repo.queries(ctx).CreateProfile(ctx, accountsqlc.CreateProfileParams{
		UserID:      profile.UserID,
		Username:    profile.Username,
		DisplayName: profile.DisplayName,
		Bio:         profile.Bio,
		AvatarPath:  profile.AvatarPath,
		BannerPath:  profile.BannerPath,
		CreatedAt:   profile.CreatedAt,
		UpdatedAt:   profile.UpdatedAt,
	})
	if err != nil {
		return domain.Profile{}, err
	}

	return mapProfile(row), nil
}

func (repo *PostgresRepository) GetProfileByUserID(ctx context.Context, userID uuid.UUID) (domain.Profile, error) {
	row, err := repo.queries(ctx).GetProfileByUserID(ctx, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.Profile{}, accountapp.NotFound("account profile not found", err)
		}

		return domain.Profile{}, err
	}

	return mapProfile(row), nil
}

func (repo *PostgresRepository) GetProfileByUsername(ctx context.Context, username string) (domain.Profile, error) {
	row, err := repo.queries(ctx).GetProfileByUsername(ctx, username)
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.Profile{}, accountapp.NotFound("account profile not found", err)
		}

		return domain.Profile{}, err
	}

	return mapProfile(row), nil
}

func (repo *PostgresRepository) UpdateProfile(ctx context.Context, profile domain.Profile) (domain.Profile, error) {
	row, err := repo.queries(ctx).UpdateProfile(ctx, accountsqlc.UpdateProfileParams{
		UserID:      profile.UserID,
		Username:    profile.Username,
		DisplayName: profile.DisplayName,
		Bio:         profile.Bio,
		AvatarPath:  profile.AvatarPath,
		BannerPath:  profile.BannerPath,
		UpdatedAt:   profile.UpdatedAt,
	})
	if err != nil {
		return domain.Profile{}, err
	}

	return mapProfile(row), nil
}

func (repo *PostgresRepository) CreateSettings(ctx context.Context, settings domain.Settings) (domain.Settings, error) {
	row, err := repo.queries(ctx).CreateSettings(ctx, accountsqlc.CreateSettingsParams{
		UserID:    settings.UserID,
		Locale:    settings.Locale,
		Timezone:  settings.Timezone,
		CreatedAt: settings.CreatedAt,
		UpdatedAt: settings.UpdatedAt,
	})
	if err != nil {
		return domain.Settings{}, err
	}

	return mapSettings(row), nil
}

func (repo *PostgresRepository) GetSettingsByUserID(ctx context.Context, userID uuid.UUID) (domain.Settings, error) {
	row, err := repo.queries(ctx).GetSettingsByUserID(ctx, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.Settings{}, accountapp.NotFound("account settings not found", err)
		}

		return domain.Settings{}, err
	}

	return mapSettings(row), nil
}

func (repo *PostgresRepository) UpdateSettings(ctx context.Context, settings domain.Settings) (domain.Settings, error) {
	row, err := repo.queries(ctx).UpdateSettings(ctx, accountsqlc.UpdateSettingsParams{
		UserID:    settings.UserID,
		Locale:    settings.Locale,
		Timezone:  settings.Timezone,
		UpdatedAt: settings.UpdatedAt,
	})
	if err != nil {
		return domain.Settings{}, err
	}

	return mapSettings(row), nil
}

func (repo *PostgresRepository) CreatePrivacySettings(ctx context.Context, settings domain.PrivacySettings) (domain.PrivacySettings, error) {
	row, err := repo.queries(ctx).CreatePrivacySettings(ctx, accountsqlc.CreatePrivacySettingsParams{
		UserID:            settings.UserID,
		ProfileVisibility: settings.ProfileVisibility,
		CreatedAt:         settings.CreatedAt,
		UpdatedAt:         settings.UpdatedAt,
	})
	if err != nil {
		return domain.PrivacySettings{}, err
	}

	return mapPrivacySettings(row), nil
}

func (repo *PostgresRepository) GetPrivacySettingsByUserID(ctx context.Context, userID uuid.UUID) (domain.PrivacySettings, error) {
	row, err := repo.queries(ctx).GetPrivacySettingsByUserID(ctx, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.PrivacySettings{}, accountapp.NotFound("account privacy settings not found", err)
		}

		return domain.PrivacySettings{}, err
	}

	return mapPrivacySettings(row), nil
}

func (repo *PostgresRepository) UpdatePrivacySettings(ctx context.Context, settings domain.PrivacySettings) (domain.PrivacySettings, error) {
	row, err := repo.queries(ctx).UpdatePrivacySettings(ctx, accountsqlc.UpdatePrivacySettingsParams{
		UserID:            settings.UserID,
		ProfileVisibility: settings.ProfileVisibility,
		UpdatedAt:         settings.UpdatedAt,
	})
	if err != nil {
		return domain.PrivacySettings{}, err
	}

	return mapPrivacySettings(row), nil
}

func (uow *PostgresUnitOfWork) WithinTransaction(ctx context.Context, fn func(ctx context.Context, repo accountapp.Repository) error) error {
	if uow.txManager == nil {
		return fn(ctx, uow.repo)
	}

	return uow.txManager.WithinTransaction(ctx, func(ctx context.Context, _ *sql.Tx) error {
		return fn(ctx, uow.repo)
	})
}

func (repo *PostgresRepository) queries(ctx context.Context) *accountsqlc.Queries {
	if tx, ok := platformdb.TxFromContext(ctx); ok {
		return accountsqlc.New(tx)
	}

	return accountsqlc.New(repo.db)
}

func mapProfile(row accountsqlc.AccountProfile) domain.Profile {
	return domain.Profile{
		UserID:      row.UserID,
		Username:    row.Username,
		DisplayName: row.DisplayName,
		Bio:         row.Bio,
		AvatarPath:  row.AvatarPath,
		BannerPath:  row.BannerPath,
		CreatedAt:   row.CreatedAt,
		UpdatedAt:   row.UpdatedAt,
	}
}

func mapSettings(row accountsqlc.AccountSetting) domain.Settings {
	return domain.Settings{
		UserID:    row.UserID,
		Locale:    row.Locale,
		Timezone:  row.Timezone,
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
	}
}

func mapPrivacySettings(row accountsqlc.AccountPrivacySetting) domain.PrivacySettings {
	return domain.PrivacySettings{
		UserID:            row.UserID,
		ProfileVisibility: row.ProfileVisibility,
		CreatedAt:         row.CreatedAt,
		UpdatedAt:         row.UpdatedAt,
	}
}
