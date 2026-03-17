package domain

import (
	"time"

	"github.com/google/uuid"
)

type Profile struct {
	UserID      uuid.UUID
	Username    string
	DisplayName string
	Bio         string
	AvatarPath  string
	BannerPath  string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
