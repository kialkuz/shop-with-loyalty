package service

import (
	"context"
	"kialkuz/shop-with-loyalty/internal/model"

	"github.com/google/uuid"
)

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

func (s *UserService) GetUserIdByToken(ctx context.Context, token string) (*uuid.UUID, error) {
	user, err := s.repository.GetByToken(ctx, token)
	if err != nil {
		return nil, err
	}

	userId := user.ID

	return &userId, nil
}

func (s *UserService) AddNewUser(ctx context.Context, user model.User) error {
	return s.repository.AddNewUser(ctx, user)
}
