package model

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/patyukin/mbs-pkg/pkg/mapping/creditmapper"
	desc "github.com/patyukin/mbs-pkg/pkg/proto/credit_v1"
)

func ToModelCredit(creditApplication CreditApplication, in *desc.CreateCreditRequest, totalPaid int64) Credit {
	startDate := time.Now()
	endDate := startDate.AddDate(0, int(in.CreditTermMonths), 0)

	return Credit{
		ID:                  uuid.New().String(),
		AccountID:           in.AccountId,
		CreditApplicationID: creditApplication.ID,
		UserID:              creditApplication.UserID,
		Amount:              creditApplication.ApprovedAmount.Int64,
		InterestRate:        creditApplication.InterestRate,
		RemainingAmount:     totalPaid,
		Status:              "ACTIVE",
		StartDate:           startDate,
		EndDate:             endDate,
		CreatedAt:           time.Now().UTC(),
	}
}

func ToProtoCreditApplication(in CreditApplication) (*desc.GetCreditApplicationResponse, error) {
	status, err := creditmapper.StringToEnumCreditApplicationStatus(in.Status)
	if err != nil {
		return nil, fmt.Errorf("failed creditmapper.StringToEnumCreditApplicationStatus: %w", err)
	}

	decisionDateStr := in.DecisionDate.Time.Format(time.DateOnly)
	r := &desc.GetCreditApplicationResponse{
		ApplicationId:  in.ID,
		Status:         status,
		DecisionDate:   decisionDateStr,
		ApprovedAmount: in.ApprovedAmount.Int64,
		Description:    in.Description,
	}

	return r, nil
}

func ToModelCreditApplicationSolution(in *desc.UpdateCreditApplicationSolutionRequest) (CreditApplicationSolution, error) {
	status, err := creditmapper.EnumToStringCreditApplicationStatus(in.Status)
	if err != nil {
		return CreditApplicationSolution{}, fmt.Errorf("failed creditmapper.EnumToStringCreditApplicationStatus: %w", err)
	}

	return CreditApplicationSolution{
		CreditApplicationID: in.ApplicationId,
		Status:              status,
		ApprovedAmount:      in.ApprovedAmount,
		DecisionNotes:       in.DecisionNotes,
	}, nil
}

func ToProtoCredit(in Credit) (*desc.GetCreditResponse, error) {
	status, err := creditmapper.StringToEnumCreditStatus(in.Status)
	if err != nil {
		return nil, fmt.Errorf("failed creditmapper.StringToEnumCreditStatus: %w", err)
	}

	r := &desc.GetCreditResponse{
		Credit: &desc.Credit{
			CreditId:            in.ID,
			AccountId:           in.AccountID,
			CreditApplicationId: in.CreditApplicationID,
			UserId:              in.UserID,
			Amount:              in.Amount,
			InterestRate:        in.InterestRate,
			RemainingAmount:     in.RemainingAmount,
			Status:              status,
			StartDate:           in.StartDate.Format(time.DateOnly),
			EndDate:             in.EndDate.Format(time.DateOnly),
			CreatedAt:           in.CreatedAt.Format(time.DateOnly),
		},
	}

	if in.UpdatedAt.Valid {
		r.Credit.UpdatedAt = in.UpdatedAt.Time.Format(time.DateOnly)
	}

	return r, nil
}

func ToProtoCredits(credits []Credit) ([]*desc.Credit, error) {
	var r []*desc.Credit
	for _, c := range credits {
		status, err := creditmapper.StringToEnumCreditStatus(c.Status)
		if err != nil {
			return nil, fmt.Errorf("failed creditmapper.StringToEnumCreditStatus: %w", err)
		}

		r = append(
			r, &desc.Credit{
				CreditId:            c.ID,
				AccountId:           c.AccountID,
				CreditApplicationId: c.CreditApplicationID,
				UserId:              c.UserID,
				Amount:              c.Amount,
				InterestRate:        c.InterestRate,
				RemainingAmount:     c.RemainingAmount,
				Status:              status,
				StartDate:           c.StartDate.Format(time.DateOnly),
				EndDate:             c.EndDate.Format(time.DateOnly),
				CreatedAt:           c.CreatedAt.Format(time.DateOnly),
			},
		)
	}

	return r, nil
}

func ToProtoPaymentSchedule(payments []PaymentSchedule) (*desc.GetPaymentScheduleResponse, error) {
	var r []*desc.PaymentSchedule

	for _, p := range payments {
		status, err := creditmapper.StringToEnumPaymentStatus(p.Status)
		if err != nil {
			return nil, fmt.Errorf("failed creditmapper.StringToEnumPaymentScheduleStatus: %w", err)
		}

		r = append(
			r, &desc.PaymentSchedule{
				PaymentId: p.ID,
				Amount:    p.Amount,
				DueDate:   p.DueDate.Format(time.DateOnly),
				Status:    status,
			},
		)
	}

	return &desc.GetPaymentScheduleResponse{
		Payments: r,
	}, nil
}
