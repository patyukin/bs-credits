package usecase

import (
	"context"
	"fmt"
	"github.com/patyukin/mbs-credits/internal/db"
	"github.com/patyukin/mbs-credits/internal/model"
	desc "github.com/patyukin/mbs-pkg/pkg/proto/credit_v1"
)

func (u *UseCase) GetCreditApplicationUseCase(ctx context.Context, in *desc.GetCreditApplicationRequest) (*desc.GetCreditApplicationResponse, error) {
	var creditApplication model.CreditApplication
	var err error

	err = u.registry.ReadCommitted(ctx, func(ctx context.Context, repo *db.Repository) error {
		creditApplication, err = repo.SelectCreditApplicationByIDAndUserID(ctx, in.ApplicationId, in.UserId)
		if err != nil {
			return fmt.Errorf("failed repo.SelectCreditApplicationByID: %w", err)
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed u.registry.ReadCommitted: %w", err)
	}

	return model.ToProtoCreditApplication(creditApplication)
}
