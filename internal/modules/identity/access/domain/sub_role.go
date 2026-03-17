package domain

import (
	"time"

	"github.com/google/uuid"
)

type SubRole struct {
	ID             uuid.UUID
	Key            string
	Name           string
	Description    string
	PermissionKeys []string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
