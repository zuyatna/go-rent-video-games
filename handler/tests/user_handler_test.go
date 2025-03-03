package tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"rent-video-game/handler"
	"rent-video-game/mocks"
	"rent-video-game/model"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestRegisterUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserUsecase := mocks.NewMockIUserUsecase(ctrl)
	mockTopupHistoryUsecase := mocks.NewMockITopupHistoryUsecase(ctrl)

	handler := handler.NewUserHandlerWithInterface(mockUserUsecase, mockTopupHistoryUsecase)

	e := echo.New()

	reqJSON := `{
        "name": "test",
        "email": "test@example.com",
        "password": "test_password",
        "address": "address"
    }`

	req := httptest.NewRequest(http.MethodPost, "/user/register", strings.NewReader(reqJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	userID := uuid.New()
	expectedUser := &model.Users{
		UserID:   userID,
		Name:     "test",
		Email:    "test@example.com",
		Password: "test_password",
		Address:  "address",
	}

	mockUserUsecase.EXPECT().GetUserByEmail(gomock.Any()).Return(nil, fmt.Errorf("user not found"))
	mockUserUsecase.EXPECT().RegisterUser(gomock.Any()).Return(expectedUser, nil)

	err := handler.RegisterUserInterface(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response model.RegisterResponse
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "success register user", response.Message)
	assert.Equal(t, expectedUser.UserID, response.Data.UserID)
	assert.Equal(t, expectedUser.Email, response.Data.Email)
	assert.Equal(t, expectedUser.Name, response.Data.Name)
	assert.Equal(t, http.StatusOK, rec.Code)
}
