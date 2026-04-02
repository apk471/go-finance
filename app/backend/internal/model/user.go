package model

import "strings"

type UserRole string

const (
	UserRoleViewer  UserRole = "viewer"
	UserRoleAnalyst UserRole = "analyst"
	UserRoleAdmin   UserRole = "admin"
)

type UserStatus string

const (
	UserStatusActive   UserStatus = "active"
	UserStatusInactive UserStatus = "inactive"
)

type User struct {
	Base
	AuthUserID string     `json:"authUserId" db:"auth_user_id"`
	Email      string     `json:"email" db:"email"`
	Name       string     `json:"name" db:"name"`
	Role       UserRole   `json:"role" db:"role"`
	Status     UserStatus `json:"status" db:"status"`
}

func ValidUserRoles() []UserRole {
	return []UserRole{UserRoleViewer, UserRoleAnalyst, UserRoleAdmin}
}

func ValidUserStatuses() []UserStatus {
	return []UserStatus{UserStatusActive, UserStatusInactive}
}

func NormalizeRole(role UserRole) UserRole {
	return UserRole(strings.ToLower(string(role)))
}

func NormalizeStatus(status UserStatus) UserStatus {
	return UserStatus(strings.ToLower(string(status)))
}

func IsValidRole(role UserRole) bool {
	switch NormalizeRole(role) {
	case UserRoleViewer, UserRoleAnalyst, UserRoleAdmin:
		return true
	default:
		return false
	}
}

func IsValidStatus(status UserStatus) bool {
	switch NormalizeStatus(status) {
	case UserStatusActive, UserStatusInactive:
		return true
	default:
		return false
	}
}

func RoleRank(role UserRole) int {
	switch NormalizeRole(role) {
	case UserRoleAdmin:
		return 3
	case UserRoleAnalyst:
		return 2
	case UserRoleViewer:
		return 1
	default:
		return 0
	}
}

func RoleAtLeast(role UserRole, required UserRole) bool {
	return RoleRank(role) >= RoleRank(required)
}
