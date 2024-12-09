package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/patyukin/mbs-credits/internal/db"
	"github.com/patyukin/mbs-credits/internal/model"
	pkgModel "github.com/patyukin/mbs-pkg/pkg/model"
	authpb "github.com/patyukin/mbs-pkg/pkg/proto/auth_v1"
	desc "github.com/patyukin/mbs-pkg/pkg/proto/credit_v1"
	"github.com/rs/zerolog/log"
)

func (u *UseCase) UpdateCreditApplicationSolutionUseCase(ctx context.Context, in *desc.UpdateCreditApplicationSolutionRequest) (*desc.UpdateCreditApplicationSolutionResponse, error) {
	err := u.registry.ReadCommitted(
		ctx, func(ctx context.Context, repo *db.Repository) error {
			pbm, err := model.ToModelCreditApplicationSolution(in)
			log.Debug().Msgf("in: %+v, pbm: %+v, err: %+v", in, pbm, err)
			if err != nil {
				return fmt.Errorf("failed model.ToModelCreditApplicationSolution: %w", err)
			}

			if err = repo.UpdateCreditApplicationSolution(ctx, pbm); err != nil {
				return fmt.Errorf("failed repo.UpdateCreditApplicationSolution: %w", err)
			}

			creditApplication, err := repo.SelectCreditApplicationByID(ctx, pbm.CreditApplicationID, pbm.CreditApplicationID)
			if err != nil {
				return fmt.Errorf("failed repo.SelectCreditApplicationByIDAndUserID: %w", err)
			}

			user, err := u.authClient.GetBriefUserByID(ctx, &authpb.GetBriefUserByIDRequest{UserId: creditApplication.UserID})
			if err != nil {
				return fmt.Errorf("failed u.authClient.GetBriefUserByID: %w", err)
			}

			msg := pkgModel.SimpleTelegramMessage{
				Message: fmt.Sprintf("обновлена заявка на кредит: %s", creditApplication.Description),
				ChatID:  user.ChatId,
			}

			msgBytes, err := json.Marshal(msg)
			if err != nil {
				return fmt.Errorf("failed to marshal message: %w", err)
			}

			if err = u.rbt.EnqueueTelegramMessage(ctx, msgBytes, nil); err != nil {
				return fmt.Errorf("failed u.rbt.EnqueueTelegramMessage: %w", err)
			}

			return nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed u.registry.ReadCommitted: %w", err)
	}

	return &desc.UpdateCreditApplicationSolutionResponse{Message: "Заявка успешно обновлена"}, nil
}
