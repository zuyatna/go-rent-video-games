package tests

import (
	"rent-video-game/mocks"
	"rent-video-game/model"
	"rent-video-game/usecase"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestRegisterUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockIUserRepository(ctrl)
	userUsecase := usecase.NewUserUsecase(mockRepo)

	testUser := &model.Users{
		UserID:   uuid.New(),
		Name:     "test",
		Email:    "test@example.com",
		Password: "test_password",
		Address:  "address",
		Amount:   0,
	}

	mockRepo.EXPECT().RegisterUser(gomock.Any()).Return(testUser, nil)

	result, err := userUsecase.RegisterUser(testUser)
	if err != nil {
		t.Error(err)
	}

	assert.NoError(t, err)
	assert.Equal(t, testUser, result)
}
