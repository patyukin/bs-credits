package model

import (
	"database/sql"
	"time"
)

type CreditApplicationSolution struct {
	CreditApplicationID string    `json:"credit_application_id"`
	Status              string    `json:"status"`
	DecisionDate        time.Time `json:"decision_date"`
	ApprovedAmount      int64     `json:"approved_amount"`
	DecisionNotes       string    `json:"decision_notes"`
}

type CreditApplication struct {
	ID              string         `json:"id,omitempty"`
	UserID          string         `json:"user_id,omitempty"`
	RequestedAmount int64          `json:"requested_amount,omitempty"`
	InterestRate    int32          `json:"interest_rate,omitempty"`
	Status          string         `json:"status,omitempty"`
	Description     string         `json:"description"`
	DecisionDate    sql.NullTime   `json:"decision_date"`
	ApprovedAmount  sql.NullInt64  `json:"approved_amount"`
	DecisionNotes   sql.NullString `json:"decision_notes"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       sql.NullTime   `json:"updated_at"`
}

type Credit struct {
	ID                  string       `json:"id,omitempty"`
	AccountID           string       `json:"account_id,omitempty"`
	CreditApplicationID string       `json:"credit_application_id,omitempty"`
	UserID              string       `json:"user_id,omitempty"`
	Amount              int64        `json:"amount,omitempty"`
	InterestRate        int32        `json:"interest_rate,omitempty"`
	RemainingAmount     int64        `json:"remaining_amount,omitempty"`
	Status              string       `json:"status,omitempty"`
	StartDate           time.Time    `json:"start_date"`
	EndDate             time.Time    `json:"end_date"`
	CreatedAt           time.Time    `json:"created_at"`
	UpdatedAt           sql.NullTime `json:"updated_at"`
}

type PaymentSchedule struct {
	ID        string       `json:"id,omitempty"`
	CreditID  string       `json:"credit_id,omitempty"`
	Amount    int64        `json:"amount,omitempty"`
	DueDate   time.Time    `json:"due_date"`
	Status    string       `json:"status,omitempty"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt sql.NullTime `json:"updated_at"`
}

type CreditPaymentSchedule struct {
	PaymentScheduleID string `json:"payment_schedule_id"`
	Amount            int64  `json:"amount"`
	AccountID         string `json:"account_id"`
}
