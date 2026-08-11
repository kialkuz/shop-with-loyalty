package service

import (
	"context"
	"kialkuz/shop-with-loyalty/internal/model"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

//go:generate go run go.uber.org/mock/mockgen -destination=mocks/token_repository_mock.go -package=mocks -typed
type TokenRepository interface {
	AddNewToken(ctx context.Context, token *model.Token) error
	AddNewTokenTx(ctx context.Context, tx pgx.Tx, token *model.Token) error
	GetTokenByJti(ctx context.Context, jti uuid.UUID) (*model.Token, error)
}
