package service

import (
	"context"
	balanceModel "kialkuz/shop-with-loyalty/internal/domain/balance/model"
	balanceServ "kialkuz/shop-with-loyalty/internal/domain/balance/service"
	mocksBalanceRepo "kialkuz/shop-with-loyalty/internal/domain/balance/service/mocks"
	drawalModel "kialkuz/shop-with-loyalty/internal/domain/drawal/model"
	drawalServ "kialkuz/shop-with-loyalty/internal/domain/drawal/service"
	mocksDrawalRepo "kialkuz/shop-with-loyalty/internal/domain/drawal/service/mocks"
	orderModel "kialkuz/shop-with-loyalty/internal/domain/order/model"
	"kialkuz/shop-with-loyalty/internal/domain/order/service/mocks"
	infrastructureMocks "kialkuz/shop-with-loyalty/internal/infrastructure/mocks"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestGetBNumber(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBalanceRepo := mocksBalanceRepo.NewMockBalanceRepository(ctrl)
	balanceService := balanceServ.NewBalanceService(mockBalanceRepo)

	mockDrawalRepo := mocksDrawalRepo.NewMockDrawalRepository(ctrl)
	drawalService := drawalServ.NewDrawalService(mockDrawalRepo)

	mockTransaction := infrastructureMocks.NewMockTransaction(ctrl)

	mockRepo := mocks.NewMockOrderRepository(ctrl)
	service := NewOrderService(mockTransaction, mockRepo, balanceService, drawalService)

	ctx := context.Background()

	order := getOrder(uuid.New(), "12345678903")

	mockRepo.EXPECT().GetByNumber(ctx, order.Number).Return(&order, nil)

	result, err := service.GetByNumber(ctx, order.Number.Value)

	assert.NoError(t, err)
	assert.NotEmpty(t, result)
}

func TestGetByUserID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBalanceRepo := mocksBalanceRepo.NewMockBalanceRepository(ctrl)
	balanceService := balanceServ.NewBalanceService(mockBalanceRepo)

	mockDrawalRepo := mocksDrawalRepo.NewMockDrawalRepository(ctrl)
	drawalService := drawalServ.NewDrawalService(mockDrawalRepo)

	mockTransaction := infrastructureMocks.NewMockTransaction(ctrl)

	mockRepo := mocks.NewMockOrderRepository(ctrl)
	service := NewOrderService(mockTransaction, mockRepo, balanceService, drawalService)

	ctx := context.Background()

	order := getOrder(uuid.New(), "12345678903")

	mockRepo.EXPECT().GetByUserID(ctx, order.UserID).Return([]orderModel.Order{
		order,
	}, nil)

	result, err := service.GetByUserID(ctx, order.UserID)

	assert.NoError(t, err)
	assert.NotEmpty(t, result)
}

func TestAdd(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBalanceRepo := mocksBalanceRepo.NewMockBalanceRepository(ctrl)
	balanceService := balanceServ.NewBalanceService(mockBalanceRepo)

	mockDrawalRepo := mocksDrawalRepo.NewMockDrawalRepository(ctrl)
	drawalService := drawalServ.NewDrawalService(mockDrawalRepo)

	mockTransaction := infrastructureMocks.NewMockTransaction(ctrl)

	mockRepo := mocks.NewMockOrderRepository(ctrl)
	service := NewOrderService(mockTransaction, mockRepo, balanceService, drawalService)

	ctx := context.Background()

	order := getOrder(uuid.New(), "12345678903")

	mockRepo.EXPECT().Add(ctx, order).Return(nil)

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

func TestWithDraw(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBalanceRepo := mocksBalanceRepo.NewMockBalanceRepository(ctrl)
	balanceService := balanceServ.NewBalanceService(mockBalanceRepo)

	mockDrawalRepo := mocksDrawalRepo.NewMockDrawalRepository(ctrl)
	drawalService := drawalServ.NewDrawalService(mockDrawalRepo)

	mockTransaction := infrastructureMocks.NewMockTransaction(ctrl)
	mockTx := infrastructureMocks.NewMockTx(ctrl)

	mockRepo := mocks.NewMockOrderRepository(ctrl)
	service := NewOrderService(mockTransaction, mockRepo, balanceService, drawalService)

	ctx := context.Background()
	userID := uuid.New()

	drawal := drawalModel.Drawal{
		ID:     uuid.New(),
		UserID: userID,
		Sum:    50,
	}

	balance := &balanceModel.Balance{
		ID:        uuid.New(),
		UserID:    userID,
		Current:   100,
		WithDrawn: 0,
	}

	expectedBalance := *balance
	expectedBalance.Current = 50
	expectedBalance.WithDrawn = 50

	mockTransaction.
		EXPECT().
		RunTransaction(ctx, gomock.All()).
		DoAndReturn(func(ctx context.Context, fn func(pgx.Tx) error) error {
			return fn(mockTx)
		})

	mockBalanceRepo.EXPECT().GetByUserIDWithBlockForUpdate(ctx, mockTx, userID).Return(balance, nil)

	mockBalanceRepo.
		EXPECT().
		UpdateTx(
			ctx,
			mockTx,
			expectedBalance,
		).
		Return(nil)

	mockDrawalRepo.
		EXPECT().
		AddTx(
			ctx,
			mockTx,
			drawal,
		).
		Return(nil)

	err := service.WithDraw(ctx, userID, drawal)

	assert.NoError(t, err)
}
