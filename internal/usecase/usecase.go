package usecase

import (
	"context"
	"github.com/patyukin/mbs-credits/internal/db"
	authpb "github.com/patyukin/mbs-pkg/pkg/proto/auth_v1"
	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/grpc"
)

type Cacher interface {
	SetCreateApplicationConfirmationCode(ctx context.Context, userID, creditApplicationID, code string) error
	GetCreateApplicationConfirmationCode(ctx context.Context, userID, code string) (string, error)
	DeleteCreateApplicationConfirmationCode(ctx context.Context, userID, creditApplicationID string) error
}

type KafkaProducer interface {
	PublishCreditCreated(ctx context.Context, value []byte) error
	PublishCreditPayments(ctx context.Context, value []byte) error
}

type RabbitMQProducer interface {
	EnqueueTelegramMessage(ctx context.Context, body []byte, headers amqp.Table) error
}

type AuthClient interface {
	GetBriefUserByID(ctx context.Context, in *authpb.GetBriefUserByIDRequest, opts ...grpc.CallOption) (*authpb.GetBriefUserByIDResponse, error)
}

type UseCase struct {
	registry      *db.Registry
	kafkaProducer KafkaProducer
	cacher        Cacher
	rbt           RabbitMQProducer
	authClient    AuthClient
}

func New(registry *db.Registry, kafkaProducer KafkaProducer, cacher Cacher, rbt RabbitMQProducer, authClient AuthClient) *UseCase {
	return &UseCase{
		registry:      registry,
		kafkaProducer: kafkaProducer,
		cacher:        cacher,
		rbt:           rbt,
		authClient:    authClient,
	}
}
