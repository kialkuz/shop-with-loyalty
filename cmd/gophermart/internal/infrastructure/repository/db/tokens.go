package db

import (
	"context"
	"kialkuz/shop-with-loyalty/internal/infrastructure/storage"
	"kialkuz/shop-with-loyalty/internal/model"
	pkgErrors "kialkuz/shop-with-loyalty/pkg/errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type TokenRepository struct {
	db *storage.DB
}

func NewTokenRepository(db *storage.DB) *TokenRepository {
	return &TokenRepository{db: db}
}

func (r *TokenRepository) AddNewToken(
	ctx context.Context,
	token *model.Token,
) error {
	return r.db.Exec(
		ctx,
		"INSERT INTO active_tokens (jti, user_id, issued_at, expires_at) VALUES ($1, $2, $3, $4)",
		token.Jti,
		token.UserId,
		token.IssuedAt,
		token.ExpiresAt,
	)
}

func (r *TokenRepository) GetTokenByJti(ctx context.Context, jti uuid.UUID) (*model.Token, error) {
	token := &model.Token{}
	err := r.db.QueryRow(ctx, func(row pgx.Row) error {
		return row.Scan(
			&token.Jti,
			&token.UserId,
			&token.IssuedAt,
			&token.ExpiresAt,
		)
	}, `SELECT jti, user_id, issued_at, expires_at
		FROM active_tokens
		WHERE jti = $1`,
		jti,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, pkgErrors.ErrNotFound
		}
		return nil, err
	}
	return token, nil
}

func (r *TokenRepository) AddNewTokenTx(
	ctx context.Context,
	tx pgx.Tx,
	token *model.Token,
) error {
	_, err := tx.Exec(
		ctx,
		"INSERT INTO active_tokens (jti, user_id, issued_at, expires_at) VALUES ($1, $2, $3, $4)",
		token.Jti,
		token.UserId,
		token.IssuedAt,
		token.ExpiresAt,
	)

	return err
}
