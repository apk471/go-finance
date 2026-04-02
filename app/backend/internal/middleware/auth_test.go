package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/apk471/go-boilerplate/internal/config"
	"github.com/apk471/go-boilerplate/internal/server"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

func TestRequireAuthUsesLocalDevIdentityWithoutAuthorizationHeader(t *testing.T) {
	auth := &AuthMiddleware{
		server: &server.Server{
			Config: &config.Config{
				Primary: config.Primary{Env: "local"},
			},
		},
	}

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/users/me", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	called := false
	handler := auth.RequireAuth(func(c echo.Context) error {
		called = true

		if got := GetAuthUserID(c); got != defaultLocalDevAuthUserID {
			t.Fatalf("expected auth user id %q, got %q", defaultLocalDevAuthUserID, got)
		}
		if got := GetUserID(c); got != defaultLocalDevAuthUserID {
			t.Fatalf("expected user id %q, got %q", defaultLocalDevAuthUserID, got)
		}
		if got, _ := c.Get("user_role").(string); got != string(defaultDevRole()) {
			t.Fatalf("expected user role %q, got %q", defaultDevRole(), got)
		}

		return nil
	})

	if err := handler(c); err != nil {
		t.Fatalf("RequireAuth() error = %v", err)
	}
	if !called {
		t.Fatal("expected next handler to be called")
	}
}

func TestRequireAuthUsesExplicitLocalDevHeaders(t *testing.T) {
	auth := &AuthMiddleware{
		server: &server.Server{
			Config: &config.Config{
				Primary: config.Primary{Env: "local"},
			},
			Logger: func() *zerolog.Logger {
				logger := zerolog.Nop()
				return &logger
			}(),
		},
	}

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/users/me", nil)
	req.Header.Set(DevAuthUserIDHeader, "dev-analyst")
	req.Header.Set(DevUserRoleHeader, "analyst")

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := auth.RequireAuth(func(c echo.Context) error {
		if got := GetAuthUserID(c); got != "dev-analyst" {
			t.Fatalf("expected auth user id %q, got %q", "dev-analyst", got)
		}
		if got, _ := c.Get("user_role").(string); got != "analyst" {
			t.Fatalf("expected user role %q, got %q", "analyst", got)
		}
		return nil
	})

	if err := handler(c); err != nil {
		t.Fatalf("RequireAuth() error = %v", err)
	}
}

func defaultDevRole() string {
	return "admin"
}
