package cacher

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
	"time"
)

type Cacher struct {
	client *redis.Client
}

func New(ctx context.Context, dsn string) (*Cacher, error) {
	c := redis.NewClient(&redis.Options{Addr: dsn})

	err := c.Ping(ctx).Err()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	log.Info().Msg("connected to redis")
	return &Cacher{client: c}, nil
}

func (c *Cacher) Close() error {
	return c.client.Close()
}

func (c *Cacher) SetCreateApplicationConfirmationCode(ctx context.Context, ID, code string) error {
	return c.client.Set(ctx, fmt.Sprintf("ca:%s", ID), code, 24*time.Hour).Err()
}

func (c *Cacher) GetCreateApplicationConfirmationCode(ctx context.Context, ID string) (string, error) {
	return c.client.Get(ctx, fmt.Sprintf("ca:%s", ID)).Result()
}

func (c *Cacher) DeleteCreateApplicationConfirmationCode(ctx context.Context, ID string) error {
	return c.client.Del(ctx, fmt.Sprintf("ca:%s", ID)).Err()
}
