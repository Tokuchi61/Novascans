package app_test

import (
	"testing"
	"time"

	"github.com/google/uuid"

	accessdomain "github.com/Tokuchi61/Novascans/internal/modules/identity/access/domain"
	authdomain "github.com/Tokuchi61/Novascans/internal/modules/identity/auth/domain"
	accountapp "github.com/Tokuchi61/Novascans/internal/modules/user/account/app"
	"github.com/Tokuchi61/Novascans/internal/modules/user/account/domain"
	accountstore "github.com/Tokuchi61/Novascans/internal/modules/user/account/store"
)

func TestProvisionDefaultsGeneratesUniqueUsername(t *testing.T) {
	repo := accountstore.NewMemoryRepository()
	service := accountapp.NewService(repo, accountstore.NewMemoryUnitOfWork(repo))
	now := time.Now().UTC()

	userOne := authdomain.User{
		ID:        uuid.New(),
		Email:     "reader@example.com",
		BaseRole:  accessdomain.BaseRoleUser,
		Status:    authdomain.StatusActive,
		CreatedAt: now,
		UpdatedAt: now,
	}
	userTwo := authdomain.User{
		ID:        uuid.New(),
		Email:     "reader@other.example",
		BaseRole:  accessdomain.BaseRoleUser,
		Status:    authdomain.StatusActive,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := service.ProvisionDefaults(t.Context(), userOne, now); err != nil {
		t.Fatalf("provision defaults for user one: %v", err)
	}

	if err := service.ProvisionDefaults(t.Context(), userTwo, now); err != nil {
		t.Fatalf("provision defaults for user two: %v", err)
	}

	profileOne, err := repo.GetProfileByUserID(t.Context(), userOne.ID)
	if err != nil {
		t.Fatalf("get profile for user one: %v", err)
	}

	profileTwo, err := repo.GetProfileByUserID(t.Context(), userTwo.ID)
	if err != nil {
		t.Fatalf("get profile for user two: %v", err)
	}

	if profileOne.Username != "reader" {
		t.Fatalf("expected first username to be reader, got %q", profileOne.Username)
	}

	if profileTwo.Username == profileOne.Username {
		t.Fatalf("expected unique username for second user, got duplicate %q", profileTwo.Username)
	}

	if _, err := repo.GetSettingsByUserID(t.Context(), userTwo.ID); err != nil {
		t.Fatalf("expected settings for user two: %v", err)
	}

	if _, err := repo.GetPrivacySettingsByUserID(t.Context(), userTwo.ID); err != nil {
		t.Fatalf("expected privacy settings for user two: %v", err)
	}
}

func TestGetPublicProfileHonorsVisibility(t *testing.T) {
	repo := accountstore.NewMemoryRepository()
	service := accountapp.NewService(repo, accountstore.NewMemoryUnitOfWork(repo))
	now := time.Now().UTC()

	user := authdomain.User{
		ID:        uuid.New(),
		Email:     "owner@example.com",
		BaseRole:  accessdomain.BaseRoleUser,
		Status:    authdomain.StatusActive,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := service.ProvisionDefaults(t.Context(), user, now); err != nil {
		t.Fatalf("provision defaults: %v", err)
	}

	profile, err := repo.GetProfileByUserID(t.Context(), user.ID)
	if err != nil {
		t.Fatalf("get profile: %v", err)
	}

	if _, err := service.GetPublicProfile(t.Context(), profile.Username, accountapp.Viewer{}); err != nil {
		t.Fatalf("expected guest to read public profile: %v", err)
	}

	if _, err := service.UpdatePrivacy(t.Context(), accountapp.UpdatePrivacyInput{
		UserID:            user.ID,
		ProfileVisibility: domain.ProfileVisibilityAuthenticated,
	}); err != nil {
		t.Fatalf("update privacy to authenticated: %v", err)
	}

	if _, err := service.GetPublicProfile(t.Context(), profile.Username, accountapp.Viewer{}); !accountapp.HasCode(err, accountapp.CodeNotFound) {
		t.Fatalf("expected guest lookup to fail for authenticated profile, got %v", err)
	}

	if _, err := service.GetPublicProfile(t.Context(), profile.Username, accountapp.Viewer{Authenticated: true}); err != nil {
		t.Fatalf("expected authenticated viewer to read profile: %v", err)
	}

	if _, err := service.UpdatePrivacy(t.Context(), accountapp.UpdatePrivacyInput{
		UserID:            user.ID,
		ProfileVisibility: domain.ProfileVisibilityPrivate,
	}); err != nil {
		t.Fatalf("update privacy to private: %v", err)
	}

	if _, err := service.GetPublicProfile(t.Context(), profile.Username, accountapp.Viewer{Authenticated: true}); !accountapp.HasCode(err, accountapp.CodeNotFound) {
		t.Fatalf("expected other authenticated viewer to fail for private profile, got %v", err)
	}

	if _, err := service.GetPublicProfile(t.Context(), profile.Username, accountapp.Viewer{
		Authenticated: true,
		UserID:        &user.ID,
	}); err != nil {
		t.Fatalf("expected owner to read private profile: %v", err)
	}
}
