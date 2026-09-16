package service

import (
	"context"
	orderModel "kialkuz/shop-with-loyalty/internal/modules/order/model"
	"kialkuz/shop-with-loyalty/internal/modules/order/service/mocks"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestGetByNumber(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockOrderRepository(ctrl)
	service := NewOrderService(mockRepo)

	ctx := context.Background()

	order := getOrder(uuid.New(), "12345678903")

	mockRepo.EXPECT().GetByNumber(ctx, gomock.Any()).Return(&order, nil)

	result, err := service.GetByNumber(ctx, order.Number.Value)

	assert.NoError(t, err)
	assert.NotEmpty(t, result)
}

func TestGetByUserID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockOrderRepository(ctrl)
	service := NewOrderService(mockRepo)

	ctx := context.Background()

	order := getOrder(uuid.New(), "12345678903")

	mockRepo.EXPECT().GetByUserID(ctx, gomock.Any()).Return([]orderModel.Order{
		order,
	}, nil)

	result, err := service.GetByUserID(ctx, order.UserID)

	assert.NoError(t, err)
	assert.NotEmpty(t, result)
}

func TestAdd(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockOrderRepository(ctrl)
	service := NewOrderService(mockRepo)

	ctx := context.Background()

	order := getOrder(uuid.New(), "12345678903")

	mockRepo.EXPECT().Add(ctx, gomock.Any()).Return(nil)

	err := service.Add(ctx, order)

	assert.NoError(t, err)
}

func getOrder(userID uuid.UUID, number string) orderModel.Order {
	status, _ := orderModel.NewStatus(orderModel.StatusNew)

	return orderModel.Order{
		ID:         uuid.New(),
		UserID:     userID,
		Number:     orderModel.NewNumber(number),
		Status:     status,
		Accrual:    nil,
		UploadedAt: time.Now(),
	}
}
