package handler

import (
	"net/http"
	"strings"

	"github.com/apk471/go-boilerplate/internal/middleware"
	"github.com/apk471/go-boilerplate/internal/model"
	"github.com/apk471/go-boilerplate/internal/server"
	"github.com/apk471/go-boilerplate/internal/service"
	"github.com/apk471/go-boilerplate/internal/validation"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	Handler
	userService *service.UserService
}

type CreateUserRequest struct {
	AuthUserID string           `json:"authUserId" validate:"required"`
	Email      string           `json:"email" validate:"required,email"`
	Name       string           `json:"name" validate:"required,min=2,max=100"`
	Role       model.UserRole   `json:"role" validate:"required"`
	Status     model.UserStatus `json:"status"`
}

func (r *CreateUserRequest) Validate() error {
	return validator.New().Struct(r)
}

type BootstrapUserRequest struct {
	Email string `json:"email" validate:"required,email"`
	Name  string `json:"name" validate:"required,min=2,max=100"`
}

func (r *BootstrapUserRequest) Validate() error {
	return validator.New().Struct(r)
}

type UpdateUserRequest struct {
	ID     string            `param:"id" validate:"required,uuid"`
	Email  *string           `json:"email" validate:"omitempty,email"`
	Name   *string           `json:"name" validate:"omitempty,min=2,max=100"`
	Role   *model.UserRole   `json:"role"`
	Status *model.UserStatus `json:"status"`
}

func (r *UpdateUserRequest) Validate() error {
	if err := validator.New().Struct(r); err != nil {
		return err
	}

	fieldErrors := validation.CustomValidationErrors{}
	if r.Role != nil && !model.IsValidRole(*r.Role) {
		fieldErrors = append(fieldErrors, validation.CustomValidationError{
			Field:   "role",
			Message: "must be one of: viewer analyst admin",
		})
	}
	if r.Status != nil && !model.IsValidStatus(*r.Status) {
		fieldErrors = append(fieldErrors, validation.CustomValidationError{
			Field:   "status",
			Message: "must be one of: active inactive",
		})
	}

	if len(fieldErrors) > 0 {
		return fieldErrors
	}

	return nil
}

type GetUserRequest struct {
	ID string `param:"id" validate:"required,uuid"`
}

func (r *GetUserRequest) Validate() error {
	return validator.New().Struct(r)
}

type EmptyRequest struct{}

func (r *EmptyRequest) Validate() error {
	return nil
}

func NewUserHandler(s *server.Server, userService *service.UserService) *UserHandler {
	return &UserHandler{
		Handler:     NewHandler(s),
		userService: userService,
	}
}

func (h *UserHandler) BootstrapUser(c echo.Context, req *BootstrapUserRequest) (*model.User, error) {
	return h.userService.BootstrapUser(c.Request().Context(), middleware.GetAuthUserID(c), service.BootstrapUserInput{
		Email: req.Email,
		Name:  req.Name,
	})
}

func (h *UserHandler) ListUsers(c echo.Context, req *EmptyRequest) ([]model.User, error) {
	return h.userService.ListUsers(c.Request().Context())
}

func (h *UserHandler) GetCurrentUser(c echo.Context, req *EmptyRequest) (*model.User, error) {
	return middleware.GetCurrentUser(c), nil
}

func (h *UserHandler) GetUser(c echo.Context, req *GetUserRequest) (*model.User, error) {
	userID, err := uuid.Parse(req.ID)
	if err != nil {
		return nil, err
	}

	return h.userService.GetUserByID(c.Request().Context(), userID)
}

func (h *UserHandler) CreateUser(c echo.Context, req *CreateUserRequest) (*model.User, error) {
	status := req.Status
	if strings.TrimSpace(string(status)) == "" {
		status = model.UserStatusActive
	}

	return h.userService.CreateUser(c.Request().Context(), service.CreateUserInput{
		AuthUserID: req.AuthUserID,
		Email:      req.Email,
		Name:       req.Name,
		Role:       req.Role,
		Status:     status,
	})
}

func (h *UserHandler) UpdateUser(c echo.Context, req *UpdateUserRequest) (*model.User, error) {
	userID, err := uuid.Parse(req.ID)
	if err != nil {
		return nil, err
	}

	return h.userService.UpdateUser(c.Request().Context(), userID, service.UpdateUserInput{
		Email:  req.Email,
		Name:   req.Name,
		Role:   req.Role,
		Status: req.Status,
	})
}

func (h *UserHandler) RegisterRoutes(group *echo.Group, auth *middleware.AuthMiddleware) {
	group.GET("/users/me", Handle(h.Handler, h.GetCurrentUser, http.StatusOK, &EmptyRequest{}), auth.RequireAuth, auth.RequireActiveUser)
	group.POST("/users/bootstrap", Handle(h.Handler, h.BootstrapUser, http.StatusCreated, &BootstrapUserRequest{}), auth.RequireAuth)

	adminGroup := group.Group("/users", auth.RequireAuth, auth.RequireActiveUser, auth.RequireRoles(model.UserRoleAdmin))
	adminGroup.GET("", Handle(h.Handler, h.ListUsers, http.StatusOK, &EmptyRequest{}))
	adminGroup.GET("/:id", Handle(h.Handler, h.GetUser, http.StatusOK, &GetUserRequest{}))
	adminGroup.POST("", Handle(h.Handler, h.CreateUser, http.StatusCreated, &CreateUserRequest{}))
	adminGroup.PATCH("/:id", Handle(h.Handler, h.UpdateUser, http.StatusOK, &UpdateUserRequest{}))
}
