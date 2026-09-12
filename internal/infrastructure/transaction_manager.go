package infrastructure

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

//go:generate go run go.uber.org/mock/mockgen -destination=mocks/pgx_tx_mock.go -package=mocks github.com/jackc/pgx/v5 Tx
//go:generate go run go.uber.org/mock/mockgen -source=transaction_manager.go -destination=mocks/transaction_mock.go -package=mocks -typed
type Transaction interface {
	RunTransaction(
		ctx context.Context,
		fn func(pgx.Tx) error,
	) error
}

type TransactionManager struct {
	pool *pgxpool.Pool
}

func NewTransactionManager(pool *pgxpool.Pool) *TransactionManager {
	return &TransactionManager{pool: pool}
}

func (m *TransactionManager) RunTransaction(
	ctx context.Context,
	fn func(tx pgx.Tx) error,
) error {
	tx, err := m.pool.Begin(ctx)
	if err != nil {
		return err
	}

	defer tx.Rollback(ctx)

	if err := fn(tx); err != nil {
		return err
	}

	return tx.Commit(ctx)
}
