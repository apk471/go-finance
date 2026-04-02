package repository

import (
	"context"
	"fmt"

	"github.com/apk471/go-boilerplate/internal/model"
	"github.com/apk471/go-boilerplate/internal/server"
	"github.com/google/uuid"
)

type UserRepository struct {
	server *server.Server
}

type CreateUserParams struct {
	AuthUserID string
	Email      string
	Name       string
	Role       model.UserRole
	Status     model.UserStatus
}

type UpdateUserParams struct {
	Email  string
	Name   string
	Role   model.UserRole
	Status model.UserStatus
}

func NewUserRepository(s *server.Server) *UserRepository {
	return &UserRepository{server: s}
}

func (r *UserRepository) CountUsers(ctx context.Context) (int, error) {
	const query = `SELECT COUNT(*) FROM users`

	var count int
	if err := r.server.DB.Pool.QueryRow(ctx, query).Scan(&count); err != nil {
		return 0, fmt.Errorf("count users: %w", err)
	}

	return count, nil
}

func (r *UserRepository) CreateUser(ctx context.Context, params CreateUserParams) (*model.User, error) {
	const query = `
		INSERT INTO users (auth_user_id, email, name, role, status)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, auth_user_id, email, name, role, status, created_at, updated_at
	`

	user := &model.User{}
	if err := r.server.DB.Pool.QueryRow(
		ctx,
		query,
		params.AuthUserID,
		params.Email,
		params.Name,
		params.Role,
		params.Status,
	).Scan(
		&user.ID,
		&user.AuthUserID,
		&user.Email,
		&user.Name,
		&user.Role,
		&user.Status,
		&user.CreatedAt,
		&user.UpdatedAt,
	); err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	return user, nil
}

func (r *UserRepository) ListUsers(ctx context.Context) ([]model.User, error) {
	const query = `
		SELECT id, auth_user_id, email, name, role, status, created_at, updated_at
		FROM users
		ORDER BY created_at ASC
	`

	rows, err := r.server.DB.Pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("list users: %w", err)
	}
	defer rows.Close()

	users := make([]model.User, 0)
	for rows.Next() {
		var user model.User
		if err := rows.Scan(
			&user.ID,
			&user.AuthUserID,
			&user.Email,
			&user.Name,
			&user.Role,
			&user.Status,
			&user.CreatedAt,
			&user.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan user: %w", err)
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate users: %w", err)
	}

	return users, nil
}

func (r *UserRepository) GetUserByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	const query = `
		SELECT id, auth_user_id, email, name, role, status, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	user := &model.User{}
	if err := r.server.DB.Pool.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.AuthUserID,
		&user.Email,
		&user.Name,
		&user.Role,
		&user.Status,
		&user.CreatedAt,
		&user.UpdatedAt,
	); err != nil {
		return nil, fmt.Errorf("get user by id: %w", err)
	}

	return user, nil
}

func (r *UserRepository) GetUserByAuthUserID(ctx context.Context, authUserID string) (*model.User, error) {
	const query = `
		SELECT id, auth_user_id, email, name, role, status, created_at, updated_at
		FROM users
		WHERE auth_user_id = $1
	`

	user := &model.User{}
	if err := r.server.DB.Pool.QueryRow(ctx, query, authUserID).Scan(
		&user.ID,
		&user.AuthUserID,
		&user.Email,
		&user.Name,
		&user.Role,
		&user.Status,
		&user.CreatedAt,
		&user.UpdatedAt,
	); err != nil {
		return nil, fmt.Errorf("get user by auth user id: %w", err)
	}

	return user, nil
}

func (r *UserRepository) UpdateUser(ctx context.Context, id uuid.UUID, params UpdateUserParams) (*model.User, error) {
	const query = `
		UPDATE users
		SET email = $2,
		    name = $3,
		    role = $4,
		    status = $5,
		    updated_at = NOW()
		WHERE id = $1
		RETURNING id, auth_user_id, email, name, role, status, created_at, updated_at
	`

	user := &model.User{}
	if err := r.server.DB.Pool.QueryRow(
		ctx,
		query,
		id,
		params.Email,
		params.Name,
		params.Role,
		params.Status,
	).Scan(
		&user.ID,
		&user.AuthUserID,
		&user.Email,
		&user.Name,
		&user.Role,
		&user.Status,
		&user.CreatedAt,
		&user.UpdatedAt,
	); err != nil {
		return nil, fmt.Errorf("update user: %w", err)
	}

	return user, nil
}
