package domain

import (
	"time"

	"github.com/google/uuid"
)

type EmailVerificationToken struct {
	ID         uuid.UUID
	UserID     uuid.UUID
	TokenHash  string
	ExpiresAt  time.Time
	ConsumedAt *time.Time
	CreatedAt  time.Time
}
