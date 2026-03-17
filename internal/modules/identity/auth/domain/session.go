package domain

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID                uuid.UUID
	UserID            uuid.UUID
	TokenHash         string
	UserAgent         string
	IPAddress         string
	ExpiresAt         time.Time
	RevokedAt         *time.Time
	ReplacedBySession *uuid.UUID
	LastUsedAt        *time.Time
	CreatedAt         time.Time
}
