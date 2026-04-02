package router

import (
	"github.com/apk471/go-boilerplate/internal/handler"
	"github.com/apk471/go-boilerplate/internal/middleware"
	"github.com/labstack/echo/v4"
)

func registerFinancialRecordRoutes(group *echo.Group, h *handler.Handlers, auth *middleware.AuthMiddleware) {
	h.FinancialRecord.RegisterRoutes(group, auth)
}
