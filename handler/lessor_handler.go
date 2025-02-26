package handler

import (
	"net/http"
	"rent-video-game/middleware"
	"rent-video-game/model"
	"rent-video-game/usecase"
	"rent-video-game/utils"

	"github.com/labstack/echo/v4"
)

type LessorHandler struct {
	lessorUsecase *usecase.LessorUsecase
}

func NewLessorHandler(lessorUsecase *usecase.LessorUsecase) *LessorHandler {
	return &LessorHandler{lessorUsecase: lessorUsecase}
}

func (u *LessorHandler) LessorRoutes(e *echo.Echo) {
	e.POST("/lessor/register", middleware.UserAuthMiddleware()(u.RegisterLessor))
	e.GET("/lessor/:lessor_id", middleware.UserAuthMiddleware()(u.GetLessorByID))
	e.PUT("/lessor/:lessor_id", middleware.UserAuthMiddleware()(u.UpdateLessor))
	e.DELETE("/lessor/:lessor_id", middleware.UserAuthMiddleware()(u.DeleteLessor))
}

func (u *LessorHandler) RegisterLessor(c echo.Context) error {
	var lessor *model.Lessors
	if err := c.Bind(&lessor); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	userID, err := UserToken(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if lessor.UserID != userID {
		return echo.NewHTTPError(http.StatusForbidden, "forbidden access")
	}

	lessor.UserID = userID

	lessor, err = u.lessorUsecase.RegisterLessor(lessor)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	response := model.LessorRegisterResponse{
		Message: "success register lessor",
	}

	response.Data.LessorID = lessor.LessorID
	response.Data.Name = lessor.Name
	response.Data.Location = lessor.Location

	return c.JSON(http.StatusOK, response)
}

func (u *LessorHandler) GetLessorByID(c echo.Context) error {
	lessorID := c.Param("lessor_id")
	id := utils.StringToInt(lessorID)

	userID, err := UserToken(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	lessor, err := u.lessorUsecase.GetLessorByID(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if lessor.UserID != userID {
		return echo.NewHTTPError(http.StatusForbidden, "forbidden access")
	}

	response := model.LessorRegisterResponse{
		Message: "success get lessor by id",
	}

	response.Data.LessorID = lessor.LessorID
	response.Data.Name = lessor.Name
	response.Data.Location = lessor.Location

	return c.JSON(http.StatusOK, response)
}

func (u *LessorHandler) UpdateLessor(c echo.Context) error {
	var lessor *model.Lessors
	if err := c.Bind(&lessor); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	lessorID := c.Param("lessor_id")
	id := utils.StringToInt(lessorID)

	userID, err := UserToken(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if lessor.UserID != userID {
		return echo.NewHTTPError(http.StatusForbidden, "forbidden access")
	}

	lessor.UserID = userID

	lessor, err = u.lessorUsecase.UpdateLessor(id, lessor)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	response := model.LessorRegisterResponse{
		Message: "success update lessor",
	}

	response.Data.LessorID = lessor.LessorID
	response.Data.Name = lessor.Name
	response.Data.Location = lessor.Location

	return c.JSON(http.StatusOK, response)
}

func (u *LessorHandler) DeleteLessor(c echo.Context) error {
	lessorID := c.Param("lessor_id")
	id := utils.StringToInt(lessorID)

	userID, err := UserToken(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	lessor, err := u.lessorUsecase.DeleteLessor(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if lessor.UserID != userID {
		return echo.NewHTTPError(http.StatusForbidden, "forbidden access")
	}

	response := model.LessorRegisterResponse{
		Message: "success delete lessor",
	}

	response.Data.LessorID = lessor.LessorID
	response.Data.Name = lessor.Name
	response.Data.Location = lessor.Location

	return c.JSON(http.StatusOK, response)
}
