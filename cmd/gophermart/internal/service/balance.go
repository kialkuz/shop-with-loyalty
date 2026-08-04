package service

import (
	"kialkuz/gophermart/internal/infrastructure/repository/db"
)

type BalanceService struct {
	repository *db.BalanceRepository
}

func NewBalanceService(repository *db.BalanceRepository) *BalanceService {
	return &BalanceService{repository: repository}
}
