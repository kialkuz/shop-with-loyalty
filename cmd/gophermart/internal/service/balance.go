package service

import (
	"kialkuz/shop-with-loyalty/internal/contracts/repository"
)

type BalanceService struct {
	repository repository.BalanceRepository
}

func NewBalanceService(repository repository.BalanceRepository) *BalanceService {
	return &BalanceService{repository: repository}
}
