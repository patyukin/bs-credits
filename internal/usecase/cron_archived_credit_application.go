package usecase

import (
	"context"
	"fmt"
)

func (u *UseCase) ArchivedCreditApplication(ctx context.Context) error {
	err := u.registry.GetRepo().UpdateCreditApplicationsToArchivedStatus(ctx)
	if err != nil {
		return fmt.Errorf("failed u.registry.GetRepo().UpdateCreditApplicationsToArchivedStatus: %w", err)
	}

	return nil
}
