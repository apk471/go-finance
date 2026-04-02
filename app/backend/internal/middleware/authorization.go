package middleware

import (
	"github.com/apk471/go-boilerplate/internal/errs"
	"github.com/apk471/go-boilerplate/internal/model"
	"github.com/labstack/echo/v4"
)

const (
	CurrentUserKey      = "current_user"
	AuthUserIDKey       = "auth_user_id"
	DevAuthUserIDHeader = "X-Dev-Auth-User-Id"
	DevUserRoleHeader   = "X-Dev-User-Role"
)

type Permission string

const (
	PermissionReadRecords     Permission = "read_records"
	PermissionAccessSummaries Permission = "access_summaries"
	PermissionManageRecords   Permission = "manage_records"
	PermissionDeleteRecords   Permission = "delete_records"
	PermissionManageUsers     Permission = "manage_users"
)

var rolePermissions = map[model.UserRole]map[Permission]struct{}{
	model.UserRoleViewer: {
		PermissionReadRecords:     {},
		PermissionAccessSummaries: {},
	},
	model.UserRoleAnalyst: {
		PermissionReadRecords:     {},
		PermissionAccessSummaries: {},
		PermissionManageRecords:   {},
	},
	model.UserRoleAdmin: {
		PermissionReadRecords:     {},
		PermissionAccessSummaries: {},
		PermissionManageRecords:   {},
		PermissionDeleteRecords:   {},
		PermissionManageUsers:     {},
	},
}

func GetCurrentUser(c echo.Context) *model.User {
	user, ok := c.Get(CurrentUserKey).(*model.User)
	if !ok {
		return nil
	}

	return user
}

func GetAuthUserID(c echo.Context) string {
	authUserID, _ := c.Get(AuthUserIDKey).(string)
	return authUserID
}

func HasPermission(role model.UserRole, permission Permission) bool {
	permissions, ok := rolePermissions[model.NormalizeRole(role)]
	if !ok {
		return false
	}

	_, ok = permissions[permission]
	return ok
}

func (auth *AuthMiddleware) RequireActiveUser(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authUserID := GetAuthUserID(c)
		if authUserID == "" {
			authUserID = GetUserID(c)
		}
		if authUserID == "" {
			return errs.NewUnauthorizedError("Unauthorized", false)
		}

		user, err := auth.userRepo.GetUserByAuthUserID(c.Request().Context(), authUserID)
		if err != nil {
			return errs.NewForbiddenError("User is not provisioned for this application", true)
		}

		if user.Status != model.UserStatusActive {
			return errs.NewForbiddenError("User account is inactive", true)
		}

		c.Set(CurrentUserKey, user)
		c.Set(UserRoleKey, string(user.Role))
		c.Set("app_user_id", user.ID.String())

		return next(c)
	}
}

func (auth *AuthMiddleware) RequireRoles(roles ...model.UserRole) echo.MiddlewareFunc {
	allowed := make(map[model.UserRole]struct{}, len(roles))
	for _, role := range roles {
		allowed[model.NormalizeRole(role)] = struct{}{}
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			user := GetCurrentUser(c)
			if user == nil {
				return errs.NewUnauthorizedError("Unauthorized", false)
			}

			if _, ok := allowed[model.NormalizeRole(user.Role)]; !ok {
				return errs.NewForbiddenError("Insufficient role for this action", true)
			}

			return next(c)
		}
	}
}

func (auth *AuthMiddleware) RequireMinimumRole(role model.UserRole) echo.MiddlewareFunc {
	required := model.NormalizeRole(role)

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			user := GetCurrentUser(c)
			if user == nil {
				return errs.NewUnauthorizedError("Unauthorized", false)
			}

			if !model.RoleAtLeast(user.Role, required) {
				return errs.NewForbiddenError("Insufficient role for this action", true)
			}

			return next(c)
		}
	}
}

func (auth *AuthMiddleware) RequirePermission(permission Permission) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			user := GetCurrentUser(c)
			if user == nil {
				return errs.NewUnauthorizedError("Unauthorized", false)
			}

			if !HasPermission(user.Role, permission) {
				return errs.NewForbiddenError("Insufficient permissions for this action", true)
			}

			return next(c)
		}
	}
}
