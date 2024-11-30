package db

import (
	"context"
	"fmt"
	"github.com/patyukin/mbs-credits/internal/model"
	desc "github.com/patyukin/mbs-pkg/pkg/proto/credit_v1"
	"time"
)

func (r *Repository) InsertCreditApplication(ctx context.Context, in *desc.CreateCreditApplicationRequest) (string, error) {
	currentTime := time.Now().UTC()
	query := `
INSERT INTO credit_applications (
	user_id,
	requested_amount,
	interest_rate,
	status,
	created_at
) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	row := r.db.QueryRowContext(
		ctx,
		query,
		in.UserId,
		in.RequestedAmount,
		in.InterestRate,
		"PENDING",
		currentTime,
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

func (r *Repository) SelectCreditApplicationByIDAndUserID(ctx context.Context, id string, userID string) (model.CreditApplication, error) {
	query := `
SELECT
	id,
	user_id,
	requested_amount,
	interest_rate,
	status,
	decision_date,
	approved_amount,
	decision_notes,
	description,
	created_at,
	updated_at
FROM credit_applications WHERE id = $1 AND user_id = $2 AND status = 'APPROVED'`
	row := r.db.QueryRowContext(ctx, query, id, userID)
	if row.Err() != nil {
		return model.CreditApplication{}, fmt.Errorf("failed r.db.QueryRowContext, row.Err(): %w", row.Err())
	}

	var ca model.CreditApplication
	if err := row.Scan(
		&ca.ID,
		&ca.UserID,
		&ca.RequestedAmount,
		&ca.InterestRate,
		&ca.Status,
		&ca.DecisionDate,
		&ca.ApprovedAmount,
		&ca.DecisionNotes,
		&ca.Description,
		&ca.CreatedAt,
		&ca.UpdatedAt,
	); err != nil {
		return model.CreditApplication{}, fmt.Errorf("failed row.Scan: %w", err)
	}

	return ca, nil
}

func (r *Repository) UpdateCreditApplicationSolution(ctx context.Context, in model.CreditApplicationSolution) error {
	query := `
UPDATE credit_applications 
SET
	status = $1,
	decision_date = $2,
	approved_amount = $3,
	decision_notes = $4,
	updated_at = $5
WHERE id = $6`
	_, err := r.db.ExecContext(
		ctx,
		query,
		in.Status,
		in.DecisionDate,
		in.ApprovedAmount,
		in.DecisionNotes,
		time.Now().UTC(),
		in.CreditApplicationID,
	)
	if err != nil {
		return fmt.Errorf("failed r.db.ExecContext: %w", err)
	}

	return nil
}

func (r *Repository) UpdateCreditApplicationStatus(ctx context.Context, id, status string) error {
	query := `UPDATE credit_applications SET status = $1 WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, status, id)
	if err != nil {
		return fmt.Errorf("failed r.db.ExecContext: %w", err)
	}

	return nil
}

func (r *Repository) UpdateCreditApplicationsToArchivedStatus(ctx context.Context) error {
	t := time.Now().Add(-24 * time.Hour).UTC()
	query := `UPDATE credit_applications SET status = 'ARCHIVED' WHERE created_at > $1 AND status = 'PENDING'`
	_, err := r.db.ExecContext(ctx, query, t)
	if err != nil {
		return fmt.Errorf("failed r.db.ExecContext: %w", err)
	}

	return nil
}
