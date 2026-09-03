package interfaces

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type TransactionManager interface {
	WithinTransaction(
		ctx context.Context,
		fn func(pgx.Tx) error,
	) error
}
