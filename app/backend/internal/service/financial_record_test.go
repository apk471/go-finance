package service

import (
	"context"
	"testing"
	"time"

	"github.com/apk471/go-boilerplate/internal/model"
	"github.com/apk471/go-boilerplate/internal/repository"
	"github.com/google/uuid"
)

type stubFinancialRecordRepository struct {
	createInput  repository.CreateFinancialRecordParams
	createResult *model.FinancialRecord
	createErr    error
	listInput    model.FinancialRecordFilters
	listResult   []model.FinancialRecord
	listErr      error
	getResult    *model.FinancialRecord
	getErr       error
	updateInput  repository.UpdateFinancialRecordParams
	updateResult *model.FinancialRecord
	updateErr    error
	deleteID     uuid.UUID
	deleteErr    error
}

func (s *stubFinancialRecordRepository) CreateFinancialRecord(ctx context.Context, params repository.CreateFinancialRecordParams) (*model.FinancialRecord, error) {
	s.createInput = params
	if s.createResult != nil {
		return s.createResult, s.createErr
	}

	return &model.FinancialRecord{
		Base:       model.Base{BaseWithId: model.BaseWithId{ID: uuid.New()}},
		UserID:     params.UserID,
		Amount:     params.Amount,
		Type:       params.Type,
		Category:   params.Category,
		RecordDate: params.RecordDate,
		Notes:      params.Notes,
	}, s.createErr
}

func (s *stubFinancialRecordRepository) ListFinancialRecords(ctx context.Context, filters model.FinancialRecordFilters) ([]model.FinancialRecord, error) {
	s.listInput = filters
	return s.listResult, s.listErr
}

func (s *stubFinancialRecordRepository) GetFinancialRecordByID(ctx context.Context, id uuid.UUID) (*model.FinancialRecord, error) {
	return s.getResult, s.getErr
}

func (s *stubFinancialRecordRepository) UpdateFinancialRecord(ctx context.Context, id uuid.UUID, params repository.UpdateFinancialRecordParams) (*model.FinancialRecord, error) {
	s.updateInput = params
	if s.updateResult != nil {
		return s.updateResult, s.updateErr
	}

	return &model.FinancialRecord{
		Base:       model.Base{BaseWithId: model.BaseWithId{ID: id}},
		Amount:     params.Amount,
		Type:       params.Type,
		Category:   params.Category,
		RecordDate: params.RecordDate,
		Notes:      params.Notes,
	}, s.updateErr
}

func (s *stubFinancialRecordRepository) DeleteFinancialRecord(ctx context.Context, id uuid.UUID) error {
	s.deleteID = id
	return s.deleteErr
}

func TestCreateFinancialRecordNormalizesInput(t *testing.T) {
	repo := &stubFinancialRecordRepository{}
	svc := NewFinancialRecordService(repo)
	userID := uuid.New()
	notes := "  Monthly salary  "

	record, err := svc.CreateFinancialRecord(context.Background(), CreateFinancialRecordInput{
		UserID:   userID,
		Amount:   "1250.5",
		Type:     "INCOME",
		Category: " Salary ",
		Date:     "2026-04-02",
		Notes:    &notes,
	})
	if err != nil {
		t.Fatalf("CreateFinancialRecord() error = %v", err)
	}

	if repo.createInput.Amount != "1250.50" {
		t.Fatalf("expected normalized amount, got %s", repo.createInput.Amount)
	}
	if repo.createInput.Type != model.RecordTypeIncome {
		t.Fatalf("expected income type, got %s", repo.createInput.Type)
	}
	if repo.createInput.Category != "Salary" {
		t.Fatalf("expected trimmed category, got %q", repo.createInput.Category)
	}
	if record.UserID != userID {
		t.Fatalf("expected user id %s, got %s", userID, record.UserID)
	}
}

func TestListFinancialRecordsNormalizesFilters(t *testing.T) {
	repo := &stubFinancialRecordRepository{}
	svc := NewFinancialRecordService(repo)
	recordType := model.RecordType("EXPENSE")
	dateFrom := "2026-04-01"
	dateTo := "2026-04-30"

	_, err := svc.ListFinancialRecords(context.Background(), ListFinancialRecordFiltersInput{
		Type:     &recordType,
		Category: " Food ",
		DateFrom: &dateFrom,
		DateTo:   &dateTo,
	})
	if err != nil {
		t.Fatalf("ListFinancialRecords() error = %v", err)
	}

	if repo.listInput.Type == nil || *repo.listInput.Type != model.RecordTypeExpense {
		t.Fatalf("expected expense filter, got %#v", repo.listInput.Type)
	}
	if repo.listInput.Category != "Food" {
		t.Fatalf("expected trimmed category, got %q", repo.listInput.Category)
	}
	if repo.listInput.DateFrom == nil || repo.listInput.DateTo == nil {
		t.Fatal("expected date filters to be set")
	}
}

func TestUpdateFinancialRecordPreservesExistingFields(t *testing.T) {
	repo := &stubFinancialRecordRepository{
		getResult: &model.FinancialRecord{
			Base:       model.Base{BaseWithId: model.BaseWithId{ID: uuid.New()}},
			Amount:     "50.00",
			Type:       model.RecordTypeExpense,
			Category:   "Transport",
			RecordDate: time.Date(2026, 4, 1, 0, 0, 0, 0, time.UTC),
		},
	}
	svc := NewFinancialRecordService(repo)
	category := " Travel "

	_, err := svc.UpdateFinancialRecord(context.Background(), repo.getResult.ID, UpdateFinancialRecordInput{
		Category: &category,
	})
	if err != nil {
		t.Fatalf("UpdateFinancialRecord() error = %v", err)
	}

	if repo.updateInput.Amount != "50.00" {
		t.Fatalf("expected amount to be preserved, got %s", repo.updateInput.Amount)
	}
	if repo.updateInput.Category != "Travel" {
		t.Fatalf("expected category to be updated, got %q", repo.updateInput.Category)
	}
}
