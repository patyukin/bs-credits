package server

import (
	"context"
	desc "github.com/patyukin/mbs-pkg/pkg/proto/credit_v1"
)

type UseCase interface {
	CreateCreditApplicationUseCase(ctx context.Context, in *desc.CreateCreditApplicationRequest) (*desc.CreateCreditApplicationResponse, error)
	CreateCreditUseCase(ctx context.Context, in *desc.CreateCreditRequest) (*desc.CreateCreditResponse, error)
	CreditApplicationConfirmationUseCase(ctx context.Context, in *desc.CreditApplicationConfirmationRequest) (*desc.CreditApplicationConfirmationResponse, error)
	GetCreditApplicationUseCase(ctx context.Context, in *desc.GetCreditApplicationRequest) (*desc.GetCreditApplicationResponse, error)
	UpdateCreditApplicationSolutionUseCase(ctx context.Context, in *desc.UpdateCreditApplicationSolutionRequest) (*desc.UpdateCreditApplicationSolutionResponse, error)
	GetCreditUseCase(ctx context.Context, in *desc.GetCreditRequest) (*desc.GetCreditResponse, error)
	GetListUserCreditsUseCase(ctx context.Context, in *desc.GetListUserCreditsRequest) (*desc.GetListUserCreditsResponse, error)
	GetPaymentScheduleUseCase(ctx context.Context, in *desc.GetPaymentScheduleRequest) (*desc.GetPaymentScheduleResponse, error)
}

type Server struct {
	desc.UnimplementedCreditsServiceV1Server
	uc UseCase
}

func New(uc UseCase) *Server {
	return &Server{
		uc: uc,
	}
}
