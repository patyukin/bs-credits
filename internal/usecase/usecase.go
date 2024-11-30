package usecase

import (
	"context"
	"github.com/patyukin/mbs-credits/internal/db"
)

type Cacher interface {
	SetCreateApplicationConfirmationCode(ctx context.Context, ID, code string) error
	GetCreateApplicationConfirmationCode(ctx context.Context, ID string) (string, error)
	DeleteCreateApplicationConfirmationCode(ctx context.Context, ID string) error
}

type KafkaProducer interface {
	PublishCreditCreated(ctx context.Context, value []byte) error
	PublishCreditPayments(ctx context.Context, value []byte) error
}

type UseCase struct {
	registry      *db.Registry
	kafkaProducer KafkaProducer
	cacher        Cacher
}

func New(registry *db.Registry, kafkaProducer KafkaProducer, cacher Cacher) *UseCase {
	return &UseCase{
		registry:      registry,
		kafkaProducer: kafkaProducer,
		cacher:        cacher,
	}
}
