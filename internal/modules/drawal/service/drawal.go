package service

import (
	"context"
	drawalDto "kialkuz/shop-with-loyalty/internal/modules/drawal/dto"
	"kialkuz/shop-with-loyalty/internal/modules/drawal/model"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

//go:generate go run go.uber.org/mock/mockgen -source=drawal.go -destination=mocks/drawal_mock.go -package=mocks -typed
type DrawalRepository interface {
	AddTx(ctx context.Context, tx pgx.Tx, drawal model.Drawal) error
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]drawalDto.UserDrawal, error)
}

type DrawalService struct {
	repository DrawalRepository
}

func NewDrawalService(repository DrawalRepository) *DrawalService {
	return &DrawalService{repository: repository}
}

func (s *DrawalService) AddTx(ctx context.Context, tx pgx.Tx, drawal model.Drawal) error {
	return s.repository.AddTx(ctx, tx, drawal)
}

func (s *DrawalService) GetByUserID(ctx context.Context, userID uuid.UUID) ([]drawalDto.UserDrawal, error) {
	return s.repository.GetByUserID(ctx, userID)
}
