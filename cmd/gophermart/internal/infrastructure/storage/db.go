package storage

import (
	"context"

	"kialkuz/shop-with-loyalty/pkg/pgerrors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	pool       *pgxpool.Pool
	classifier *pgerrors.PostgresErrorClassifier
}

func NewDB(pool *pgxpool.Pool) *DB {
	return &DB{
		pool:       pool,
		classifier: pgerrors.NewPostgresErrorClassifier(),
	}
}

func (db *DB) Exec(ctx context.Context, sql string, args ...any) error {
	_, err := db.pool.Exec(ctx, sql, args...)

	return err
}

func (db *DB) QueryRow(
	ctx context.Context,
	scan func(pgx.Row) error,
	sql string,
	args ...any,
) error {
	row := db.pool.QueryRow(ctx, sql, args...)

	return scan(row)
}

func (db *DB) QueryRows(
	ctx context.Context,
	sql string,
	args ...any,
) (pgx.Rows, error) {
	rows, err := db.pool.Query(ctx, sql, args...)

	return rows, err
}
