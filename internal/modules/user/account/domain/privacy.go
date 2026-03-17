package domain

import (
	"time"

	"github.com/google/uuid"
)

const (
	ProfileVisibilityPublic        = "public"
	ProfileVisibilityAuthenticated = "authenticated"
	ProfileVisibilityPrivate       = "private"
)

type PrivacySettings struct {
	UserID            uuid.UUID
	ProfileVisibility string
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

func IsValidProfileVisibility(value string) bool {
	switch value {
	case ProfileVisibilityPublic, ProfileVisibilityAuthenticated, ProfileVisibilityPrivate:
		return true
	default:
		return false
	}
}

func CanViewProfile(ownerID uuid.UUID, visibility string, viewerID *uuid.UUID, authenticated bool) bool {
	if viewerID != nil && *viewerID == ownerID {
		return true
	}

	switch visibility {
	case ProfileVisibilityPublic:
		return true
	case ProfileVisibilityAuthenticated:
		return authenticated
	case ProfileVisibilityPrivate:
		return false
	default:
		return false
	}
}
