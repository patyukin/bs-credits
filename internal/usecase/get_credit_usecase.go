package usecase

import (
	"context"
	"fmt"
	"github.com/patyukin/mbs-credits/internal/db"
	"github.com/patyukin/mbs-credits/internal/model"
	desc "github.com/patyukin/mbs-pkg/pkg/proto/credit_v1"
)

func (u *UseCase) GetCreditUseCase(ctx context.Context, in *desc.GetCreditRequest) (*desc.GetCreditResponse, error) {
	var credit model.Credit
	var err error

	err = u.registry.ReadCommitted(ctx, func(ctx context.Context, repo *db.Repository) error {
		credit, err = repo.SelectCreditByIDAndUserID(ctx, in.CreditId, in.UserId)
		if err != nil {
			return fmt.Errorf("failed repo.SelectCreditByID: %w", err)
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed u.registry.ReadCommitted: %w", err)
	}

	return model.ToProtoCredit(credit)
}
