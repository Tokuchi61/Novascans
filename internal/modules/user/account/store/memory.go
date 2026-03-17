package store

import (
	"context"
	"sync"

	"github.com/google/uuid"

	accountapp "github.com/Tokuchi61/Novascans/internal/modules/user/account/app"
	"github.com/Tokuchi61/Novascans/internal/modules/user/account/domain"
)

type MemoryRepository struct {
	mu              sync.RWMutex
	profiles        map[uuid.UUID]domain.Profile
	profilesByName  map[string]uuid.UUID
	settings        map[uuid.UUID]domain.Settings
	privacySettings map[uuid.UUID]domain.PrivacySettings
}

type MemoryUnitOfWork struct {
	repo *MemoryRepository
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		profiles:        make(map[uuid.UUID]domain.Profile),
		profilesByName:  make(map[string]uuid.UUID),
		settings:        make(map[uuid.UUID]domain.Settings),
		privacySettings: make(map[uuid.UUID]domain.PrivacySettings),
	}
}

func NewMemoryUnitOfWork(repo *MemoryRepository) *MemoryUnitOfWork {
	return &MemoryUnitOfWork{repo: repo}
}

func (repo *MemoryRepository) CreateProfile(_ context.Context, profile domain.Profile) (domain.Profile, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	repo.profiles[profile.UserID] = profile
	repo.profilesByName[profile.Username] = profile.UserID
	return profile, nil
}

func (repo *MemoryRepository) GetProfileByUserID(_ context.Context, userID uuid.UUID) (domain.Profile, error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	profile, ok := repo.profiles[userID]
	if !ok {
		return domain.Profile{}, accountapp.NotFound("account profile not found", nil)
	}

	return profile, nil
}

func (repo *MemoryRepository) GetProfileByUsername(_ context.Context, username string) (domain.Profile, error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	userID, ok := repo.profilesByName[username]
	if !ok {
		return domain.Profile{}, accountapp.NotFound("account profile not found", nil)
	}

	profile, ok := repo.profiles[userID]
	if !ok {
		return domain.Profile{}, accountapp.NotFound("account profile not found", nil)
	}

	return profile, nil
}

func (repo *MemoryRepository) UpdateProfile(_ context.Context, profile domain.Profile) (domain.Profile, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	current, ok := repo.profiles[profile.UserID]
	if !ok {
		return domain.Profile{}, accountapp.NotFound("account profile not found", nil)
	}

	if current.Username != profile.Username {
		delete(repo.profilesByName, current.Username)
		repo.profilesByName[profile.Username] = profile.UserID
	}

	repo.profiles[profile.UserID] = profile
	return profile, nil
}

func (repo *MemoryRepository) CreateSettings(_ context.Context, settings domain.Settings) (domain.Settings, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	repo.settings[settings.UserID] = settings
	return settings, nil
}

func (repo *MemoryRepository) GetSettingsByUserID(_ context.Context, userID uuid.UUID) (domain.Settings, error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	settings, ok := repo.settings[userID]
	if !ok {
		return domain.Settings{}, accountapp.NotFound("account settings not found", nil)
	}

	return settings, nil
}

func (repo *MemoryRepository) UpdateSettings(_ context.Context, settings domain.Settings) (domain.Settings, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	if _, ok := repo.settings[settings.UserID]; !ok {
		return domain.Settings{}, accountapp.NotFound("account settings not found", nil)
	}

	repo.settings[settings.UserID] = settings
	return settings, nil
}

func (repo *MemoryRepository) CreatePrivacySettings(_ context.Context, settings domain.PrivacySettings) (domain.PrivacySettings, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	repo.privacySettings[settings.UserID] = settings
	return settings, nil
}

func (repo *MemoryRepository) GetPrivacySettingsByUserID(_ context.Context, userID uuid.UUID) (domain.PrivacySettings, error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	settings, ok := repo.privacySettings[userID]
	if !ok {
		return domain.PrivacySettings{}, accountapp.NotFound("account privacy settings not found", nil)
	}

	return settings, nil
}

func (repo *MemoryRepository) UpdatePrivacySettings(_ context.Context, settings domain.PrivacySettings) (domain.PrivacySettings, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	if _, ok := repo.privacySettings[settings.UserID]; !ok {
		return domain.PrivacySettings{}, accountapp.NotFound("account privacy settings not found", nil)
	}

	repo.privacySettings[settings.UserID] = settings
	return settings, nil
}

func (uow *MemoryUnitOfWork) WithinTransaction(ctx context.Context, fn func(ctx context.Context, repo accountapp.Repository) error) error {
	return fn(ctx, uow.repo)
}
