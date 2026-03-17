package db

import (
	"context"
	"database/sql"
	"fmt"
)

type txContextKey struct{}

type TxManager interface {
	WithinTransaction(ctx context.Context, fn func(ctx context.Context, tx *sql.Tx) error) error
}

type SQLTxManager struct {
	db *sql.DB
}

func NewTxManager(db *sql.DB) *SQLTxManager {
	return &SQLTxManager{db: db}
}

func ContextWithTx(ctx context.Context, tx *sql.Tx) context.Context {
	if tx == nil {
		return ctx
	}

	return context.WithValue(ctx, txContextKey{}, tx)
}

func TxFromContext(ctx context.Context) (*sql.Tx, bool) {
	tx, ok := ctx.Value(txContextKey{}).(*sql.Tx)
	return tx, ok
}

func (manager *SQLTxManager) WithinTransaction(ctx context.Context, fn func(ctx context.Context, tx *sql.Tx) error) error {
	tx, err := manager.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}

	txCtx := ContextWithTx(ctx, tx)
	if err := fn(txCtx, tx); err != nil {
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
