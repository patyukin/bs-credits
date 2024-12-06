package cacher

import (
	"context"
	"errors"
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

func (c *Cacher) SetCreateApplicationConfirmationCode(ctx context.Context, userID, creditApplicationID, code string) error {
	err := c.client.Set(ctx, fmt.Sprintf("u:%s:ca:%s", userID, creditApplicationID), code, 50*time.Minute).Err()
	if err != nil {
		return fmt.Errorf("failed to set create application confirmation code: %w", err)
	}

	return nil
}

func (c *Cacher) GetCreateApplicationConfirmationCode(ctx context.Context, userID, code string) (string, error) {
	pattern := fmt.Sprintf("u:%s:ca:*", userID)
	keys, err := c.client.Keys(ctx, pattern).Result()
	if err != nil {
		return "", fmt.Errorf("failed to get keys: %w", err)
	}

	var storedCode string
	for _, key := range keys {
		storedCode, err = c.client.Get(ctx, key).Result()
		if errors.Is(err, redis.Nil) {
			continue
		}

		if err != nil {
			return "", fmt.Errorf("failed to get value for key %s: %w", key, err)
		}

		if storedCode == code {
			return key[len(fmt.Sprintf("u:%s:pc:", userID)):], nil
		}
	}

	return "", fmt.Errorf("payment confirmation code not found for userID: %s and code: %s", userID, code)
}

func (c *Cacher) DeleteCreateApplicationConfirmationCode(ctx context.Context, userID, creditApplicationID string) error {
	err := c.client.Del(ctx, fmt.Sprintf("u:%s:ca:%s", userID, creditApplicationID)).Err()
	if err != nil {
		return fmt.Errorf("failed to delete payment confirmation code: %w", err)
	}

	return nil
}
