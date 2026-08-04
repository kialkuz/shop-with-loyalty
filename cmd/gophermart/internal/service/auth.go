package service

import (
	"kialkuz/gophermart/internal/infrastructure/repository/db"
)

type AuthService struct {
	repository *db.UsersRepository
}

func NewAuthService(repository *db.UsersRepository) *AuthService {
	return &AuthService{repository: repository}
}
