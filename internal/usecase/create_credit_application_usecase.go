package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/patyukin/mbs-credits/internal/db"
	"github.com/patyukin/mbs-pkg/pkg/model"
	authpb "github.com/patyukin/mbs-pkg/pkg/proto/auth_v1"
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

		err = u.cacher.SetCreateApplicationConfirmationCode(ctx, in.UserId, creditApplicationID, code.String())
		if err != nil {
			return fmt.Errorf("failed repo.UpdateCreditApplicationCode: %w", err)
		}

		user, err := u.authClient.GetBriefUserByID(ctx, &authpb.GetBriefUserByIDRequest{UserId: in.UserId})
		if err != nil {
			return fmt.Errorf("failed u.authClient.GetBriefUserByID: %w", err)
		}

		msg := model.SimpleTelegramMessage{
			ChatID:  user.ChatId,
			Message: fmt.Sprintf("Ваш код подтверждения заявки на кредит: %s", code.String()),
		}

		bytes, err := json.Marshal(msg)
		if err != nil {
			return fmt.Errorf("failed to marshal message: %w", err)
		}

		err = u.rbt.EnqueueTelegramMessage(ctx, bytes, nil)
		if err != nil {
			return fmt.Errorf("failed u.rbt.EnqueueTelegramMessage: %w", err)
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed u.registry.ReadCommitted: %w", err)
	}

	return &desc.CreateCreditApplicationResponse{Message: "Код подтверждения отправлен в telegram"}, nil
}
