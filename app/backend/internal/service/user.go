package service

import (
	"context"
	"strings"

	"github.com/apk471/go-boilerplate/internal/errs"
	"github.com/apk471/go-boilerplate/internal/model"
	"github.com/apk471/go-boilerplate/internal/repository"
	"github.com/apk471/go-boilerplate/internal/server"
	"github.com/google/uuid"
)

type userRepository interface {
	CountUsers(ctx context.Context) (int, error)
	CreateUser(ctx context.Context, params repository.CreateUserParams) (*model.User, error)
	ListUsers(ctx context.Context) ([]model.User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*model.User, error)
	GetUserByAuthUserID(ctx context.Context, authUserID string) (*model.User, error)
	UpdateUser(ctx context.Context, id uuid.UUID, params repository.UpdateUserParams) (*model.User, error)
}

type UserService struct {
	server   *server.Server
	userRepo userRepository
}

type CreateUserInput struct {
	AuthUserID string
	Email      string
	Name       string
	Role       model.UserRole
	Status     model.UserStatus
}

type BootstrapUserInput struct {
	Email string
	Name  string
}

type UpdateUserInput struct {
	Email  *string
	Name   *string
	Role   *model.UserRole
	Status *model.UserStatus
}

func NewUserService(s *server.Server, userRepo userRepository) *UserService {
	return &UserService{
		server:   s,
		userRepo: userRepo,
	}
}

func (s *UserService) BootstrapUser(ctx context.Context, authUserID string, input BootstrapUserInput) (*model.User, error) {
	count, err := s.userRepo.CountUsers(ctx)
	if err != nil {
		return nil, err
	}

	if count > 0 {
		return nil, errs.NewForbiddenError("Initial admin has already been created", true)
	}

	return s.createUser(ctx, CreateUserInput{
		AuthUserID: authUserID,
		Email:      input.Email,
		Name:       input.Name,
		Role:       model.UserRoleAdmin,
		Status:     model.UserStatusActive,
	})
}

func (s *UserService) CreateUser(ctx context.Context, input CreateUserInput) (*model.User, error) {
	return s.createUser(ctx, input)
}

func (s *UserService) createUser(ctx context.Context, input CreateUserInput) (*model.User, error) {
	role := model.NormalizeRole(input.Role)
	status := model.NormalizeStatus(input.Status)
	if status == "" {
		status = model.UserStatusActive
	}

	if !model.IsValidRole(role) {
		return nil, errs.NewBadRequestError("Invalid role", true, nil, []errs.FieldError{
			{Field: "role", Error: "must be one of: viewer analyst admin"},
		}, nil)
	}

	if !model.IsValidStatus(status) {
		return nil, errs.NewBadRequestError("Invalid status", true, nil, []errs.FieldError{
			{Field: "status", Error: "must be one of: active inactive"},
		}, nil)
	}

	return s.userRepo.CreateUser(ctx, repository.CreateUserParams{
		AuthUserID: strings.TrimSpace(input.AuthUserID),
		Email:      strings.ToLower(strings.TrimSpace(input.Email)),
		Name:       strings.TrimSpace(input.Name),
		Role:       role,
		Status:     status,
	})
}

func (s *UserService) ListUsers(ctx context.Context) ([]model.User, error) {
	return s.userRepo.ListUsers(ctx)
}

func (s *UserService) GetUserByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	return s.userRepo.GetUserByID(ctx, id)
}

func (s *UserService) GetUserByAuthUserID(ctx context.Context, authUserID string) (*model.User, error) {
	return s.userRepo.GetUserByAuthUserID(ctx, authUserID)
}

func (s *UserService) UpdateUser(ctx context.Context, id uuid.UUID, input UpdateUserInput) (*model.User, error) {
	existing, err := s.userRepo.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}

	email := existing.Email
	name := existing.Name
	role := existing.Role
	status := existing.Status

	if input.Email != nil {
		email = strings.ToLower(strings.TrimSpace(*input.Email))
	}
	if input.Name != nil {
		name = strings.TrimSpace(*input.Name)
	}
	if input.Role != nil {
		role = model.NormalizeRole(*input.Role)
	}
	if input.Status != nil {
		status = model.NormalizeStatus(*input.Status)
	}

	if !model.IsValidRole(role) {
		return nil, errs.NewBadRequestError("Invalid role", true, nil, []errs.FieldError{
			{Field: "role", Error: "must be one of: viewer analyst admin"},
		}, nil)
	}

	if !model.IsValidStatus(status) {
		return nil, errs.NewBadRequestError("Invalid status", true, nil, []errs.FieldError{
			{Field: "status", Error: "must be one of: active inactive"},
		}, nil)
	}

	return s.userRepo.UpdateUser(ctx, id, repository.UpdateUserParams{
		Email:  email,
		Name:   name,
		Role:   role,
		Status: status,
	})
}

func (s *UserService) EnsureUserIsActive(ctx context.Context, authUserID string) (*model.User, error) {
	user, err := s.userRepo.GetUserByAuthUserID(ctx, authUserID)
	if err != nil {
		return nil, errs.NewForbiddenError("User is not provisioned for this application", true)
	}

	if user.Status != model.UserStatusActive {
		return nil, errs.NewForbiddenError("User account is inactive", true)
	}

	return user, nil
}
