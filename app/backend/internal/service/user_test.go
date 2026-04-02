package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/apk471/go-boilerplate/internal/errs"
	"github.com/apk471/go-boilerplate/internal/model"
	"github.com/apk471/go-boilerplate/internal/repository"
	"github.com/google/uuid"
)

type stubUserRepository struct {
	countUsersResult int
	countUsersErr    error
	createInput      repository.CreateUserParams
	createResult     *model.User
	createErr        error
	getByIDResult    *model.User
	getByIDErr       error
	getByAuthResult  *model.User
	getByAuthErr     error
	updateInput      repository.UpdateUserParams
	updateResult     *model.User
	updateErr        error
}

func (s *stubUserRepository) CountUsers(ctx context.Context) (int, error) {
	return s.countUsersResult, s.countUsersErr
}

func (s *stubUserRepository) CreateUser(ctx context.Context, params repository.CreateUserParams) (*model.User, error) {
	s.createInput = params
	if s.createResult != nil {
		return s.createResult, s.createErr
	}

	return &model.User{
		Base: model.Base{
			BaseWithId:        model.BaseWithId{ID: uuid.New()},
			BaseWithCreatedAt: model.BaseWithCreatedAt{CreatedAt: time.Now()},
			BaseWithUpdatedAt: model.BaseWithUpdatedAt{UpdatedAt: time.Now()},
		},
		AuthUserID: params.AuthUserID,
		Email:      params.Email,
		Name:       params.Name,
		Role:       params.Role,
		Status:     params.Status,
	}, s.createErr
}

func (s *stubUserRepository) ListUsers(ctx context.Context) ([]model.User, error) {
	return nil, nil
}

func (s *stubUserRepository) GetUserByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	return s.getByIDResult, s.getByIDErr
}

func (s *stubUserRepository) GetUserByAuthUserID(ctx context.Context, authUserID string) (*model.User, error) {
	return s.getByAuthResult, s.getByAuthErr
}

func (s *stubUserRepository) UpdateUser(ctx context.Context, id uuid.UUID, params repository.UpdateUserParams) (*model.User, error) {
	s.updateInput = params
	if s.updateResult != nil {
		return s.updateResult, s.updateErr
	}

	return &model.User{
		Base: model.Base{
			BaseWithId:        model.BaseWithId{ID: id},
			BaseWithCreatedAt: model.BaseWithCreatedAt{CreatedAt: time.Now()},
			BaseWithUpdatedAt: model.BaseWithUpdatedAt{UpdatedAt: time.Now()},
		},
		AuthUserID: "auth_123",
		Email:      params.Email,
		Name:       params.Name,
		Role:       params.Role,
		Status:     params.Status,
	}, s.updateErr
}

func TestBootstrapUserCreatesInitialAdmin(t *testing.T) {
	repo := &stubUserRepository{}
	svc := NewUserService(nil, repo)

	user, err := svc.BootstrapUser(context.Background(), "auth_123", BootstrapUserInput{
		Email: "Admin@Example.com",
		Name:  "Admin User",
	})
	if err != nil {
		t.Fatalf("BootstrapUser() error = %v", err)
	}

	if user.Role != model.UserRoleAdmin {
		t.Fatalf("expected admin role, got %s", user.Role)
	}

	if repo.createInput.Status != model.UserStatusActive {
		t.Fatalf("expected active status, got %s", repo.createInput.Status)
	}

	if repo.createInput.Email != "admin@example.com" {
		t.Fatalf("expected normalized email, got %s", repo.createInput.Email)
	}
}

func TestBootstrapUserFailsWhenUsersExist(t *testing.T) {
	repo := &stubUserRepository{countUsersResult: 1}
	svc := NewUserService(nil, repo)

	_, err := svc.BootstrapUser(context.Background(), "auth_123", BootstrapUserInput{
		Email: "admin@example.com",
		Name:  "Admin User",
	})

	var httpErr *errs.HTTPError
	if err == nil {
		t.Fatal("expected error")
	}
	if !errors.As(err, &httpErr) {
		t.Fatalf("expected HTTPError, got %T", err)
	}
	if httpErr.Status != 403 {
		t.Fatalf("expected 403 status, got %d", httpErr.Status)
	}
}

func TestUpdateUserMergesFieldsAndNormalizesRole(t *testing.T) {
	repo := &stubUserRepository{
		getByIDResult: &model.User{
			Base:   model.Base{BaseWithId: model.BaseWithId{ID: uuid.New()}},
			Email:  "viewer@example.com",
			Name:   "Viewer",
			Role:   model.UserRoleViewer,
			Status: model.UserStatusActive,
		},
	}
	svc := NewUserService(nil, repo)

	role := model.UserRole("ADMIN")
	name := "Updated Viewer"

	user, err := svc.UpdateUser(context.Background(), repo.getByIDResult.ID, UpdateUserInput{
		Name: &name,
		Role: &role,
	})
	if err != nil {
		t.Fatalf("UpdateUser() error = %v", err)
	}

	if user.Role != model.UserRoleAdmin {
		t.Fatalf("expected admin role, got %s", user.Role)
	}

	if repo.updateInput.Email != "viewer@example.com" {
		t.Fatalf("expected existing email to be preserved, got %s", repo.updateInput.Email)
	}
}

func TestEnsureUserIsActiveRejectsInactiveUsers(t *testing.T) {
	repo := &stubUserRepository{
		getByAuthResult: &model.User{
			Base:   model.Base{BaseWithId: model.BaseWithId{ID: uuid.New()}},
			Status: model.UserStatusInactive,
		},
	}
	svc := NewUserService(nil, repo)

	_, err := svc.EnsureUserIsActive(context.Background(), "auth_123")
	if err == nil {
		t.Fatal("expected error")
	}
}
