package service

import (
	"context"
	"kialkuz/shop-with-loyalty/internal/auth/model"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type TokenService struct {
	repository TokenRepository
}

func NewTokenService(repository TokenRepository) *TokenService {
	return &TokenService{repository: repository}
}

func (s TokenService) AddNewToken(ctx context.Context, token model.Token) error {
	return s.repository.AddNewToken(ctx, token)
}

func (s TokenService) AddTx(ctx context.Context, tx pgx.Tx, token model.Token) error {
	return s.repository.AddTx(ctx, tx, token)
}

func (s TokenService) GetTokenByJti(ctx context.Context, jti uuid.UUID) (*model.Token, error) {
	return s.repository.GetTokenByJti(ctx, jti)
}
