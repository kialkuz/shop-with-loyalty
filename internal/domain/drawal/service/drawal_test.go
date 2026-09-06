package service

import (
	"context"
	drawalDto "kialkuz/shop-with-loyalty/internal/domain/drawal/dto"
	"kialkuz/shop-with-loyalty/internal/domain/drawal/service/mocks"
	orderModel "kialkuz/shop-with-loyalty/internal/domain/order/model"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestGetByUserID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockDrawalRepository(ctrl)
	service := NewDrawalService(mockRepo)

	ctx := context.Background()
	userID := uuid.New()

	mockRepo.EXPECT().GetByUserID(ctx, gomock.Any()).Return([]drawalDto.UserDrawal{
		{
			OrderNumber: orderModel.NewNumber("12345678903"),
			Sum:         100,
			ProcessedAt: time.Now(),
		},
	}, nil)

	result, err := service.GetByUserID(ctx, userID)

	assert.NoError(t, err)
	assert.NotEmpty(t, result)
}
