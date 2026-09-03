package service

import (
	"context"
	"kialkuz/shop-with-loyalty/internal/domain/balance/model"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type BalanceService struct {
	repository BalanceRepository
}

func NewBalanceService(repository BalanceRepository) *BalanceService {
	return &BalanceService{repository: repository}
}

func (s *BalanceService) GetByUserId(ctx context.Context, userId uuid.UUID) (*model.Balance, error) {
	return s.repository.GetByUserId(ctx, userId)
}

func (s *BalanceService) GetByUserIdWithBlockForUpdate(
	ctx context.Context,
	tx pgx.Tx,
	userId uuid.UUID) (*model.Balance, error) {
	return s.repository.GetByUserIdWithBlockForUpdate(ctx, tx, userId)
}

func (s *BalanceService) GetByUsersId(ctx context.Context, usersId []uuid.UUID) (map[uuid.UUID]model.Balance, error) {
	return s.repository.GetByUsersId(ctx, usersId)
}

func (s *BalanceService) UpdateTx(ctx context.Context, tx pgx.Tx, balance model.Balance) error {
	return s.repository.UpdateTx(ctx, tx, balance)
}

func (s *BalanceService) AddTx(ctx context.Context, tx pgx.Tx, user model.Balance) error {
	return s.repository.AddTx(ctx, tx, user)
}
