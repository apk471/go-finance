package model

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

type RecordType string
type TrendInterval string

const (
	RecordTypeIncome  RecordType = "income"
	RecordTypeExpense RecordType = "expense"

	TrendIntervalWeekly  TrendInterval = "weekly"
	TrendIntervalMonthly TrendInterval = "monthly"
)

type FinancialRecord struct {
	Base
	UserID     uuid.UUID  `json:"userId" db:"user_id"`
	Amount     string     `json:"amount" db:"amount"`
	Type       RecordType `json:"type" db:"type"`
	Category   string     `json:"category" db:"category"`
	RecordDate time.Time  `json:"date" db:"record_date"`
	Notes      *string    `json:"notes,omitempty" db:"notes"`
}

type FinancialRecordFilters struct {
	Type     *RecordType
	Category string
	DateFrom *time.Time
	DateTo   *time.Time
}

type DashboardCategoryTotal struct {
	Type     RecordType `json:"type"`
	Category string     `json:"category"`
	Total    string     `json:"total"`
}

type DashboardTrendPoint struct {
	Interval    TrendInterval `json:"interval"`
	PeriodStart time.Time     `json:"periodStart"`
	PeriodEnd   time.Time     `json:"periodEnd"`
	Income      string        `json:"income"`
	Expenses    string        `json:"expenses"`
	NetBalance  string        `json:"netBalance"`
}

type DashboardSummary struct {
	DateFrom       *time.Time               `json:"dateFrom,omitempty"`
	DateTo         *time.Time               `json:"dateTo,omitempty"`
	TrendInterval  TrendInterval            `json:"trendInterval"`
	TotalIncome    string                   `json:"totalIncome"`
	TotalExpenses  string                   `json:"totalExpenses"`
	NetBalance     string                   `json:"netBalance"`
	CategoryTotals []DashboardCategoryTotal `json:"categoryTotals"`
	RecentActivity []FinancialRecord        `json:"recentActivity"`
	Trends         []DashboardTrendPoint    `json:"trends"`
}

func NormalizeRecordType(recordType RecordType) RecordType {
	return RecordType(strings.ToLower(string(recordType)))
}

func IsValidRecordType(recordType RecordType) bool {
	switch NormalizeRecordType(recordType) {
	case RecordTypeIncome, RecordTypeExpense:
		return true
	default:
		return false
	}
}

func NormalizeTrendInterval(interval TrendInterval) TrendInterval {
	return TrendInterval(strings.ToLower(string(interval)))
}

func IsValidTrendInterval(interval TrendInterval) bool {
	switch NormalizeTrendInterval(interval) {
	case TrendIntervalWeekly, TrendIntervalMonthly:
		return true
	default:
		return false
	}
}
