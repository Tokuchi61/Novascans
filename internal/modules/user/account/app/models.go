package app

import (
	"github.com/google/uuid"

	"github.com/Tokuchi61/Novascans/internal/modules/user/account/domain"
)

type Viewer struct {
	UserID        *uuid.UUID
	Authenticated bool
}

type AccountData struct {
	Profile  domain.Profile
	Settings domain.Settings
	Privacy  domain.PrivacySettings
}

type UpdateProfileInput struct {
	UserID      uuid.UUID
	Username    *string
	DisplayName *string
	Bio         *string
	AvatarPath  *string
	BannerPath  *string
}

type UpdateSettingsInput struct {
	UserID   uuid.UUID
	Locale   *string
	Timezone *string
}

type UpdatePrivacyInput struct {
	UserID            uuid.UUID
	ProfileVisibility string
}
