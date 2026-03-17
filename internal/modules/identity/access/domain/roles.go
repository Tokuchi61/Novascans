package domain

const (
	BaseRoleGuest     = "guest"
	BaseRoleUser      = "user"
	BaseRoleModerator = "moderator"
	BaseRoleAdmin     = "admin"
)

func IsAssignableBaseRole(role string) bool {
	switch role {
	case BaseRoleUser, BaseRoleModerator, BaseRoleAdmin:
		return true
	default:
		return false
	}
}
