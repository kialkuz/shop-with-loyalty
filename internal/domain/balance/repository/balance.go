package db

import (
	"context"
	"fmt"
	"kialkuz/shop-with-loyalty/internal/domain/balance/model"
	"kialkuz/shop-with-loyalty/internal/infrastructure"
	pkgErrors "kialkuz/shop-with-loyalty/pkg/errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type BalanceRepository struct {
	db *infrastructure.DB
}

func NewBalanceRepository(db *infrastructure.DB) *BalanceRepository {
	return &BalanceRepository{db: db}
}

func (r *BalanceRepository) GetByUserId(ctx context.Context, userId uuid.UUID) (*model.Balance, error) {
	balance := &model.Balance{}
	err := r.db.QueryRow(ctx, func(row pgx.Row) error {
		return row.Scan(
			&balance.ID,
			&balance.UserId,
			&balance.Current,
			&balance.WithDrawn,
		)
	}, "SELECT id, user_id, current, withdrawn FROM balance WHERE user_id = $1", userId)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, pkgErrors.ErrNotFound
		}
		return nil, err
	}
	return balance, nil
}

func (r *BalanceRepository) GetByUserIdWithBlockForUpdate(
	ctx context.Context,
	tx pgx.Tx,
	userId uuid.UUID,
) (*model.Balance, error) {
	balance := &model.Balance{}
	row := tx.QueryRow(
		ctx,
		`SELECT id, user_id, current, withdrawn
		FROM balance
		WHERE user_id = $1 FOR UPDATE`,
		userId,
	)
	err := row.Scan(
		&balance.ID,
		&balance.UserId,
		&balance.Current,
		&balance.WithDrawn,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, pkgErrors.ErrNotFound
		}
		return nil, err
	}
	return balance, nil
}

func (r *BalanceRepository) getByUserId(ctx context.Context, query string, userId uuid.UUID) (*model.Balance, error) {
	balance := &model.Balance{}
	err := r.db.QueryRow(ctx, func(row pgx.Row) error {
		return row.Scan(
			&balance.ID,
			&balance.UserId,
			&balance.Current,
			&balance.WithDrawn,
		)
	}, query, userId)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, pkgErrors.ErrNotFound
		}
		return nil, err
	}
	return balance, nil
}

func (r *BalanceRepository) GetByUsersId(ctx context.Context, usersId []uuid.UUID) (map[uuid.UUID]model.Balance, error) {
	rows, err := r.db.QueryRows(
		ctx,
		"SELECT id, user_id, current, withdrawn FROM balance WHERE user_id = ANY($1::uuid[])",
		usersId,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, pkgErrors.ErrNotFound
		}
		return nil, err
	}

	defer rows.Close()

	usersBalance := make(map[uuid.UUID]model.Balance)

	for rows.Next() {
		o := model.Balance{}

		if err := rows.Scan(
			&o.ID,
			&o.UserId,
			&o.Current,
			&o.WithDrawn,
		); err != nil {
			return nil, fmt.Errorf("scan order: %w", err)
		}

		usersBalance[o.UserId] = o
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	if len(usersBalance) == 0 {
		return nil, pkgErrors.ErrNotFound
	}

	return usersBalance, nil
}

func (r *BalanceRepository) AddTx(ctx context.Context, tx pgx.Tx, balance model.Balance) error {
	_, err := tx.Exec(
		ctx,
		"INSERT INTO balance (id, user_id, current, withdrawn) VALUES ($1, $2, $3, $4)",
		balance.ID,
		balance.UserId,
		balance.Current,
		balance.WithDrawn,
	)

	return err
}

func (r *BalanceRepository) UpdateTx(ctx context.Context, tx pgx.Tx, balance model.Balance) error {
	_, err := tx.Exec(
		ctx,
		"UPDATE balance SET current = $1, withdrawn = $2 WHERE user_id = $3",
		balance.Current,
		balance.WithDrawn,
		balance.UserId,
	)

	return err
}
