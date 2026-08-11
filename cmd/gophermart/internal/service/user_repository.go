package service

import (
	"context"
	"kialkuz/shop-with-loyalty/internal/model"

	"github.com/jackc/pgx/v5"
)

//go:generate go run go.uber.org/mock/mockgen -destination=mocks/user_repository_mock.go -package=mocks -typed kialkuz/shop-with-loyalty/internal/contracts/repository UserRepository
type UserRepository interface {
	CheckExistLogin(ctx context.Context, login string) (bool, error)
	AddNewUser(ctx context.Context, user model.User) error
	AddNewUserTx(ctx context.Context, tx pgx.Tx, user model.User) error
	GetByLogin(ctx context.Context, login string) (*model.User, error)
	GetByToken(ctx context.Context, token string) (*model.User, error)
}
