package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/patyukin/mbs-credits/internal/model"
	"github.com/rs/zerolog/log"
	"strings"
	"time"
)

func (r *Repository) InsertPaymentSchedules(ctx context.Context, in []model.PaymentSchedule) error {
	currentTime := time.Now().UTC()
	query := `
	INSERT INTO payment_schedules (credit_id, amount, due_date, status, created_at)
	VALUES 
`
	var (
		placeholders []string
		args         []interface{}
	)

	for i, schedule := range in {
		n := i * 5
		placeholders = append(placeholders, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d)", n+1, n+2, n+3, n+4, n+5))
		args = append(
			args,
			schedule.CreditID,
			schedule.Amount,
			schedule.DueDate,
			schedule.Status,
			currentTime,
		)

		query += strings.Join(placeholders, ", ")
	}

	_, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed r.db.ExecContext: %w", err)
	}

	return nil
}

func (r *Repository) SelectPaymentScheduleByUserIDAndCreditID(ctx context.Context, userID, creditID string) ([]model.PaymentSchedule, error) {
	var result []model.PaymentSchedule
	query := `
SELECT id, credit_id, amount, due_date, status, created_at, updated_at
FROM payment_schedules
WHERE credit_id = $1 AND credit_id IN (SELECT id FROM credits WHERE user_id = $2)
ORDER BY due_date
		`
	rows, err := r.db.QueryContext(ctx, query, userID, creditID)
	if err != nil {
		return nil, fmt.Errorf("failed r.db.QueryContext: %w", err)
	}

	for rows.Next() {
		var schedule model.PaymentSchedule
		if err = rows.Scan(
			&schedule.ID,
			&schedule.CreditID,
			&schedule.Amount,
			&schedule.DueDate,
			&schedule.Status,
			&schedule.CreatedAt,
			&schedule.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed rows.Scan: %w", err)
		}

		result = append(result, schedule)
	}

	return result, nil
}

func (r *Repository) SelectPaymentScheduleToPay(ctx context.Context) ([]model.CreditPaymentSchedule, error) {
	var result []model.CreditPaymentSchedule
	query := `
SELECT ps.id, ps.amount, c.account_id
FROM payment_schedules AS ps
INNER JOIN credits AS c ON c.id = ps.credit_id
WHERE status = 'SCHEDULED'
ORDER BY due_date
`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed r.db.QueryContext: %w", err)
	}

	if err = rows.Err(); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return result, nil
		}

		return nil, fmt.Errorf("failed rows.Err(): %w", err)
	}

	defer func(rows *sql.Rows) {
		if err = rows.Close(); err != nil {
			log.Error().Msgf("failed rows.Close: %v", err)
		}
	}(rows)

	for rows.Next() {
		var schedule model.CreditPaymentSchedule
		if err = rows.Scan(
			&schedule.PaymentScheduleID,
			&schedule.Amount,
			&schedule.AccountID,
		); err != nil {
			return nil, fmt.Errorf("failed rows.Scan: %w", err)
		}
		result = append(result, schedule)
	}

	return result, nil
}

func (r *Repository) UpdatePaymentScheduleStatus(ctx context.Context, id, status string) error {
	query := `UPDATE payment_schedules SET status = $1 WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, status, id)
	if err != nil {
		return fmt.Errorf("failed r.db.ExecContext: %w", err)
	}

	return nil
}
