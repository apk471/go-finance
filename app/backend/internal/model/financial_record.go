package model

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

type RecordType string

const (
	RecordTypeIncome  RecordType = "income"
	RecordTypeExpense RecordType = "expense"
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
