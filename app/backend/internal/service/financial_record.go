package service

import (
	"context"
	"strings"
	"time"

	"github.com/apk471/go-boilerplate/internal/errs"
	"github.com/apk471/go-boilerplate/internal/model"
	"github.com/apk471/go-boilerplate/internal/repository"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

const financialRecordDateLayout = "2006-01-02"

const (
	defaultDashboardRecentLimit  = 5
	defaultDashboardTrendPeriods = 6
	maxDashboardRecentLimit      = 20
	maxDashboardTrendPeriods     = 24
)

var currentTime = time.Now

type financialRecordRepository interface {
	CreateFinancialRecord(ctx context.Context, params repository.CreateFinancialRecordParams) (*model.FinancialRecord, error)
	ListFinancialRecords(ctx context.Context, filters model.FinancialRecordFilters) ([]model.FinancialRecord, error)
	GetFinancialRecordByID(ctx context.Context, id uuid.UUID) (*model.FinancialRecord, error)
	UpdateFinancialRecord(ctx context.Context, id uuid.UUID, params repository.UpdateFinancialRecordParams) (*model.FinancialRecord, error)
	DeleteFinancialRecord(ctx context.Context, id uuid.UUID) error
	GetDashboardSummary(ctx context.Context, params repository.DashboardSummaryParams) (*model.DashboardSummary, error)
}

type FinancialRecordService struct {
	recordRepo financialRecordRepository
}

type CreateFinancialRecordInput struct {
	UserID   uuid.UUID
	Amount   string
	Type     model.RecordType
	Category string
	Date     string
	Notes    *string
}

type UpdateFinancialRecordInput struct {
	Amount   *string
	Type     *model.RecordType
	Category *string
	Date     *string
	Notes    *string
}

type ListFinancialRecordFiltersInput struct {
	Type     *model.RecordType
	Category string
	DateFrom *string
	DateTo   *string
}

type GetDashboardSummaryInput struct {
	UserID        uuid.UUID
	DateFrom      *string
	DateTo        *string
	TrendInterval *model.TrendInterval
	TrendPeriods  *int
	RecentLimit   *int
}

func NewFinancialRecordService(recordRepo financialRecordRepository) *FinancialRecordService {
	return &FinancialRecordService{recordRepo: recordRepo}
}

func (s *FinancialRecordService) CreateFinancialRecord(ctx context.Context, input CreateFinancialRecordInput) (*model.FinancialRecord, error) {
	amount, err := normalizeAmount(input.Amount)
	if err != nil {
		return nil, err
	}

	recordType, err := normalizeRecordType(input.Type)
	if err != nil {
		return nil, err
	}

	recordDate, err := parseRecordDate(input.Date, "date")
	if err != nil {
		return nil, err
	}

	category := strings.TrimSpace(input.Category)
	if category == "" {
		return nil, errs.NewBadRequestError("Validation failed", true, nil, []errs.FieldError{
			{Field: "category", Error: "is required"},
		}, nil)
	}

	return s.recordRepo.CreateFinancialRecord(ctx, repository.CreateFinancialRecordParams{
		UserID:     input.UserID,
		Amount:     amount,
		Type:       recordType,
		Category:   category,
		RecordDate: recordDate,
		Notes:      normalizeOptionalString(input.Notes),
	})
}

func (s *FinancialRecordService) ListFinancialRecords(ctx context.Context, input ListFinancialRecordFiltersInput) ([]model.FinancialRecord, error) {
	filters := model.FinancialRecordFilters{
		Category: strings.TrimSpace(input.Category),
	}

	if input.Type != nil {
		recordType, err := normalizeRecordType(*input.Type)
		if err != nil {
			return nil, err
		}
		filters.Type = &recordType
	}

	if input.DateFrom != nil && strings.TrimSpace(*input.DateFrom) != "" {
		dateFrom, err := parseRecordDate(*input.DateFrom, "dateFrom")
		if err != nil {
			return nil, err
		}
		filters.DateFrom = &dateFrom
	}

	if input.DateTo != nil && strings.TrimSpace(*input.DateTo) != "" {
		dateTo, err := parseRecordDate(*input.DateTo, "dateTo")
		if err != nil {
			return nil, err
		}
		filters.DateTo = &dateTo
	}

	return s.recordRepo.ListFinancialRecords(ctx, filters)
}

func (s *FinancialRecordService) GetFinancialRecordByID(ctx context.Context, id uuid.UUID) (*model.FinancialRecord, error) {
	return s.recordRepo.GetFinancialRecordByID(ctx, id)
}

func (s *FinancialRecordService) GetDashboardSummary(ctx context.Context, input GetDashboardSummaryInput) (*model.DashboardSummary, error) {
	var (
		dateFrom *time.Time
		dateTo   *time.Time
		err      error
	)

	if input.DateFrom != nil && strings.TrimSpace(*input.DateFrom) != "" {
		parsedDateFrom, parseErr := parseRecordDate(*input.DateFrom, "dateFrom")
		if parseErr != nil {
			return nil, parseErr
		}
		dateFrom = &parsedDateFrom
	}

	if input.DateTo != nil && strings.TrimSpace(*input.DateTo) != "" {
		parsedDateTo, parseErr := parseRecordDate(*input.DateTo, "dateTo")
		if parseErr != nil {
			return nil, parseErr
		}
		dateTo = &parsedDateTo
	}

	if dateFrom != nil && dateTo != nil && dateFrom.After(*dateTo) {
		return nil, errs.NewBadRequestError("Validation failed", true, nil, []errs.FieldError{
			{Field: "dateFrom", Error: "must be before or equal to dateTo"},
		}, nil)
	}

	trendInterval := model.TrendIntervalMonthly
	if input.TrendInterval != nil && strings.TrimSpace(string(*input.TrendInterval)) != "" {
		trendInterval, err = normalizeTrendInterval(*input.TrendInterval)
		if err != nil {
			return nil, err
		}
	}

	trendPeriods := defaultDashboardTrendPeriods
	if input.TrendPeriods != nil {
		if *input.TrendPeriods < 1 || *input.TrendPeriods > maxDashboardTrendPeriods {
			return nil, errs.NewBadRequestError("Validation failed", true, nil, []errs.FieldError{
				{Field: "trendPeriods", Error: "must be between 1 and 24"},
			}, nil)
		}
		trendPeriods = *input.TrendPeriods
	}

	recentLimit := defaultDashboardRecentLimit
	if input.RecentLimit != nil {
		if *input.RecentLimit < 1 || *input.RecentLimit > maxDashboardRecentLimit {
			return nil, errs.NewBadRequestError("Validation failed", true, nil, []errs.FieldError{
				{Field: "recentLimit", Error: "must be between 1 and 20"},
			}, nil)
		}
		recentLimit = *input.RecentLimit
	}

	trendStart, trendEnd := deriveTrendRange(dateFrom, dateTo, trendInterval, trendPeriods)

	return s.recordRepo.GetDashboardSummary(ctx, repository.DashboardSummaryParams{
		UserID:        input.UserID,
		DateFrom:      dateFrom,
		DateTo:        dateTo,
		TrendInterval: trendInterval,
		TrendStart:    trendStart,
		TrendEnd:      trendEnd,
		RecentLimit:   recentLimit,
	})
}

func (s *FinancialRecordService) UpdateFinancialRecord(ctx context.Context, id uuid.UUID, input UpdateFinancialRecordInput) (*model.FinancialRecord, error) {
	existing, err := s.recordRepo.GetFinancialRecordByID(ctx, id)
	if err != nil {
		return nil, err
	}

	amount := existing.Amount
	recordType := existing.Type
	category := existing.Category
	recordDate := existing.RecordDate
	notes := existing.Notes

	if input.Amount != nil {
		amount, err = normalizeAmount(*input.Amount)
		if err != nil {
			return nil, err
		}
	}
	if input.Type != nil {
		recordType, err = normalizeRecordType(*input.Type)
		if err != nil {
			return nil, err
		}
	}
	if input.Category != nil {
		category = strings.TrimSpace(*input.Category)
		if category == "" {
			return nil, errs.NewBadRequestError("Validation failed", true, nil, []errs.FieldError{
				{Field: "category", Error: "is required"},
			}, nil)
		}
	}
	if input.Date != nil && strings.TrimSpace(*input.Date) != "" {
		recordDate, err = parseRecordDate(*input.Date, "date")
		if err != nil {
			return nil, err
		}
	}
	if input.Notes != nil {
		notes = normalizeOptionalString(input.Notes)
	}

	return s.recordRepo.UpdateFinancialRecord(ctx, id, repository.UpdateFinancialRecordParams{
		Amount:     amount,
		Type:       recordType,
		Category:   category,
		RecordDate: recordDate,
		Notes:      notes,
	})
}

func (s *FinancialRecordService) DeleteFinancialRecord(ctx context.Context, id uuid.UUID) error {
	if _, err := s.recordRepo.GetFinancialRecordByID(ctx, id); err != nil {
		return err
	}

	return s.recordRepo.DeleteFinancialRecord(ctx, id)
}

func normalizeAmount(raw string) (string, error) {
	parsed, err := decimal.NewFromString(strings.TrimSpace(raw))
	if err != nil {
		return "", errs.NewBadRequestError("Validation failed", true, nil, []errs.FieldError{
			{Field: "amount", Error: "must be a valid decimal amount"},
		}, nil)
	}
	if parsed.IsNegative() {
		return "", errs.NewBadRequestError("Validation failed", true, nil, []errs.FieldError{
			{Field: "amount", Error: "must be greater than or equal to 0"},
		}, nil)
	}

	return parsed.StringFixed(2), nil
}

func normalizeRecordType(recordType model.RecordType) (model.RecordType, error) {
	normalized := model.NormalizeRecordType(recordType)
	if !model.IsValidRecordType(normalized) {
		return "", errs.NewBadRequestError("Validation failed", true, nil, []errs.FieldError{
			{Field: "type", Error: "must be one of: income expense"},
		}, nil)
	}

	return normalized, nil
}

func parseRecordDate(value string, field string) (time.Time, error) {
	recordDate, err := time.Parse(financialRecordDateLayout, strings.TrimSpace(value))
	if err != nil {
		return time.Time{}, errs.NewBadRequestError("Validation failed", true, nil, []errs.FieldError{
			{Field: field, Error: "must be in YYYY-MM-DD format"},
		}, nil)
	}

	return recordDate.UTC(), nil
}

func normalizeOptionalString(value *string) *string {
	if value == nil {
		return nil
	}

	trimmed := strings.TrimSpace(*value)
	if trimmed == "" {
		return nil
	}

	return &trimmed
}

func normalizeTrendInterval(interval model.TrendInterval) (model.TrendInterval, error) {
	normalized := model.NormalizeTrendInterval(interval)
	if !model.IsValidTrendInterval(normalized) {
		return "", errs.NewBadRequestError("Validation failed", true, nil, []errs.FieldError{
			{Field: "trendInterval", Error: "must be one of: weekly monthly"},
		}, nil)
	}

	return normalized, nil
}

func deriveTrendRange(dateFrom *time.Time, dateTo *time.Time, interval model.TrendInterval, periods int) (time.Time, time.Time) {
	if dateFrom != nil && dateTo != nil {
		return *dateFrom, *dateTo
	}

	now := currentTime().UTC()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

	end := today
	if dateTo != nil {
		end = *dateTo
	}

	if dateFrom != nil {
		return *dateFrom, end
	}

	switch interval {
	case model.TrendIntervalWeekly:
		start := end.AddDate(0, 0, -7*(periods-1))
		return start, end
	default:
		start := end.AddDate(0, -(periods - 1), 0)
		return start, end
	}
}
