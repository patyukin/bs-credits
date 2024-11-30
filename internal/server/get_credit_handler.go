package server

import (
	"context"
	"fmt"
	"github.com/patyukin/mbs-pkg/pkg/errs"
	desc "github.com/patyukin/mbs-pkg/pkg/proto/credit_v1"
)

func (s Server) GetCredit(ctx context.Context, in *desc.GetCreditRequest) (*desc.GetCreditResponse, error) {
	result, err := s.uc.GetCreditUseCase(ctx, in)
	if err != nil {
		return &desc.GetCreditResponse{
			Error: errs.ToErrorResponse(fmt.Errorf("failed s.uc.GetCreditUseCase: %w", err)),
		}, nil
	}

	return result, nil
}
