package domain

import "github.com/google/uuid"

type Principal struct {
	IsGuest        bool
	UserID         *uuid.UUID
	SessionID      *uuid.UUID
	Email          string
	BaseRole       string
	SubRoles       []SubRole
	PermissionKeys []string
}

func GuestPrincipal() Principal {
	return Principal{
		IsGuest:  true,
		BaseRole: "guest",
	}
}

func (principal Principal) IsAuthenticated() bool {
	return !principal.IsGuest && principal.UserID != nil
}

func (principal Principal) HasBaseRole(role string) bool {
	if principal.BaseRole == "admin" {
		return true
	}

	return principal.BaseRole == role
}

func (principal Principal) HasAnyBaseRole(roles ...string) bool {
	if principal.BaseRole == "admin" {
		return true
	}

	for _, role := range roles {
		if principal.BaseRole == role {
			return true
		}
	}

	return false
}

func (principal Principal) HasPermission(permission string) bool {
	if principal.BaseRole == "admin" {
		return true
	}

	for _, candidate := range principal.PermissionKeys {
		if candidate == permission {
			return true
		}
	}

	return false
}
