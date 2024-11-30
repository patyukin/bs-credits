package server

import (
	"context"
	desc "github.com/patyukin/mbs-pkg/pkg/proto/credit_v1"
)

func (s Server) CreditApplicationConfirmation(ctx context.Context, in *desc.CreditApplicationConfirmationRequest) (*desc.CreditApplicationConfirmationResponse, error) {
	result, err := s.uc.CreditApplicationConfirmationUseCase(ctx, in)
}
