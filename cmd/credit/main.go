package main

import (
	"context"
	"fmt"
	grpcPrometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/patyukin/mbs-credits/internal/cacher"
	"github.com/patyukin/mbs-credits/internal/config"
	"github.com/patyukin/mbs-credits/internal/cronjob"
	"github.com/patyukin/mbs-credits/internal/db"
	"github.com/patyukin/mbs-credits/internal/metrics"
	"github.com/patyukin/mbs-credits/internal/server"
	"github.com/patyukin/mbs-credits/internal/usecase"
	"github.com/patyukin/mbs-pkg/pkg/dbconn"
	"github.com/patyukin/mbs-pkg/pkg/grpc_client"
	"github.com/patyukin/mbs-pkg/pkg/grpc_server"
	"github.com/patyukin/mbs-pkg/pkg/kafka"
	"github.com/patyukin/mbs-pkg/pkg/migrator"
	"github.com/patyukin/mbs-pkg/pkg/mux_server"
	authpb "github.com/patyukin/mbs-pkg/pkg/proto/auth_v1"
	desc "github.com/patyukin/mbs-pkg/pkg/proto/credit_v1"
	"github.com/patyukin/mbs-pkg/pkg/rabbitmq"
	"github.com/patyukin/mbs-pkg/pkg/tracing"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const ServiceName = "CreditService"

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal().Msgf("failed to load config, error: %v", err)
	}

	if err = metrics.Init(); err != nil {
		log.Fatal().Msgf("failed to init metrics: %v", err)
	}

	_, closer, err := tracing.InitJaeger(fmt.Sprintf("jaeger:6831"), ServiceName)
	if err != nil {
		log.Fatal().Msgf("failed to initialize tracer: %v", err)
	}

	defer closer()

	log.Info().Msg("Jaeger connected")

	log.Info().Msg("Opentracing connected")
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.GRPCServer.Port))
	if err != nil {
		log.Fatal().Msgf("failed to listen: %v", err)
	}

	dbConn, err := dbconn.New(ctx, cfg.PostgreSQLDSN)
	if err != nil {
		log.Fatal().Msgf("failed to connect to db: %v", err)
	}

	if err = migrator.UpMigrations(ctx, dbConn); err != nil {
		log.Fatal().Msgf("failed to up migrations: %v", err)
	}

	rbt, err := rabbitmq.New(cfg.RabbitMQUrl, rabbitmq.Exchange)
	if err != nil {
		log.Fatal().Msgf("failed to create rabbit producer: %v", err)
	}

	err = rbt.BindQueueToExchange(
		rabbitmq.Exchange,
		rabbitmq.TelegramMessageQueue,
		[]string{rabbitmq.TelegramMessageRouteKey},
	)
	if err != nil {
		log.Fatal().Msgf("failed to bind TelegramMessageQueue to exchange with - TelegramMessageRouteKey: %v", err)
	}

	chr, err := cacher.New(ctx, cfg.RedisDSN)
	if err != nil {
		log.Fatal().Msgf("failed to create redis cacher: %v", err)
	}

	kfk, err := kafka.NewProducer(cfg.Kafka.Brokers)
	if err != nil {
		log.Fatal().Msgf("failed to create kafka producer: %v", err)
	}

	kafkaConsumer, err := kafka.NewConsumer(cfg.Kafka.Brokers, cfg.Kafka.ConsumerGroup, cfg.Kafka.Topics)
	if err != nil {
		log.Fatal().Msgf("failed to create kafka consumer: %v", err)
	}

	registry := db.New(dbConn)

	// auth service init
	authConn, err := grpc_client.NewGRPCClientServiceConn(cfg.GRPC.AuthService)
	if err != nil {
		log.Fatal().Msgf("failed to connect to auth service: %v", err)
	}

	defer func(authConn *grpc.ClientConn) {
		if err = authConn.Close(); err != nil {
			log.Error().Msgf("failed to close auth service connection: %v", err)
		}
	}(authConn)

	authClient := authpb.NewAuthServiceClient(authConn)

	// Создаем Dispatcher

	uc := usecase.New(registry, kfk, chr, rbt, authClient)
	srv := server.New(uc)

	// grpc server
	s := grpc_server.NewGRPCServer()
	reflection.Register(s)
	desc.RegisterCreditsServiceV1Server(s, srv)
	grpcPrometheus.Register(s)

	// mux server
	muxServer := mux_server.New()

	// cron job
	cj := cronjob.New(uc)

	log.Printf("server listening at %v", lis.Addr())

	errCh := make(chan error)

	// run consumer
	go func() {
		if err = kafkaConsumer.ProcessMessages(ctx, uc.ConsumePaymentScheduleSolution); err != nil {
			log.Error().Msgf("failed to process messages: %v", err)
			errCh <- err
		}
	}()

	go func() {
		if err = cj.Run(ctx); err != nil {
			log.Error().Msgf("failed to run cronjob: %v", err)
			errCh <- err
		}
	}()

	// GRPC server
	go func() {
		log.Info().Msgf("GRPC started on :%d", cfg.GRPCServer.Port)
		if err = s.Serve(lis); err != nil {
			log.Error().Msgf("failed to serve: %v", err)
			errCh <- err
		}
	}()

	// metrics + pprof server
	go func() {
		if err = muxServer.Run(cfg.HTTPServer.Port); err != nil {
			log.Error().Msgf("Failed to serve Prometheus metrics: %v", err)
			errCh <- err
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)

	select {
	case err = <-errCh:
		log.Error().Msgf("Failed to run, err: %v", err)
	case res := <-sigChan:
		if res == syscall.SIGINT || res == syscall.SIGTERM {
			log.Info().Msg("Signal received")
		} else if res == syscall.SIGHUP {
			log.Info().Msg("Signal received")
		}
	}

	log.Info().Msg("Shutting Down")

	// stop server
	s.GracefulStop()

	if err = muxServer.Shutdown(ctx); err != nil {
		log.Error().Msgf("failed to shutdown http server: %s", err.Error())
	}

	if err = rbt.Close(); err != nil {
		log.Error().Msgf("failed rabbit connection close: %s", err.Error())
	}

	if err = dbConn.Close(); err != nil {
		log.Error().Msgf("failed db connection close: %s", err.Error())
	}

	if err = chr.Close(); err != nil {
		log.Error().Msgf("failed redis connection close: %s", err.Error())
	}

	cj.Stop()
}
