package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/apk471/go-boilerplate/internal/model"
	"github.com/apk471/go-boilerplate/internal/server"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
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

type DashboardSummaryParams struct {
	UserID        uuid.UUID
	DateFrom      *time.Time
	DateTo        *time.Time
	TrendInterval model.TrendInterval
	TrendStart    time.Time
	TrendEnd      time.Time
	RecentLimit   int
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

func (r *FinancialRecordRepository) GetDashboardSummary(ctx context.Context, params DashboardSummaryParams) (*model.DashboardSummary, error) {
	summary := &model.DashboardSummary{
		DateFrom:      params.DateFrom,
		DateTo:        params.DateTo,
		TrendInterval: params.TrendInterval,
	}

	totalIncome, totalExpenses, err := r.getDashboardTotals(ctx, params)
	if err != nil {
		return nil, err
	}

	summary.TotalIncome = totalIncome
	summary.TotalExpenses = totalExpenses
	summary.NetBalance = calculateNetBalance(totalIncome, totalExpenses)

	categoryTotals, err := r.getDashboardCategoryTotals(ctx, params)
	if err != nil {
		return nil, err
	}
	summary.CategoryTotals = categoryTotals

	recentActivity, err := r.getDashboardRecentActivity(ctx, params)
	if err != nil {
		return nil, err
	}
	summary.RecentActivity = recentActivity

	trends, err := r.getDashboardTrends(ctx, params)
	if err != nil {
		return nil, err
	}
	summary.Trends = trends

	return summary, nil
}

func (r *FinancialRecordRepository) getDashboardTotals(ctx context.Context, params DashboardSummaryParams) (string, string, error) {
	const query = `
		SELECT
			COALESCE(SUM(CASE WHEN type = 'income' THEN amount ELSE 0 END), 0)::text AS total_income,
			COALESCE(SUM(CASE WHEN type = 'expense' THEN amount ELSE 0 END), 0)::text AS total_expenses
		FROM financial_records
		WHERE user_id = $1
		  AND ($2::date IS NULL OR record_date >= $2::date)
		  AND ($3::date IS NULL OR record_date <= $3::date)
	`

	var totalIncome string
	var totalExpenses string
	if err := r.server.DB.Pool.QueryRow(
		ctx,
		query,
		params.UserID,
		params.DateFrom,
		params.DateTo,
	).Scan(&totalIncome, &totalExpenses); err != nil {
		return "", "", fmt.Errorf("get dashboard totals: %w", err)
	}

	return totalIncome, totalExpenses, nil
}

func (r *FinancialRecordRepository) getDashboardCategoryTotals(ctx context.Context, params DashboardSummaryParams) ([]model.DashboardCategoryTotal, error) {
	const query = `
		SELECT type, category, COALESCE(SUM(amount), 0)::text AS total
		FROM financial_records
		WHERE user_id = $1
		  AND ($2::date IS NULL OR record_date >= $2::date)
		  AND ($3::date IS NULL OR record_date <= $3::date)
		GROUP BY type, category
		ORDER BY SUM(amount) DESC, category ASC
	`

	rows, err := r.server.DB.Pool.Query(ctx, query, params.UserID, params.DateFrom, params.DateTo)
	if err != nil {
		return nil, fmt.Errorf("get dashboard category totals: %w", err)
	}
	defer rows.Close()

	categoryTotals := make([]model.DashboardCategoryTotal, 0)
	for rows.Next() {
		var total model.DashboardCategoryTotal
		if err := rows.Scan(&total.Type, &total.Category, &total.Total); err != nil {
			return nil, fmt.Errorf("scan dashboard category total: %w", err)
		}
		categoryTotals = append(categoryTotals, total)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate dashboard category totals: %w", err)
	}

	return categoryTotals, nil
}

func (r *FinancialRecordRepository) getDashboardRecentActivity(ctx context.Context, params DashboardSummaryParams) ([]model.FinancialRecord, error) {
	const query = `
		SELECT id, user_id, amount::text, type, category, record_date, notes, created_at, updated_at
		FROM financial_records
		WHERE user_id = $1
		  AND ($2::date IS NULL OR record_date >= $2::date)
		  AND ($3::date IS NULL OR record_date <= $3::date)
		ORDER BY record_date DESC, created_at DESC
		LIMIT $4
	`

	rows, err := r.server.DB.Pool.Query(ctx, query, params.UserID, params.DateFrom, params.DateTo, params.RecentLimit)
	if err != nil {
		return nil, fmt.Errorf("get dashboard recent activity: %w", err)
	}
	defer rows.Close()

	records := make([]model.FinancialRecord, 0, params.RecentLimit)
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
			return nil, fmt.Errorf("scan dashboard recent activity: %w", err)
		}
		records = append(records, record)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate dashboard recent activity: %w", err)
	}

	return records, nil
}

func (r *FinancialRecordRepository) getDashboardTrends(ctx context.Context, params DashboardSummaryParams) ([]model.DashboardTrendPoint, error) {
	query := weeklyDashboardTrendsQuery
	if params.TrendInterval == model.TrendIntervalMonthly {
		query = monthlyDashboardTrendsQuery
	}

	rows, err := r.server.DB.Pool.Query(ctx, query, params.UserID, params.TrendStart, params.TrendEnd)
	if err != nil {
		return nil, fmt.Errorf("get dashboard trends: %w", err)
	}
	defer rows.Close()

	trends := make([]model.DashboardTrendPoint, 0)
	for rows.Next() {
		var point model.DashboardTrendPoint
		point.Interval = params.TrendInterval
		if err := rows.Scan(&point.PeriodStart, &point.PeriodEnd, &point.Income, &point.Expenses); err != nil {
			return nil, fmt.Errorf("scan dashboard trend: %w", err)
		}
		point.NetBalance = calculateNetBalance(point.Income, point.Expenses)
		trends = append(trends, point)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate dashboard trends: %w", err)
	}

	return trends, nil
}

func calculateNetBalance(totalIncome string, totalExpenses string) string {
	income, err := parseDecimal(totalIncome)
	if err != nil {
		return "0.00"
	}

	expenses, err := parseDecimal(totalExpenses)
	if err != nil {
		return income.StringFixed(2)
	}

	return income.Sub(expenses).StringFixed(2)
}

func parseDecimal(value string) (decimal.Decimal, error) {
	return decimal.NewFromString(strings.TrimSpace(value))
}

const weeklyDashboardTrendsQuery = `
	WITH periods AS (
		SELECT
			generate_series(
				date_trunc('week', $2::date),
				date_trunc('week', $3::date),
				interval '1 week'
			)::date AS period_start
	),
	aggregated AS (
		SELECT
			date_trunc('week', record_date)::date AS period_start,
			COALESCE(SUM(CASE WHEN type = 'income' THEN amount ELSE 0 END), 0)::text AS income,
			COALESCE(SUM(CASE WHEN type = 'expense' THEN amount ELSE 0 END), 0)::text AS expenses
		FROM financial_records
		WHERE user_id = $1
		  AND record_date >= $2::date
		  AND record_date <= $3::date
		GROUP BY 1
	)
	SELECT
		periods.period_start,
		(periods.period_start + interval '6 days')::date AS period_end,
		COALESCE(aggregated.income, '0') AS income,
		COALESCE(aggregated.expenses, '0') AS expenses
	FROM periods
	LEFT JOIN aggregated ON aggregated.period_start = periods.period_start
	ORDER BY periods.period_start ASC
`

const monthlyDashboardTrendsQuery = `
	WITH periods AS (
		SELECT
			generate_series(
				date_trunc('month', $2::date),
				date_trunc('month', $3::date),
				interval '1 month'
			)::date AS period_start
	),
	aggregated AS (
		SELECT
			date_trunc('month', record_date)::date AS period_start,
			COALESCE(SUM(CASE WHEN type = 'income' THEN amount ELSE 0 END), 0)::text AS income,
			COALESCE(SUM(CASE WHEN type = 'expense' THEN amount ELSE 0 END), 0)::text AS expenses
		FROM financial_records
		WHERE user_id = $1
		  AND record_date >= $2::date
		  AND record_date <= $3::date
		GROUP BY 1
	)
	SELECT
		periods.period_start,
		(date_trunc('month', periods.period_start) + interval '1 month - 1 day')::date AS period_end,
		COALESCE(aggregated.income, '0') AS income,
		COALESCE(aggregated.expenses, '0') AS expenses
	FROM periods
	LEFT JOIN aggregated ON aggregated.period_start = periods.period_start
	ORDER BY periods.period_start ASC
`
