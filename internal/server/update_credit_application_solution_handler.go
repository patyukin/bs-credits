package server

import (
	"context"
	"fmt"
	"github.com/patyukin/mbs-pkg/pkg/errs"
	desc "github.com/patyukin/mbs-pkg/pkg/proto/credit_v1"
)

func (s Server) UpdateCreditApplicationSolution(ctx context.Context, in *desc.UpdateCreditApplicationSolutionRequest) (*desc.UpdateCreditApplicationSolutionResponse, error) {
	result, err := s.uc.UpdateCreditApplicationSolutionUseCase(ctx, in)
	if err != nil {
		return &desc.UpdateCreditApplicationSolutionResponse{
			Error: errs.ToErrorResponse(fmt.Errorf("failed s.uc.UpdateCreditApplicationSolutionUseCase: %w", err)),
		}, nil
	}

	return result, nil
}
