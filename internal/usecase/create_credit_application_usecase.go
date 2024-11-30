package usecase

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/patyukin/mbs-credits/internal/db"
	desc "github.com/patyukin/mbs-pkg/pkg/proto/credit_v1"
)

func (u *UseCase) CreateCreditApplicationUseCase(ctx context.Context, in *desc.CreateCreditApplicationRequest) (*desc.CreateCreditApplicationResponse, error) {
	err := u.registry.ReadCommitted(ctx, func(ctx context.Context, repo *db.Repository) error {
		creditApplicationID, err := repo.InsertCreditApplication(ctx, in)
		if err != nil {
			return fmt.Errorf("failed repo.CreateCreditApplication: %w", err)
		}

		code, err := uuid.NewUUID()
		if err != nil {
			return fmt.Errorf("failed uuid.NewUUID: %w", err)
		}

		err = u.cacher.SetCreateApplicationConfirmationCode(ctx, creditApplicationID, code.String())
		if err != nil {
			return fmt.Errorf("failed repo.UpdateCreditApplicationCode: %w", err)
		}

		// send to notify

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed u.registry.ReadCommitted: %w", err)
	}

	return &desc.CreateCreditApplicationResponse{
		Message: "Ваша заявка на кредит успешно создана",
	}, nil
}
