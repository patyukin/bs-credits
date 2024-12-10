package usecase

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/patyukin/mbs-credits/internal/db"
	"github.com/patyukin/mbs-pkg/pkg/model"
	"github.com/rs/zerolog/log"
)

func (u *UseCase) CronCreditPayment(ctx context.Context) error {
	err := u.registry.ReadCommitted(
		ctx, func(ctx context.Context, repo *db.Repository) error {
			paymentSchedules, err := repo.SelectPaymentScheduleToPay(ctx)
			if err != nil {
				return fmt.Errorf("failed repo.SelectPaymentScheduleToPay: %w", err)
			}

			if len(paymentSchedules) == 0 {
				log.Info().Msg("no payment schedules to pay")
				return nil
			}

			if err = repo.UpdatePaymentSchedulesStatus(ctx, paymentSchedules); err != nil {
				return fmt.Errorf("failed repo.UpdatePaymentSchedulesStatus: %w", err)
			}

			var msgs []model.CreditPayment
			for _, paymentSchedule := range paymentSchedules {
				msgs = append(
					msgs, model.CreditPayment{
						AccountID:         paymentSchedule.AccountID,
						Amount:            paymentSchedule.Amount,
						PaymentScheduleID: paymentSchedule.PaymentScheduleID,
					},
				)
			}

			bytesMsgs, err := json.Marshal(msgs)
			if err != nil {
				return fmt.Errorf("failed to marshal messages: %w", err)
			}

			if err = u.kafkaProducer.PublishCreditPayments(ctx, bytesMsgs); err != nil {
				return fmt.Errorf("failed u.kafkaProducer.PublishCreditPayments: %w", err)
			}

			return nil
		},
	)
	if err != nil {
		return fmt.Errorf("failed u.registry.ReadCommitted: %w", err)
	}

	return nil
}
