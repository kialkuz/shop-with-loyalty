package db

import (
	"context"
	"fmt"
	"kialkuz/shop-with-loyalty/internal/infrastructure/storage"
	pkgErrors "kialkuz/shop-with-loyalty/pkg/errors"

	"kialkuz/shop-with-loyalty/internal/model/order"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type OrderRepository struct {
	db *storage.DB
}

func NewOrderRepository(db *storage.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) GetByNumber(ctx context.Context, number string) (*order.Order, error) {
	order := &order.Order{}
	err := r.db.QueryRow(ctx, func(row pgx.Row) error {
		return row.Scan(
			&order.ID,
			&order.UserId,
			&order.Number.Value,
		)
	}, "SELECT id, user_id, number FROM orders WHERE number = $1", number)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, pkgErrors.ErrNotFound
		}
		return nil, err
	}
	return order, nil
}

func (r *OrderRepository) GetByUserId(ctx context.Context, userId uuid.UUID) ([]order.Order, error) {
	rows, err := r.db.QueryRows(
		ctx,
		"SELECT id, user_id, number, status, accrual, uploaded_at FROM orders WHERE user_id = $1",
		userId,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, pkgErrors.ErrNotFound
		}
		return nil, err
	}

	defer rows.Close()

	var orders []order.Order

	for rows.Next() {
		o := order.Order{}

		if err := rows.Scan(
			&o.ID,
			&o.UserId,
			&o.Number.Value,
			&o.Status.Value,
			&o.Accrual,
			&o.UploadedAt,
		); err != nil {
			return nil, fmt.Errorf("scan order: %w", err)
		}

		orders = append(orders, o)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	if len(orders) == 0 {
		return nil, pkgErrors.ErrNotFound
	}

	return orders, nil
}

func (r *OrderRepository) Add(
	ctx context.Context,
	order order.Order,
) error {
	return r.db.Exec(
		ctx,
		"INSERT INTO orders (id, user_id, number, status, accrual, uploaded_at) VALUES ($1, $2, $3, $4, $5, $6)",
		order.ID,
		order.UserId,
		order.Number.Value,
		order.Status.Value,
		order.Accrual,
		order.UploadedAt,
	)
}
