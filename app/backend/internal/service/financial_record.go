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

type financialRecordRepository interface {
	CreateFinancialRecord(ctx context.Context, params repository.CreateFinancialRecordParams) (*model.FinancialRecord, error)
	ListFinancialRecords(ctx context.Context, filters model.FinancialRecordFilters) ([]model.FinancialRecord, error)
	GetFinancialRecordByID(ctx context.Context, id uuid.UUID) (*model.FinancialRecord, error)
	UpdateFinancialRecord(ctx context.Context, id uuid.UUID, params repository.UpdateFinancialRecordParams) (*model.FinancialRecord, error)
	DeleteFinancialRecord(ctx context.Context, id uuid.UUID) error
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
