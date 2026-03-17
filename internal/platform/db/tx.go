package db

import (
	"context"
	"database/sql"
	"fmt"
)

type TxManager interface {
	WithinTransaction(ctx context.Context, fn func(ctx context.Context, tx *sql.Tx) error) error
}

type SQLTxManager struct {
	db *sql.DB
}

func NewTxManager(db *sql.DB) *SQLTxManager {
	return &SQLTxManager{db: db}
}

func (manager *SQLTxManager) WithinTransaction(ctx context.Context, fn func(ctx context.Context, tx *sql.Tx) error) error {
	tx, err := manager.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}

	if err := fn(ctx, tx); err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return fmt.Errorf("rollback transaction: %v (original error: %w)", rollbackErr, err)
		}

		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}
