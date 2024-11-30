package usecase

import (
	"context"
	"fmt"
	"github.com/patyukin/mbs-credits/internal/db"
	desc "github.com/patyukin/mbs-pkg/pkg/proto/credit_v1"
)

func (u *UseCase) CreditApplicationConfirmationUseCase(ctx context.Context, in *desc.CreditApplicationConfirmationRequest) (*desc.CreditApplicationConfirmationResponse, error) {
	err := u.registry.ReadCommitted(ctx, func(ctx context.Context, repo *db.Repository) error {
		// get code from redis
		// update credit application

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed u.registry.ReadCommitted: %w", err)
	}

	return &desc.CreditApplicationConfirmationResponse{Message: "Заявка успешно подтверждена"}, nil
}
