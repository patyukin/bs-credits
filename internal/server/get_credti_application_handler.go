package server

import (
	"context"
	"fmt"
	"github.com/patyukin/mbs-pkg/pkg/errs"
	desc "github.com/patyukin/mbs-pkg/pkg/proto/credit_v1"
)

func (s Server) GetCreditApplication(ctx context.Context, in *desc.GetCreditApplicationRequest) (*desc.GetCreditApplicationResponse, error) {
	result, err := s.uc.GetCreditApplicationUseCase(ctx, in)
	if err != nil {
		return &desc.GetCreditApplicationResponse{
			Error: errs.ToErrorResponse(fmt.Errorf("failed s.uc.GetCreditApplicationUseCase: %w", err)),
		}, nil
	}

	return result, nil
}
