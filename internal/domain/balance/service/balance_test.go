package service

import (
	"context"
	"kialkuz/shop-with-loyalty/internal/domain/balance/model"
	"kialkuz/shop-with-loyalty/internal/domain/balance/service/mocks"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestGetByUserID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockBalanceRepository(ctrl)
	service := NewBalanceService(mockRepo)

	ctx := context.Background()
	userID := uuid.New()

	mockRepo.EXPECT().GetByUserID(ctx, gomock.Any()).Return(&model.Balance{
		ID:        uuid.New(),
		UserID:    userID,
		Current:   0,
		WithDrawn: 0,
	}, nil)

	result, err := service.GetByUserID(ctx, userID)

	assert.NoError(t, err)
	assert.NotEmpty(t, result)
}
