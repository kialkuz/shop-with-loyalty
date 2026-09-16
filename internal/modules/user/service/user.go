package service

import (
	"context"
	"kialkuz/shop-with-loyalty/internal/modules/user/model"

	"github.com/jackc/pgx/v5"
)

//go:generate go run go.uber.org/mock/mockgen -source=user.go -destination=mocks/user_repository_mock.go -package=mocks -typed
type UserRepository interface {
	CheckExistLogin(ctx context.Context, login string) (bool, error)
	AddNewUser(ctx context.Context, user model.User) error
	AddTx(ctx context.Context, tx pgx.Tx, user model.User) error
	GetByLogin(ctx context.Context, login string) (*model.User, error)
}

type UserService struct {
	repository UserRepository
}

func NewUserService(repository UserRepository) *UserService {
	return &UserService{repository: repository}
}

func (s *UserService) CheckExistLogin(ctx context.Context, login string) (bool, error) {
	return s.repository.CheckExistLogin(ctx, login)
}

func (s *UserService) GetByLogin(ctx context.Context, login string) (*model.User, error) {
	user, err := s.repository.GetByLogin(ctx, login)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) AddNewUser(ctx context.Context, user model.User) error {
	return s.repository.AddNewUser(ctx, user)
}

func (s *UserService) AddTx(ctx context.Context, tx pgx.Tx, user model.User) error {
	return s.repository.AddTx(ctx, tx, user)
}
