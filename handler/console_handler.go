package handler

import (
	"net/http"
	"rent-video-game/model"
	"rent-video-game/usecase"

	"github.com/labstack/echo/v4"
)

type ConsoleHandler struct {
	consoleUsecase *usecase.ConsoleUsecase
}

func NewConsoleHandler(consoleUsecase *usecase.ConsoleUsecase) *ConsoleHandler {
	return &ConsoleHandler{consoleUsecase: consoleUsecase}
}

func (u *ConsoleHandler) ConsoleRoutes(e *echo.Echo) {
	e.GET("/consoles", u.GetAllConsole)
}

func (u *ConsoleHandler) GetAllConsole(c echo.Context) error {
	console, err := u.consoleUsecase.GetAllConsole()
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	var consoleData []model.ConsoleData
	for _, value := range console {
		consoleData = append(consoleData, model.ConsoleData{
			ConsoleID: value.ConsoleID,
			Name:      value.Name,
		})

	}

	consoleResponse := model.ConsoleResponse{
		Message: "success get all console",
		Data:    consoleData,
	}

	return c.JSON(http.StatusOK, consoleResponse)
}
