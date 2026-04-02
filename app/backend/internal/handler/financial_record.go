package handler

import (
	"net/http"

	"github.com/apk471/go-boilerplate/internal/middleware"
	"github.com/apk471/go-boilerplate/internal/model"
	"github.com/apk471/go-boilerplate/internal/server"
	"github.com/apk471/go-boilerplate/internal/service"
	"github.com/apk471/go-boilerplate/internal/validation"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type FinancialRecordHandler struct {
	Handler
	recordService *service.FinancialRecordService
}

type CreateFinancialRecordRequest struct {
	Amount   string           `json:"amount" validate:"required"`
	Type     model.RecordType `json:"type" validate:"required"`
	Category string           `json:"category" validate:"required,min=2,max=100"`
	Date     string           `json:"date" validate:"required"`
	Notes    *string          `json:"notes" validate:"omitempty,max=500"`
}

func (r *CreateFinancialRecordRequest) Validate() error {
	return validator.New().Struct(r)
}

type UpdateFinancialRecordRequest struct {
	ID       string            `param:"id" validate:"required,uuid"`
	Amount   *string           `json:"amount"`
	Type     *model.RecordType `json:"type"`
	Category *string           `json:"category" validate:"omitempty,min=2,max=100"`
	Date     *string           `json:"date"`
	Notes    *string           `json:"notes" validate:"omitempty,max=500"`
}

func (r *UpdateFinancialRecordRequest) Validate() error {
	fieldErrors := validation.CustomValidationErrors{}
	if err := validator.New().Struct(r); err != nil {
		return err
	}

	if r.Amount == nil && r.Type == nil && r.Category == nil && r.Date == nil && r.Notes == nil {
		fieldErrors = append(fieldErrors, validation.CustomValidationError{
			Field:   "body",
			Message: "at least one field must be provided",
		})
	}

	if len(fieldErrors) > 0 {
		return fieldErrors
	}

	return nil
}

type GetFinancialRecordRequest struct {
	ID string `param:"id" validate:"required,uuid"`
}

func (r *GetFinancialRecordRequest) Validate() error {
	return validator.New().Struct(r)
}

type ListFinancialRecordsRequest struct {
	Type     *model.RecordType `query:"type"`
	Category string            `query:"category"`
	DateFrom *string           `query:"dateFrom"`
	DateTo   *string           `query:"dateTo"`
}

func (r *ListFinancialRecordsRequest) Validate() error {
	return nil
}

type GetDashboardSummaryRequest struct {
	DateFrom      *string              `query:"dateFrom"`
	DateTo        *string              `query:"dateTo"`
	TrendInterval *model.TrendInterval `query:"trendInterval" validate:"omitempty,oneof=weekly monthly"`
	TrendPeriods  *int                 `query:"trendPeriods" validate:"omitempty,min=1,max=24"`
	RecentLimit   *int                 `query:"recentLimit" validate:"omitempty,min=1,max=20"`
}

func (r *GetDashboardSummaryRequest) Validate() error {
	return validator.New().Struct(r)
}

func NewFinancialRecordHandler(s *server.Server, recordService *service.FinancialRecordService) *FinancialRecordHandler {
	return &FinancialRecordHandler{
		Handler:       NewHandler(s),
		recordService: recordService,
	}
}

func (h *FinancialRecordHandler) CreateFinancialRecord(c echo.Context, req *CreateFinancialRecordRequest) (*model.FinancialRecord, error) {
	currentUser := middleware.GetCurrentUser(c)

	return h.recordService.CreateFinancialRecord(c.Request().Context(), service.CreateFinancialRecordInput{
		UserID:   currentUser.ID,
		Amount:   req.Amount,
		Type:     req.Type,
		Category: req.Category,
		Date:     req.Date,
		Notes:    req.Notes,
	})
}

func (h *FinancialRecordHandler) ListFinancialRecords(c echo.Context, req *ListFinancialRecordsRequest) ([]model.FinancialRecord, error) {
	return h.recordService.ListFinancialRecords(c.Request().Context(), service.ListFinancialRecordFiltersInput{
		Type:     req.Type,
		Category: req.Category,
		DateFrom: req.DateFrom,
		DateTo:   req.DateTo,
	})
}

func (h *FinancialRecordHandler) GetFinancialRecord(c echo.Context, req *GetFinancialRecordRequest) (*model.FinancialRecord, error) {
	recordID, err := uuid.Parse(req.ID)
	if err != nil {
		return nil, err
	}

	return h.recordService.GetFinancialRecordByID(c.Request().Context(), recordID)
}

func (h *FinancialRecordHandler) GetDashboardSummary(c echo.Context, req *GetDashboardSummaryRequest) (*model.DashboardSummary, error) {
	currentUser := middleware.GetCurrentUser(c)

	return h.recordService.GetDashboardSummary(c.Request().Context(), service.GetDashboardSummaryInput{
		UserID:        currentUser.ID,
		DateFrom:      req.DateFrom,
		DateTo:        req.DateTo,
		TrendInterval: req.TrendInterval,
		TrendPeriods:  req.TrendPeriods,
		RecentLimit:   req.RecentLimit,
	})
}

func (h *FinancialRecordHandler) UpdateFinancialRecord(c echo.Context, req *UpdateFinancialRecordRequest) (*model.FinancialRecord, error) {
	recordID, err := uuid.Parse(req.ID)
	if err != nil {
		return nil, err
	}

	return h.recordService.UpdateFinancialRecord(c.Request().Context(), recordID, service.UpdateFinancialRecordInput{
		Amount:   req.Amount,
		Type:     req.Type,
		Category: req.Category,
		Date:     req.Date,
		Notes:    req.Notes,
	})
}

func (h *FinancialRecordHandler) DeleteFinancialRecord(c echo.Context, req *GetFinancialRecordRequest) error {
	recordID, err := uuid.Parse(req.ID)
	if err != nil {
		return err
	}

	return h.recordService.DeleteFinancialRecord(c.Request().Context(), recordID)
}

func (h *FinancialRecordHandler) RegisterRoutes(group *echo.Group, auth *middleware.AuthMiddleware) {
	readGroup := group.Group("/records", auth.RequireAuth, auth.RequireActiveUser, auth.RequireMinimumRole(model.UserRoleViewer))
	readGroup.GET("", Handle(h.Handler, h.ListFinancialRecords, http.StatusOK, &ListFinancialRecordsRequest{}))
	readGroup.GET("/:id", Handle(h.Handler, h.GetFinancialRecord, http.StatusOK, &GetFinancialRecordRequest{}))

	dashboardGroup := group.Group("/dashboard", auth.RequireAuth, auth.RequireActiveUser, auth.RequireMinimumRole(model.UserRoleViewer))
	dashboardGroup.GET("/summary", Handle(h.Handler, h.GetDashboardSummary, http.StatusOK, &GetDashboardSummaryRequest{}))

	writeGroup := group.Group("/records", auth.RequireAuth, auth.RequireActiveUser, auth.RequireMinimumRole(model.UserRoleAnalyst))
	writeGroup.POST("", Handle(h.Handler, h.CreateFinancialRecord, http.StatusCreated, &CreateFinancialRecordRequest{}))
	writeGroup.PATCH("/:id", Handle(h.Handler, h.UpdateFinancialRecord, http.StatusOK, &UpdateFinancialRecordRequest{}))

	adminGroup := group.Group("/records", auth.RequireAuth, auth.RequireActiveUser, auth.RequireRoles(model.UserRoleAdmin))
	adminGroup.DELETE("/:id", HandleNoContent(h.Handler, h.DeleteFinancialRecord, http.StatusNoContent, &GetFinancialRecordRequest{}))
}
