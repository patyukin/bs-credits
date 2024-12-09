package usecase

import (
	"context"
	"fmt"
	"github.com/patyukin/mbs-credits/internal/db"
	"github.com/patyukin/mbs-credits/internal/model"
	desc "github.com/patyukin/mbs-pkg/pkg/proto/credit_v1"
	"github.com/rs/zerolog/log"
)

func (u *UseCase) UpdateCreditApplicationSolutionUseCase(ctx context.Context, in *desc.UpdateCreditApplicationSolutionRequest) (*desc.UpdateCreditApplicationSolutionResponse, error) {
	err := u.registry.ReadCommitted(ctx, func(ctx context.Context, repo *db.Repository) error {
		pbm, err := model.ToModelCreditApplicationSolution(in)
		log.Debug().Msgf("in: %+v, pbm: %+v, err: %+v", in, pbm, err)
		if err != nil {
			return fmt.Errorf("failed model.ToModelCreditApplicationSolution: %w", err)
		}

		if err = repo.UpdateCreditApplicationSolution(ctx, pbm); err != nil {
			return fmt.Errorf("failed repo.UpdateCreditApplicationSolution: %w", err)
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed u.registry.ReadCommitted: %w", err)
	}

	return &desc.UpdateCreditApplicationSolutionResponse{Message: "Заявка успешно обновлена"}, nil
}
