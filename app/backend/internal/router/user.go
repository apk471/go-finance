package router

import (
	"github.com/apk471/go-boilerplate/internal/handler"
	"github.com/apk471/go-boilerplate/internal/middleware"
	"github.com/labstack/echo/v4"
)

func registerUserRoutes(group *echo.Group, h *handler.Handlers, auth *middleware.AuthMiddleware) {
	h.User.RegisterRoutes(group, auth)
}
