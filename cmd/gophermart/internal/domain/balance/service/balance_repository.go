package service

import (
	"context"
	"kialkuz/shop-with-loyalty/internal/domain/balance/model"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

//go:generate go run go.uber.org/mock/mockgen -destination=mocks/token_repository_mock.go -package=mocks -typed
type BalanceRepository interface {
	GetByUserId(ctx context.Context, userId uuid.UUID) (*model.Balance, error)
	GetByUserIdWithBlockForUpdate(ctx context.Context, tx pgx.Tx, userId uuid.UUID) (*model.Balance, error)
	GetByUsersId(ctx context.Context, usersId []uuid.UUID) (map[uuid.UUID]model.Balance, error)
	AddTx(ctx context.Context, tx pgx.Tx, balance model.Balance) error
	UpdateTx(ctx context.Context, tx pgx.Tx, balance model.Balance) error
}
