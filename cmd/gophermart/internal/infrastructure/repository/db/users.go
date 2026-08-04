package db

import (
	"kialkuz/gophermart/internal/infrastructure/storage"
)

type UsersRepository struct {
	db *storage.DB
}

func NewUsersRepository(db *storage.DB) *UsersRepository {
	return &UsersRepository{db: db}
}
