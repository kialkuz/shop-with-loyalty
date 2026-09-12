package db

import (
	"context"
	"fmt"
	"kialkuz/shop-with-loyalty/internal/infrastructure"
	pkgErrors "kialkuz/shop-with-loyalty/pkg/errors"

	orderModel "kialkuz/shop-with-loyalty/internal/modules/order/model"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type OrderRepository struct {
	db *infrastructure.DB
}

func NewOrderRepository(db *infrastructure.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) GetByNumber(ctx context.Context, number string) (*orderModel.Order, error) {
	order := &orderModel.Order{}
	err := r.db.QueryRow(ctx, func(row pgx.Row) error {
		return row.Scan(
			&order.ID,
			&order.UserID,
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

func (r *OrderRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]orderModel.Order, error) {
	rows, err := r.db.QueryRows(
		ctx,
		"SELECT id, user_id, number, status, accrual, uploaded_at FROM orders WHERE user_id = $1",
		userID,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, pkgErrors.ErrNotFound
		}
		return nil, err
	}

	defer rows.Close()

	return r.scanRows(rows)
}

func (r *OrderRepository) GetOrdersForGetAccrual(ctx context.Context) (map[string]orderModel.Order, error) {
	queryRows, err := r.db.QueryRows(
		ctx,
		"SELECT id, user_id, number, status, accrual, uploaded_at FROM orders WHERE status = $1 LIMIT 30",
		orderModel.StatusNew,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, pkgErrors.ErrNotFound
		}
		return nil, err
	}

	defer queryRows.Close()

	rows, err := r.scanRows(queryRows)
	if err != nil {
		return nil, err
	}

	orders := make(map[string]orderModel.Order)
	for _, row := range rows {
		orders[row.Number.Value] = row
	}

	return orders, nil
}

func (r *OrderRepository) scanRows(rows pgx.Rows) ([]orderModel.Order, error) {
	var orders []orderModel.Order

	for rows.Next() {
		o := orderModel.Order{
			Status: &orderModel.Status{},
		}

		if err := rows.Scan(
			&o.ID,
			&o.UserID,
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
	order orderModel.Order,
) error {
	return r.db.Exec(
		ctx,
		"INSERT INTO orders (id, user_id, number, status, accrual, uploaded_at) VALUES ($1, $2, $3, $4, $5, $6)",
		order.ID,
		order.UserID,
		order.Number.Value,
		order.Status.Value,
		order.Accrual,
		order.UploadedAt,
	)
}

func (r *OrderRepository) UpdateTx(
	ctx context.Context,
	tx pgx.Tx,
	order orderModel.Order,
) error {
	_, err := tx.Exec(
		ctx,
		"UPDATE orders SET status = $1, accrual = $2 WHERE number = $3",
		order.Status.Value,
		order.Accrual,
		order.Number.Value,
	)

	return err
}
