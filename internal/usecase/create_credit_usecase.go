package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/patyukin/mbs-credits/internal/db"
	"github.com/patyukin/mbs-credits/internal/model"
	kafkaModel "github.com/patyukin/mbs-pkg/pkg/model"
	desc "github.com/patyukin/mbs-pkg/pkg/proto/credit_v1"
	"math"
	"time"
)

func (u *UseCase) CreateCreditUseCase(ctx context.Context, in *desc.CreateCreditRequest) (*desc.CreateCreditResponse, error) {
	err := u.registry.ReadCommitted(ctx, func(ctx context.Context, repo *db.Repository) error {
		creditApplication, err := repo.SelectCreditApplicationByIDAndUserID(ctx, in.ApplicationId, in.UserId)
		if err != nil {
			return fmt.Errorf("failed repo.SelectCreditApplicationByIDAndUserID: %w", err)
		}

		c := model.ToModelCredit(creditApplication, in)
		creditID, err := repo.InsertCredit(ctx, c)
		if err != nil {
			return fmt.Errorf("failed repo.InsertCredit: %w", err)
		}

		numPayments := int(in.CreditTermMonths)
		schedules := generatePaymentSchedule(creditID, creditApplication.RequestedAmount, creditApplication.InterestRate, c.StartDate, numPayments)
		err = repo.InsertPaymentSchedules(ctx, schedules)
		if err != nil {
			return fmt.Errorf("failed repo.InsertPaymentSchedules: %w", err)
		}

		err = repo.UpdateCreditApplicationStatus(ctx, creditApplication.ID, "ARCHIVED")
		if err != nil {
			return fmt.Errorf("failed repo.UpdateCreditApplicationStatus: %w", err)
		}

		msg := kafkaModel.CreditCreated{AccountID: in.AccountId, Amount: creditApplication.RequestedAmount}
		value, err := json.Marshal(msg)
		if err != nil {
			return fmt.Errorf("failed json.Marshal: %w", err)
		}

		err = u.kafkaProducer.PublishCreditCreated(ctx, value)
		if err != nil {
			return fmt.Errorf("failed u.kafkaProducer.PublishCreditCreated: %w", err)
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed u.registry.ReadCommitted: %w", err)
	}

	return &desc.CreateCreditResponse{Message: "Кредит успешно создан"}, nil
}

// calculateMonthlyPayment - calculate monthly payment
func calculateMonthlyPayment(principal int64, annualInterestRate int64, numPayments int) int64 {
	P := float64(principal)
	R := float64(annualInterestRate)
	r := R / 12 / 100   // Ежемесячная процентная ставка в процентах
	rDecimal := r / 100 // Ежемесячная процентная ставка в десятичном виде
	n := float64(numPayments)

	pmt := P * rDecimal / (1 - math.Pow(1+rDecimal, -n))
	return int64(math.Round(pmt))
}

// generatePaymentSchedule - generate payment schedule
func generatePaymentSchedule(creditID string, principal int64, annualInterestRate int64, startDate time.Time, numPayments int) []model.PaymentSchedule {
	var schedules []model.PaymentSchedule

	monthlyPayment := calculateMonthlyPayment(principal, annualInterestRate, numPayments)
	remainingPrincipal := float64(principal)
	R := float64(annualInterestRate)
	r := R / 12 / 100   // Ежемесячная процентная ставка в процентах
	rDecimal := r / 100 // Ежемесячная процентная ставка в десятичном виде

	for i := 0; i < numPayments; i++ {
		// Расчет процентов за текущий месяц
		interestPayment := remainingPrincipal * rDecimal
		// Расчет погашения основного долга
		principalPayment := float64(monthlyPayment) - interestPayment
		// Обновление остатка долга
		remainingPrincipal -= principalPayment

		// Дата следующего платежа
		dueDate := startDate.AddDate(0, i+1, 0)
		schedule := model.PaymentSchedule{
			CreditID: creditID,
			Amount:   monthlyPayment,
			DueDate:  dueDate,
			Status:   "SCHEDULED",
		}

		schedules = append(schedules, schedule)
	}

	return schedules
}
