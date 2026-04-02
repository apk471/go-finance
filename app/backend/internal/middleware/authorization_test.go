package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/apk471/go-boilerplate/internal/errs"
	"github.com/apk471/go-boilerplate/internal/model"
	"github.com/labstack/echo/v4"
)

func TestHasPermission(t *testing.T) {
	testCases := []struct {
		name       string
		role       model.UserRole
		permission Permission
		want       bool
	}{
		{
			name:       "viewer can read records",
			role:       model.UserRoleViewer,
			permission: PermissionReadRecords,
			want:       true,
		},
		{
			name:       "viewer cannot manage records",
			role:       model.UserRoleViewer,
			permission: PermissionManageRecords,
			want:       false,
		},
		{
			name:       "analyst can access summaries",
			role:       model.UserRoleAnalyst,
			permission: PermissionAccessSummaries,
			want:       true,
		},
		{
			name:       "analyst cannot delete records",
			role:       model.UserRoleAnalyst,
			permission: PermissionDeleteRecords,
			want:       false,
		},
		{
			name:       "admin can manage users",
			role:       model.UserRoleAdmin,
			permission: PermissionManageUsers,
			want:       true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if got := HasPermission(tc.role, tc.permission); got != tc.want {
				t.Fatalf("HasPermission(%q, %q) = %v, want %v", tc.role, tc.permission, got, tc.want)
			}
		})
	}
}

func TestRequirePermission(t *testing.T) {
	testCases := []struct {
		name       string
		role       model.UserRole
		permission Permission
		wantStatus int
		wantCalled bool
	}{
		{
			name:       "viewer is denied record writes",
			role:       model.UserRoleViewer,
			permission: PermissionManageRecords,
			wantStatus: http.StatusForbidden,
			wantCalled: false,
		},
		{
			name:       "analyst is allowed dashboard summaries",
			role:       model.UserRoleAnalyst,
			permission: PermissionAccessSummaries,
			wantStatus: http.StatusOK,
			wantCalled: true,
		},
		{
			name:       "analyst is denied user management",
			role:       model.UserRoleAnalyst,
			permission: PermissionManageUsers,
			wantStatus: http.StatusForbidden,
			wantCalled: false,
		},
		{
			name:       "admin can delete records",
			role:       model.UserRoleAdmin,
			permission: PermissionDeleteRecords,
			wantStatus: http.StatusOK,
			wantCalled: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			auth := &AuthMiddleware{}
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.Set(CurrentUserKey, &model.User{Role: tc.role})

			called := false
			handler := auth.RequirePermission(tc.permission)(func(c echo.Context) error {
				called = true
				return c.NoContent(http.StatusOK)
			})

			err := handler(c)
			if tc.wantCalled {
				if err != nil {
					t.Fatalf("RequirePermission() error = %v", err)
				}
				if rec.Code != tc.wantStatus {
					t.Fatalf("expected status %d, got %d", tc.wantStatus, rec.Code)
				}
			} else {
				if err == nil {
					t.Fatal("expected error but got nil")
				}
				appErr, ok := err.(*errs.HTTPError)
				if !ok {
					t.Fatalf("expected *errs.HTTPError, got %T", err)
				}
				if appErr.Status != tc.wantStatus {
					t.Fatalf("expected status %d, got %d", tc.wantStatus, appErr.Status)
				}
			}

			if called != tc.wantCalled {
				t.Fatalf("expected next handler called=%v, got %v", tc.wantCalled, called)
			}
		})
	}
}
