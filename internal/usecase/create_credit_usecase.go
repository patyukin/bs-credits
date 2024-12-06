package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/patyukin/mbs-credits/internal/db"
	"github.com/patyukin/mbs-credits/internal/model"
	kafkaModel "github.com/patyukin/mbs-pkg/pkg/model"
	desc "github.com/patyukin/mbs-pkg/pkg/proto/credit_v1"
	"github.com/rs/zerolog/log"
	"math"
	"time"
)

type Loan struct {
	CreditID       string    // Идентификатор кредита
	Principal      int64     // Основная сумма кредита в копейках
	InterestRate   int32     // Годовая процентная ставка в процентах (например, 20.0 для 20%)
	NumPayments    int       // Общее количество платежей (месяцев)
	StartDate      time.Time // Дата начала кредита
	MonthlyPayment int64     // Ежемесячный платеж в копейках
}

func (u *UseCase) CreateCreditUseCase(ctx context.Context, in *desc.CreateCreditRequest) (*desc.CreateCreditResponse, error) {
	err := u.registry.ReadCommitted(ctx, func(ctx context.Context, repo *db.Repository) error {
		creditApplication, err := repo.SelectCreditApplicationByIDAndUserID(ctx, in.ApplicationId, in.UserId)
		if err != nil {
			return fmt.Errorf("failed repo.SelectCreditApplicationByIDAndUserID: %w", err)
		}

		numPayments := int(in.CreditTermMonths)
		monthlyPayment := calculateMonthlyPayment(creditApplication.RequestedAmount, creditApplication.InterestRate, numPayments)
		totalPaid := calculateTotalPaid(monthlyPayment, numPayments)

		c := model.ToModelCredit(creditApplication, in, totalPaid)
		creditID, err := repo.InsertCredit(ctx, c)
		if err != nil {
			return fmt.Errorf("failed repo.InsertCredit: %w", err)
		}

		loan := Loan{
			CreditID:       creditID,
			Principal:      creditApplication.RequestedAmount,
			InterestRate:   creditApplication.InterestRate,
			NumPayments:    numPayments,
			StartDate:      c.StartDate,
			MonthlyPayment: monthlyPayment,
		}

		schedules := generatePaymentSchedule(loan)
		if err = repo.InsertPaymentSchedules(ctx, schedules); err != nil {
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

		if err = u.kafkaProducer.PublishCreditCreated(ctx, value); err != nil {
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
func calculateMonthlyPayment(principal int64, annualInterestRate int32, numPayments int) int64 {
	log.Debug().Msgf("principal: %v, annualInterestRate: %v, numPayments: %v", principal, annualInterestRate, numPayments)

	P := float64(principal)          // Основная сумма кредита в копейках
	R := float64(annualInterestRate) // Годовая процентная ставка в процентах
	r := R / 12 / 100                // Ежемесячная процентная ставка в десятичном виде (например, 20% годовых -> 0.016666...)
	n := float64(numPayments)        // Общее количество платежей

	// Формула аннуитетного платежа
	pmt := P * r / (1 - math.Pow(1+r, -n))
	return int64(math.Round(pmt)) // Округляем до ближайшей копейки
}

// calculateTotalPaid - вычисляет общую сумму выплат по кредиту
func calculateTotalPaid(monthlyPayment int64, numPayments int) int64 {
	total := float64(monthlyPayment) * float64(numPayments)
	result := math.Round(total*100) / 100 // Округляем до копеек
	return int64(result)
}

// generatePaymentSchedule - генерирует график платежей по кредиту и рассчитывает общую сумму выплат
func generatePaymentSchedule(loan Loan) []model.PaymentSchedule {
	var schedules []model.PaymentSchedule

	remainingPrincipal := float64(loan.Principal)
	R := float64(loan.InterestRate)
	r := R / 12 / 100 // Ежемесячная процентная ставка в десятичном виде

	for i := 0; i < loan.NumPayments; i++ {
		// Расчет процентов за текущий месяц
		interestPayment := remainingPrincipal * r
		// Расчет погашения основного долга
		principalPayment := float64(loan.MonthlyPayment) - interestPayment
		// Обновление остатка долга
		remainingPrincipal -= principalPayment

		// Обеспечиваем, что остаток долга не станет отрицательным из-за округления
		if remainingPrincipal < 0 {
			principalPayment += remainingPrincipal
			remainingPrincipal = 0
		}

		// Дата следующего платежа
		log.Debug().Msgf("startDate %v, i: %v", loan.StartDate, i)
		dueDate := loan.StartDate.AddDate(0, i+1, 0)
		log.Debug().Msgf("startDate %v, i: %v, dueDate: %v", loan.StartDate, i, dueDate)
		schedule := model.PaymentSchedule{
			CreditID: loan.CreditID,
			Amount:   loan.MonthlyPayment,
			DueDate:  dueDate,
			Status:   "SCHEDULED",
		}

		schedules = append(schedules, schedule)
	}

	log.Debug().Msgf("schedules: %v", schedules)

	return schedules
}
