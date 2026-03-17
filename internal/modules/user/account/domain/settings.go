package domain

import (
	"time"

	"github.com/google/uuid"
)

type Settings struct {
	UserID    uuid.UUID
	Locale    string
	Timezone  string
	CreatedAt time.Time
	UpdatedAt time.Time
}
