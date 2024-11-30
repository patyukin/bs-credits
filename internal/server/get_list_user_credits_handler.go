package server

import (
	"context"
	"fmt"
	"github.com/patyukin/mbs-pkg/pkg/errs"
	desc "github.com/patyukin/mbs-pkg/pkg/proto/credit_v1"
)

func (s Server) GetListUserCredits(ctx context.Context, in *desc.GetListUserCreditsRequest) (*desc.GetListUserCreditsResponse, error) {
	result, err := s.uc.GetListUserCreditsUseCase(ctx, in)
	if err != nil {
		return &desc.GetListUserCreditsResponse{
			Error: errs.ToErrorResponse(fmt.Errorf("failed s.uc.GetListUserCreditsUseCase: %w", err)),
		}, nil
	}

	return result, nil
}
