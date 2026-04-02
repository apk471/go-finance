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
	createInput   repository.CreateFinancialRecordParams
	createResult  *model.FinancialRecord
	createErr     error
	listInput     model.FinancialRecordFilters
	listResult    []model.FinancialRecord
	listErr       error
	getResult     *model.FinancialRecord
	getErr        error
	updateInput   repository.UpdateFinancialRecordParams
	updateResult  *model.FinancialRecord
	updateErr     error
	deleteID      uuid.UUID
	deleteErr     error
	summaryInput  repository.DashboardSummaryParams
	summaryResult *model.DashboardSummary
	summaryErr    error
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

func (s *stubFinancialRecordRepository) GetDashboardSummary(ctx context.Context, params repository.DashboardSummaryParams) (*model.DashboardSummary, error) {
	s.summaryInput = params
	if s.summaryResult != nil {
		return s.summaryResult, s.summaryErr
	}

	return &model.DashboardSummary{
		TrendInterval: params.TrendInterval,
		TotalIncome:   "0.00",
		TotalExpenses: "0.00",
		NetBalance:    "0.00",
	}, s.summaryErr
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

func TestGetDashboardSummaryNormalizesFiltersAndDefaults(t *testing.T) {
	repo := &stubFinancialRecordRepository{}
	svc := NewFinancialRecordService(repo)
	userID := uuid.New()
	now := time.Date(2026, 4, 2, 10, 0, 0, 0, time.UTC)
	previousNow := currentTime
	currentTime = func() time.Time { return now }
	defer func() { currentTime = previousNow }()

	trendInterval := model.TrendInterval("WEEKLY")

	_, err := svc.GetDashboardSummary(context.Background(), GetDashboardSummaryInput{
		UserID:        userID,
		TrendInterval: &trendInterval,
	})
	if err != nil {
		t.Fatalf("GetDashboardSummary() error = %v", err)
	}

	if repo.summaryInput.UserID != userID {
		t.Fatalf("expected user id %s, got %s", userID, repo.summaryInput.UserID)
	}
	if repo.summaryInput.TrendInterval != model.TrendIntervalWeekly {
		t.Fatalf("expected weekly interval, got %s", repo.summaryInput.TrendInterval)
	}
	if repo.summaryInput.RecentLimit != defaultDashboardRecentLimit {
		t.Fatalf("expected default recent limit %d, got %d", defaultDashboardRecentLimit, repo.summaryInput.RecentLimit)
	}

	expectedTrendStart := time.Date(2026, 2, 26, 0, 0, 0, 0, time.UTC)
	expectedTrendEnd := time.Date(2026, 4, 2, 0, 0, 0, 0, time.UTC)
	if !repo.summaryInput.TrendStart.Equal(expectedTrendStart) {
		t.Fatalf("expected trend start %s, got %s", expectedTrendStart, repo.summaryInput.TrendStart)
	}
	if !repo.summaryInput.TrendEnd.Equal(expectedTrendEnd) {
		t.Fatalf("expected trend end %s, got %s", expectedTrendEnd, repo.summaryInput.TrendEnd)
	}
}

func TestGetDashboardSummaryUsesProvidedDates(t *testing.T) {
	repo := &stubFinancialRecordRepository{}
	svc := NewFinancialRecordService(repo)
	userID := uuid.New()
	dateFrom := "2026-01-01"
	dateTo := "2026-03-31"
	recentLimit := 10
	trendPeriods := 3

	_, err := svc.GetDashboardSummary(context.Background(), GetDashboardSummaryInput{
		UserID:       userID,
		DateFrom:     &dateFrom,
		DateTo:       &dateTo,
		RecentLimit:  &recentLimit,
		TrendPeriods: &trendPeriods,
	})
	if err != nil {
		t.Fatalf("GetDashboardSummary() error = %v", err)
	}

	if repo.summaryInput.DateFrom == nil || repo.summaryInput.DateTo == nil {
		t.Fatal("expected date filters to be set")
	}
	if repo.summaryInput.RecentLimit != recentLimit {
		t.Fatalf("expected recent limit %d, got %d", recentLimit, repo.summaryInput.RecentLimit)
	}
	if !repo.summaryInput.TrendStart.Equal(*repo.summaryInput.DateFrom) {
		t.Fatalf("expected trend start to match dateFrom, got %s", repo.summaryInput.TrendStart)
	}
	if !repo.summaryInput.TrendEnd.Equal(*repo.summaryInput.DateTo) {
		t.Fatalf("expected trend end to match dateTo, got %s", repo.summaryInput.TrendEnd)
	}
}
