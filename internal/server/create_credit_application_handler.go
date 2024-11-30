package server

import (
	"context"
	"fmt"
	"github.com/patyukin/mbs-pkg/pkg/errs"
	desc "github.com/patyukin/mbs-pkg/pkg/proto/credit_v1"
)

func (s Server) CreateCreditApplication(ctx context.Context, in *desc.CreateCreditApplicationRequest) (*desc.CreateCreditApplicationResponse, error) {
	result, err := s.uc.CreateCreditApplicationUseCase(ctx, in)
	if err != nil {
		return &desc.CreateCreditApplicationResponse{
			Error: errs.ToErrorResponse(fmt.Errorf("failed s.uc.CreateCreditApplicationUseCase: %w", err)),
		}, nil
	}

	return result, nil
}
