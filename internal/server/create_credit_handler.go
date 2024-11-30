package server

import (
	"context"
	"fmt"
	"github.com/patyukin/mbs-pkg/pkg/errs"
	desc "github.com/patyukin/mbs-pkg/pkg/proto/credit_v1"
)

func (s Server) CreateCredit(ctx context.Context, in *desc.CreateCreditRequest) (*desc.CreateCreditResponse, error) {
	result, err := s.uc.CreateCreditUseCase(ctx, in)
	if err != nil {
		return &desc.CreateCreditResponse{
			Error: errs.ToErrorResponse(fmt.Errorf("failed s.uc.CreateCreditUseCase: %w", err)),
		}, nil
	}

	return result, nil
}
