package db

import (
	"context"
	"fmt"
	drawalDto "kialkuz/shop-with-loyalty/internal/domain/drawal/dto"
	modelDrawal "kialkuz/shop-with-loyalty/internal/domain/drawal/model"
	"kialkuz/shop-with-loyalty/internal/infrastructure"
	pkgErrors "kialkuz/shop-with-loyalty/pkg/errors"

	"github.com/google/uuid"

	"github.com/jackc/pgx/v5"
)

type DrawalRepository struct {
	db *infrastructure.DB
}

func NewDrawalRepository(db *infrastructure.DB) *DrawalRepository {
	return &DrawalRepository{db: db}
}

func (r *DrawalRepository) AddTx(ctx context.Context, tx pgx.Tx, drawal modelDrawal.Drawal) error {
	_, err := tx.Exec(
		ctx,
		"INSERT INTO drawals (id, user_id, order_number, sum, processed_at) VALUES ($1, $2, $3, $4, $5)",
		drawal.ID,
		drawal.UserID,
		drawal.OrderNumber.Value,
		drawal.Sum,
		drawal.ProcessedAt,
	)

	return err
}

func (r *DrawalRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]drawalDto.UserDrawal, error) {
	rows, err := r.db.QueryRows(
		ctx,
		`SELECT order_number, sum, processed_at
		FROM drawals
		WHERE user_id = $1
		ORDER BY processed_at DESC`,
		userID,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, pkgErrors.ErrNotFound
		}
		return nil, err
	}

	defer rows.Close()

	var drawals []drawalDto.UserDrawal

	for rows.Next() {
		o := drawalDto.UserDrawal{}

		if err := rows.Scan(
			&o.OrderNumber,
			&o.Sum,
			&o.ProcessedAt,
		); err != nil {
			return nil, fmt.Errorf("scan order: %w", err)
		}

		drawals = append(drawals, o)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	if len(drawals) == 0 {
		return nil, pkgErrors.ErrNotFound
	}

	return drawals, nil
}
