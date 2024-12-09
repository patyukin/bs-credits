package usecase

import (
	"context"
	"fmt"
	"github.com/patyukin/mbs-credits/internal/db"
	"github.com/patyukin/mbs-credits/internal/model"
	desc "github.com/patyukin/mbs-pkg/pkg/proto/credit_v1"
)

func (u *UseCase) GetPaymentScheduleUseCase(ctx context.Context, in *desc.GetPaymentScheduleRequest) (*desc.GetPaymentScheduleResponse, error) {
	var err error
	var payments []model.PaymentSchedule

	err = u.registry.ReadCommitted(ctx, func(ctx context.Context, repo *db.Repository) error {
		payments, err = repo.SelectPaymentScheduleByUserIDAndCreditID(ctx, in.UserId, in.CreditId)
		if err != nil {
			return fmt.Errorf("failed repo.SelectPaymentScheduleByUserIDAndCreditID: %w", err)
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed u.registry.ReadCommitted: %w", err)
	}

	return model.ToProtoPaymentSchedule(payments)
}
