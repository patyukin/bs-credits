package server

import (
	"context"
	"fmt"
	"github.com/patyukin/mbs-pkg/pkg/errs"
	desc "github.com/patyukin/mbs-pkg/pkg/proto/credit_v1"
)

func (s Server) CreditApplicationConfirmation(ctx context.Context, in *desc.CreditApplicationConfirmationRequest) (*desc.CreditApplicationConfirmationResponse, error) {
	result, err := s.uc.CreditApplicationConfirmationUseCase(ctx, in)
	if err != nil {
		return &desc.CreditApplicationConfirmationResponse{
			Error: errs.ToErrorResponse(fmt.Errorf("failed s.uc.CreditApplicationConfirmationUseCase: %w", err)),
		}, nil
	}

	return result, nil
}
