package db

import (
	"context"
	"fmt"
	"github.com/patyukin/mbs-credits/internal/model"
)

func (r *Repository) InsertCredit(ctx context.Context, c model.Credit) (string, error) {
	query := `
INSERT INTO credits (
	account_id,
	credit_application_id,
	user_id,
	amount,
	interest_rate,
	remaining_amount,
	status,
	start_date,
	end_date,
	created_at
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
  RETURNING id`
	row := r.db.QueryRowContext(
		ctx,
		query,
		c.AccountID,
		c.CreditApplicationID,
		c.UserID,
		c.Amount,
		c.InterestRate,
		c.RemainingAmount,
		c.Status,
		c.StartDate,
		c.EndDate,
		c.CreatedAt,
	)

	if row.Err() != nil {
		return "", fmt.Errorf("failed r.db.QueryRowContext: %w", row.Err())
	}

	var id string
	if err := row.Scan(&id); err != nil {
		return "", fmt.Errorf("failed row.Scan: %w", err)
	}

	return id, nil
}

func (r *Repository) SelectCreditByIDAndUserID(ctx context.Context, id, userID string) (model.Credit, error) {
	var c model.Credit
	query := `
SELECT
	id,
	account_id,
	credit_application_id,
	user_id,
	amount,
	interest_rate,
	remaining_amount,
	status,
	start_date,
	end_date,
	created_at,
	updated_at
FROM credits WHERE id = $1 AND user_id = $2`
	row := r.db.QueryRowContext(ctx, query, id, userID)
	if row.Err() != nil {
		return model.Credit{}, fmt.Errorf("failed r.db.QueryRowContext: %w", row.Err())
	}

	if err := row.Scan(
		&c.ID,
		&c.AccountID,
		&c.CreditApplicationID,
		&c.UserID,
		&c.Amount,
		&c.InterestRate,
		&c.RemainingAmount,
		&c.Status,
		&c.StartDate,
		&c.EndDate,
		&c.CreatedAt,
		&c.UpdatedAt,
	); err != nil {
		return model.Credit{}, fmt.Errorf("failed row.Scan: %w", err)
	}

	return c, nil
}

func (r *Repository) SelectCountCreditByUserID(ctx context.Context, userID string) (int32, error) {
	var count int32
	query := `SELECT count(*) FROM credits WHERE user_id = $1`
	row := r.db.QueryRowContext(ctx, query, userID)
	if row.Err() != nil {
		return 0, fmt.Errorf("failed r.db.QueryRowContext: %w", row.Err())
	}

	if err := row.Scan(&count); err != nil {
		return 0, fmt.Errorf("failed row.Scan: %w", err)
	}

	return count, nil
}

func (r *Repository) SelectCreditByUserID(ctx context.Context, userID string, page, limit int32) ([]model.Credit, error) {
	var credits []model.Credit
	query := `
SELECT
	id,
	account_id,
	credit_application_id,
	user_id,
	amount,
	interest_rate,
	remaining_amount,
	status,
	start_date,
	end_date,
	created_at,
	updated_at
FROM credits WHERE user_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`
	rows, err := r.db.QueryContext(ctx, query, userID, limit, page)
	if err != nil {
		return nil, fmt.Errorf("failed r.db.QueryContext: %w", err)
	}

	for rows.Next() {
		var c model.Credit
		if err = rows.Scan(
			&c.ID,
			&c.AccountID,
			&c.CreditApplicationID,
			&c.UserID,
			&c.Amount,
			&c.InterestRate,
			&c.RemainingAmount,
			&c.Status,
			&c.StartDate,
			&c.EndDate,
			&c.CreatedAt,
			&c.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed rows.Scan: %w", err)
		}

		credits = append(credits, c)
	}

	return credits, nil
}
