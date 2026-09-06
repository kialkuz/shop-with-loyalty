package service

import (
	"context"
	"kialkuz/shop-with-loyalty/internal/domain/user/model"
	"kialkuz/shop-with-loyalty/internal/domain/user/service/mocks"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestCheckExistLogin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	service := NewUserService(mockRepo)

	ctx := context.Background()
	login := "test"

	mockRepo.EXPECT().CheckExistLogin(ctx, login).Return(true, nil)

	result, err := service.CheckExistLogin(ctx, login)

	assert.NoError(t, err)
	assert.True(t, result)
}

func TestGetByLogin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	service := NewUserService(mockRepo)

	ctx := context.Background()

	user := &model.User{
		ID:       uuid.New(),
		Login:    "test",
		Password: "12345",
	}

	mockRepo.EXPECT().GetByLogin(ctx, gomock.Any()).Return(user, nil)

	result, err := service.GetByLogin(ctx, user.Login)

	assert.NoError(t, err)
	assert.NotEmpty(t, result)
}

func TestAddNewUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	service := NewUserService(mockRepo)

	ctx := context.Background()

	user := model.User{
		ID:       uuid.New(),
		Login:    "test",
		Password: "12345",
	}

	mockRepo.EXPECT().AddNewUser(ctx, gomock.Any()).Return(nil)

	err := service.AddNewUser(ctx, user)

	assert.NoError(t, err)
}
