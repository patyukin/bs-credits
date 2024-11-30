package server

import (
	"context"
	"fmt"
	"github.com/patyukin/mbs-pkg/pkg/errs"
	desc "github.com/patyukin/mbs-pkg/pkg/proto/credit_v1"
)

func (s Server) GetPaymentSchedule(ctx context.Context, in *desc.GetPaymentScheduleRequest) (*desc.GetPaymentScheduleResponse, error) {
	result, err := s.uc.GetPaymentScheduleUseCase(ctx, in)
	if err != nil {
		return &desc.GetPaymentScheduleResponse{
			Error: errs.ToErrorResponse(fmt.Errorf("failed s.uc.GetPaymentScheduleUseCase: %w", err)),
		}, nil
	}

	return result, nil
}
