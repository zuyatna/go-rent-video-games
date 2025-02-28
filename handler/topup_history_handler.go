package handler

import (
	"net/http"
	"rent-video-game/middleware"
	"rent-video-game/model"
	"rent-video-game/usecase"
	"rent-video-game/utils"

	"github.com/labstack/echo/v4"
)

type TopupHistoryHandler struct {
	TopupHistoryUsecase *usecase.TopupHistoryUsecase
}

func NewTopupHistoryHandler(topupHistoryUsecase *usecase.TopupHistoryUsecase) *TopupHistoryHandler {
	return &TopupHistoryHandler{TopupHistoryUsecase: topupHistoryUsecase}
}

func (u *TopupHistoryHandler) TopupHistoryRoutes(e *echo.Echo) {
	e.GET("user/topup-history/:topup_history_id", middleware.UserAuthMiddleware()(u.GetTopupHistoryByID))
	e.GET("user/topup-histories", middleware.UserAuthMiddleware()(u.GetAllTopupHistory))
}

func (u *TopupHistoryHandler) GetTopupHistoryByID(c echo.Context) error {
	topupHistoryID := c.Param("topup_history_id")
	id := utils.StringToInt(topupHistoryID)

	topupHistory, err := u.TopupHistoryUsecase.GetTopupHistoryByID(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	topupHistoryData := model.TopupHistoryData{
		TopupHistoryID: topupHistory.TopupHistoryID,
		PaymentID:      topupHistory.PaymentID,
		Amount:         topupHistory.Amount,
		CreatedAt:      topupHistory.CreatedAt,
	}

	response := model.TopupHistoryResponse{
		Message: "success get topup history by user id",
		Data:    []model.TopupHistoryData{topupHistoryData},
	}

	return c.JSON(http.StatusOK, response)
}

func (u *TopupHistoryHandler) GetAllTopupHistory(c echo.Context) error {
	userID, err := UserToken(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	topupHistory, err := u.TopupHistoryUsecase.GetAllTopupHistoryByUser(userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	var topupHistoryData []model.TopupHistoryData
	for _, value := range *topupHistory {
		topupHistoryData = append(topupHistoryData, model.TopupHistoryData{
			TopupHistoryID: value.TopupHistoryID,
			PaymentID:      value.PaymentID,
			Amount:         value.Amount,
			CreatedAt:      value.CreatedAt,
		})
	}

	response := model.TopupHistoryResponse{
		Message: "success get all topup history",
		Data:    topupHistoryData,
	}

	return c.JSON(http.StatusOK, response)
}
