package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
)

type QueryExecutor interface {
	ExecContext(ctx context.Context, q string, args ...interface{}) (sql.Result, error)
	QueryContext(ctx context.Context, q string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, q string, args ...interface{}) *sql.Row
}

type Registry struct {
	db *sql.DB
}

func (registry *Registry) GetRepo() *Repository {
	return &Repository{
		db: registry.db,
	}
}

type Handler func(ctx context.Context, repo *Repository) error

func New(db *sql.DB) *Registry {
	return &Registry{db: db}
}

func (registry *Registry) ReadCommitted(ctx context.Context, f Handler) error {
	tx, err := registry.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		return fmt.Errorf("failed registry.db.BeginTx: %w", err)
	}

	defer func() {
		if err != nil && !errors.Is(err, sql.ErrTxDone) {
			if errRollback := tx.Rollback(); errRollback != nil {
				log.Error().Msgf("failed to rollback transaction: %v", errRollback)
			}
		}
	}()

	repo := &Repository{db: tx}

	if err = f(ctx, repo); err != nil {
		return fmt.Errorf("failed to execute handler: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (registry *Registry) Close() error {
	return registry.db.Close()
}
