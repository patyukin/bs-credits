package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/patyukin/mbs-credits/internal/db"
	"github.com/patyukin/mbs-pkg/pkg/model"
	"github.com/rs/zerolog/log"
	"github.com/twmb/franz-go/pkg/kgo"
)

func (u *UseCase) PaymentScheduleSolutionProcessUseCase(ctx context.Context, record *kgo.Record) error {
	var msgs []model.CreditPaymentSolutionMessage

	log.Debug().Msgf("Received record: %v", string(record.Value))

	if err := json.Unmarshal(record.Value, &msgs); err != nil {
		return fmt.Errorf("failed to unmarshal message: %w", err)
	}

	err := u.registry.ReadCommitted(ctx, func(ctx context.Context, repo *db.Repository) error {
		for _, msg := range msgs {
			err := repo.UpdatePaymentScheduleStatus(ctx, msg.PaymentScheduleID, msg.Status)
			if err != nil {
				return fmt.Errorf("failed repo.UpdatePaymentScheduleSolution: %w", err)
			}
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("failed u.registry.ReadCommitted: %w", err)
	}

	return nil
}
