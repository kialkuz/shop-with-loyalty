package service

import (
	"testing"

	"kialkuz/shop-with-loyalty/internal/service/mocks"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestHashPassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	authMockRepo := mocks.NewMockAuthRepository(ctrl)
	userMockRepo := mocks.NewMockUserRepository(ctrl)
	service := NewAuthService(authMockRepo, userMockRepo)

	mockRepo.EXPECT().HashPassword("12345").Return(true, nil)

	isExistLogin, err := service.HashPassword("12345")

	assert.NoError(t, err)
	assert.True(t, isExistLogin)
}

// func TestExistLogin(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	mockRepo := mocks.NewMockUserRepository(ctrl)
// 	service := NewUserService(mockRepo)

// 	mockRepo.EXPECT().GetByLogin(gomock.Any(), "test").Return(&model.User{
// 		ID:       uuid.New(),
// 		Login:    "test",
// 		Password: "12345",
// 	}, nil)

// 	user, err := service.GetByLogin(context.TODO(), "test")

// 	assert.NoError(t, err)
// 	assert.NotEmpty(t, user)
// }

// func TestAddNewUser(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	mockRepo := mocks.NewMockUserRepository(ctrl)
// 	service := NewUserService(mockRepo)

// 	user := model.User{
// 		ID:       uuid.New(),
// 		Login:    "test",
// 		Password: "12345",
// 	}

// 	mockRepo.EXPECT().AddNewUser(gomock.Any(), user).Return(nil)

// 	err := service.AddNewUser(context.TODO(), user)

// 	assert.NoError(t, err)
// }
