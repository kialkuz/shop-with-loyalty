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

func (s *BalanceService) GetByUserID(ctx context.Context, userID uuid.UUID) (*model.Balance, error) {
	return s.repository.GetByUserID(ctx, userID)
}

func (s *BalanceService) GetByUserIDWithBlockForUpdate(
	ctx context.Context,
	tx pgx.Tx,
	userID uuid.UUID) (*model.Balance, error) {
	return s.repository.GetByUserIDWithBlockForUpdate(ctx, tx, userID)
}

func (s *BalanceService) GetByUsersID(ctx context.Context, usersID []uuid.UUID) (map[uuid.UUID]model.Balance, error) {
	return s.repository.GetByUsersID(ctx, usersID)
}

func (s *BalanceService) UpdateTx(ctx context.Context, tx pgx.Tx, balance model.Balance) error {
	return s.repository.UpdateTx(ctx, tx, balance)
}

func (s *BalanceService) AddTx(ctx context.Context, tx pgx.Tx, user model.Balance) error {
	return s.repository.AddTx(ctx, tx, user)
}
