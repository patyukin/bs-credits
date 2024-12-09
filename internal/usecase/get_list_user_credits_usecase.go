package usecase

import (
	"context"
	"fmt"
	"github.com/patyukin/mbs-credits/internal/db"
	"github.com/patyukin/mbs-credits/internal/model"
	desc "github.com/patyukin/mbs-pkg/pkg/proto/credit_v1"
)

func (u *UseCase) GetListUserCreditsUseCase(ctx context.Context, in *desc.GetListUserCreditsRequest) (*desc.GetListUserCreditsResponse, error) {
	var err error
	var credits []model.Credit
	var cnt int32

	err = u.registry.ReadCommitted(ctx, func(ctx context.Context, repo *db.Repository) error {
		cnt, err = repo.SelectCountCreditByUserID(ctx, in.UserId)
		if err != nil {
			return fmt.Errorf("failed repo.SelectCountCreditByUserID: %w", err)
		}

		if cnt == 0 {
			return nil
		}

		credits, err = repo.SelectCreditByUserID(ctx, in.UserId, in.Page, in.Limit)
		if err != nil {
			return fmt.Errorf("failed repo.SelectCreditByUserID: %w", err)
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed u.registry.ReadCommitted: %w", err)
	}

	creditsProto, err := model.ToProtoCredits(credits)
	if err != nil {
		return nil, fmt.Errorf("failed model.ToProtoCredits: %w", err)
	}

	return &desc.GetListUserCreditsResponse{
		Credits: creditsProto,
		Total:   cnt,
	}, nil
}
