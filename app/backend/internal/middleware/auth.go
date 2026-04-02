package middleware

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/apk471/go-boilerplate/internal/errs"
	"github.com/apk471/go-boilerplate/internal/model"
	"github.com/apk471/go-boilerplate/internal/repository"
	"github.com/apk471/go-boilerplate/internal/server"
	"github.com/clerk/clerk-sdk-go/v2"
	clerkhttp "github.com/clerk/clerk-sdk-go/v2/http"
	"github.com/labstack/echo/v4"
)

const defaultLocalDevAuthUserID = "local-dev-user"

type AuthMiddleware struct {
	server   *server.Server
	userRepo *repository.UserRepository
}

func NewAuthMiddleware(s *server.Server) *AuthMiddleware {
	return &AuthMiddleware{
		server:   s,
		userRepo: repository.NewUserRepository(s),
	}
}

func (auth *AuthMiddleware) RequireAuth(next echo.HandlerFunc) echo.HandlerFunc {
	if auth.server.Config.Primary.Env == "local" {
		return func(c echo.Context) error {
			if auth.shouldUseLocalDevAuth(c) {
				auth.applyLocalDevAuth(c)
				return next(c)
			}

			return auth.requireClerkAuth(next)(c)
		}
	}

	return auth.requireClerkAuth(next)
}

func (auth *AuthMiddleware) requireClerkAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return echo.WrapMiddleware(
		clerkhttp.WithHeaderAuthorization(
			clerkhttp.AuthorizationFailureHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				start := time.Now()

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)

				response := map[string]string{
					"code":     "UNAUTHORIZED",
					"message":  "Unauthorized",
					"override": "false",
					"status":   "401",
				}

				if err := json.NewEncoder(w).Encode(response); err != nil {
					auth.server.Logger.Error().Err(err).Str("function", "RequireAuth").Dur(
						"duration", time.Since(start)).Msg("failed to write JSON response")
				} else {
					auth.server.Logger.Error().Str("function", "RequireAuth").Dur("duration", time.Since(start)).Msg(
						"could not get session claims from context")
				}
			}))))(func(c echo.Context) error {
		start := time.Now()
		claims, ok := clerk.SessionClaimsFromContext(c.Request().Context())

		if !ok {
			auth.server.Logger.Error().
				Str("function", "RequireAuth").
				Str("request_id", GetRequestID(c)).
				Dur("duration", time.Since(start)).
				Msg("could not get session claims from context")
			return errs.NewUnauthorizedError("Unauthorized", false)
		}

		c.Set("user_id", claims.Subject)
		c.Set(AuthUserIDKey, claims.Subject)
		c.Set("user_role", claims.ActiveOrganizationRole)
		c.Set("permissions", claims.Claims.ActiveOrganizationPermissions)

		auth.server.Logger.Info().
			Str("function", "RequireAuth").
			Str("user_id", claims.Subject).
			Str("request_id", GetRequestID(c)).
			Dur("duration", time.Since(start)).
			Msg("user authenticated successfully")

		return next(c)
	})
}

func (auth *AuthMiddleware) shouldUseLocalDevAuth(c echo.Context) bool {
	authHeader := strings.TrimSpace(c.Request().Header.Get(echo.HeaderAuthorization))
	devAuthUserID := strings.TrimSpace(c.Request().Header.Get(DevAuthUserIDHeader))

	return authHeader == "" || devAuthUserID != ""
}

func (auth *AuthMiddleware) applyLocalDevAuth(c echo.Context) {
	authUserID := strings.TrimSpace(c.Request().Header.Get(DevAuthUserIDHeader))
	if authUserID == "" {
		authUserID = defaultLocalDevAuthUserID
	}

	devRole := model.NormalizeRole(model.UserRole(strings.TrimSpace(c.Request().Header.Get(DevUserRoleHeader))))
	if !model.IsValidRole(devRole) {
		devRole = model.UserRoleAdmin
	}

	c.Set("user_id", authUserID)
	c.Set(AuthUserIDKey, authUserID)
	c.Set("user_role", string(devRole))
	c.Set("permissions", []string{"dev:*"})
}
