package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/apk471/go-boilerplate/internal/model"
	"github.com/apk471/go-boilerplate/internal/server"
	"github.com/google/uuid"
)

type FinancialRecordRepository struct {
	server *server.Server
}

type CreateFinancialRecordParams struct {
	UserID     uuid.UUID
	Amount     string
	Type       model.RecordType
	Category   string
	RecordDate time.Time
	Notes      *string
}

type UpdateFinancialRecordParams struct {
	Amount     string
	Type       model.RecordType
	Category   string
	RecordDate time.Time
	Notes      *string
}

func NewFinancialRecordRepository(s *server.Server) *FinancialRecordRepository {
	return &FinancialRecordRepository{server: s}
}

func (r *FinancialRecordRepository) CreateFinancialRecord(ctx context.Context, params CreateFinancialRecordParams) (*model.FinancialRecord, error) {
	const query = `
		INSERT INTO financial_records (user_id, amount, type, category, record_date, notes)
		VALUES ($1, $2::numeric(14,2), $3, $4, $5, $6)
		RETURNING id, user_id, amount::text, type, category, record_date, notes, created_at, updated_at
	`

	record := &model.FinancialRecord{}
	if err := r.server.DB.Pool.QueryRow(
		ctx,
		query,
		params.UserID,
		params.Amount,
		params.Type,
		params.Category,
		params.RecordDate,
		params.Notes,
	).Scan(
		&record.ID,
		&record.UserID,
		&record.Amount,
		&record.Type,
		&record.Category,
		&record.RecordDate,
		&record.Notes,
		&record.CreatedAt,
		&record.UpdatedAt,
	); err != nil {
		return nil, fmt.Errorf("create financial record: %w", err)
	}

	return record, nil
}

func (r *FinancialRecordRepository) ListFinancialRecords(ctx context.Context, filters model.FinancialRecordFilters) ([]model.FinancialRecord, error) {
	var query strings.Builder
	query.WriteString(`
		SELECT id, user_id, amount::text, type, category, record_date, notes, created_at, updated_at
		FROM financial_records
		WHERE 1 = 1
	`)

	args := make([]any, 0, 4)
	argIndex := 1

	if filters.Type != nil {
		query.WriteString(fmt.Sprintf(" AND type = $%d", argIndex))
		args = append(args, *filters.Type)
		argIndex++
	}
	if filters.Category != "" {
		query.WriteString(fmt.Sprintf(" AND category = $%d", argIndex))
		args = append(args, filters.Category)
		argIndex++
	}
	if filters.DateFrom != nil {
		query.WriteString(fmt.Sprintf(" AND record_date >= $%d", argIndex))
		args = append(args, *filters.DateFrom)
		argIndex++
	}
	if filters.DateTo != nil {
		query.WriteString(fmt.Sprintf(" AND record_date <= $%d", argIndex))
		args = append(args, *filters.DateTo)
		argIndex++
	}

	query.WriteString(" ORDER BY record_date DESC, created_at DESC")

	rows, err := r.server.DB.Pool.Query(ctx, query.String(), args...)
	if err != nil {
		return nil, fmt.Errorf("list financial records: %w", err)
	}
	defer rows.Close()

	records := make([]model.FinancialRecord, 0)
	for rows.Next() {
		var record model.FinancialRecord
		if err := rows.Scan(
			&record.ID,
			&record.UserID,
			&record.Amount,
			&record.Type,
			&record.Category,
			&record.RecordDate,
			&record.Notes,
			&record.CreatedAt,
			&record.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan financial record: %w", err)
		}
		records = append(records, record)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate financial records: %w", err)
	}

	return records, nil
}

func (r *FinancialRecordRepository) GetFinancialRecordByID(ctx context.Context, id uuid.UUID) (*model.FinancialRecord, error) {
	const query = `
		SELECT id, user_id, amount::text, type, category, record_date, notes, created_at, updated_at
		FROM financial_records
		WHERE id = $1
	`

	record := &model.FinancialRecord{}
	if err := r.server.DB.Pool.QueryRow(ctx, query, id).Scan(
		&record.ID,
		&record.UserID,
		&record.Amount,
		&record.Type,
		&record.Category,
		&record.RecordDate,
		&record.Notes,
		&record.CreatedAt,
		&record.UpdatedAt,
	); err != nil {
		return nil, fmt.Errorf("get financial record by id: %w", err)
	}

	return record, nil
}

func (r *FinancialRecordRepository) UpdateFinancialRecord(ctx context.Context, id uuid.UUID, params UpdateFinancialRecordParams) (*model.FinancialRecord, error) {
	const query = `
		UPDATE financial_records
		SET amount = $2::numeric(14,2),
		    type = $3,
		    category = $4,
		    record_date = $5,
		    notes = $6,
		    updated_at = NOW()
		WHERE id = $1
		RETURNING id, user_id, amount::text, type, category, record_date, notes, created_at, updated_at
	`

	record := &model.FinancialRecord{}
	if err := r.server.DB.Pool.QueryRow(
		ctx,
		query,
		id,
		params.Amount,
		params.Type,
		params.Category,
		params.RecordDate,
		params.Notes,
	).Scan(
		&record.ID,
		&record.UserID,
		&record.Amount,
		&record.Type,
		&record.Category,
		&record.RecordDate,
		&record.Notes,
		&record.CreatedAt,
		&record.UpdatedAt,
	); err != nil {
		return nil, fmt.Errorf("update financial record: %w", err)
	}

	return record, nil
}

func (r *FinancialRecordRepository) DeleteFinancialRecord(ctx context.Context, id uuid.UUID) error {
	const query = `DELETE FROM financial_records WHERE id = $1`

	commandTag, err := r.server.DB.Pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete financial record: %w", err)
	}
	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("delete financial record: no rows in result set")
	}

	return nil
}
