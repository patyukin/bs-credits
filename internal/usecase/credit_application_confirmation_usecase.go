package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/patyukin/mbs-credits/internal/db"
	"github.com/patyukin/mbs-pkg/pkg/model"
	authpb "github.com/patyukin/mbs-pkg/pkg/proto/auth_v1"
	desc "github.com/patyukin/mbs-pkg/pkg/proto/credit_v1"
)

func (u *UseCase) CreditApplicationConfirmationUseCase(ctx context.Context, in *desc.CreditApplicationConfirmationRequest) (*desc.CreditApplicationConfirmationResponse, error) {
	err := u.registry.ReadCommitted(ctx, func(ctx context.Context, repo *db.Repository) error {
		creditApplicationID, err := u.cacher.GetCreateApplicationConfirmationCode(ctx, in.UserId, in.Code)
		if err != nil {
			return fmt.Errorf("failed GetCreateApplicationConfirmationCode: %w", err)
		}

		if err = u.cacher.DeleteCreateApplicationConfirmationCode(ctx, in.UserId, creditApplicationID); err != nil {
			return fmt.Errorf("failed DeleteCreateApplicationConfirmationCode: %w", err)
		}

		user, err := u.authClient.GetBriefUserByID(ctx, &authpb.GetBriefUserByIDRequest{UserId: in.UserId})
		if err != nil {
			return fmt.Errorf("failed u.authClient.GetBriefUserByID: %w", err)
		}

		msg := model.SimpleTelegramMessage{ChatID: user.ChatId, Message: "Ваша заявка на кредит успешно подтверждена"}
		bytes, err := json.Marshal(msg)
		if err != nil {
			return fmt.Errorf("failed to marshal message: %w", err)
		}

		if err = u.rbt.EnqueueTelegramMessage(ctx, bytes, nil); err != nil {
			return fmt.Errorf("failed u.rbt.EnqueueTelegramMessage: %w", err)
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed u.registry.ReadCommitted: %w", err)
	}

	return &desc.CreditApplicationConfirmationResponse{Message: "Заявка успешно подтверждена"}, nil
}
