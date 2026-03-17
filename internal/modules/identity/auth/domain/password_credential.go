package domain

import (
	"time"

	"github.com/google/uuid"
)

type PasswordCredential struct {
	UserID       uuid.UUID
	PasswordHash string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
