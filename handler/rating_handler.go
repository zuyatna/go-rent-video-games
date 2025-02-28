package handler

import (
	"errors"
	"net/http"
	"rent-video-game/middleware"
	"rent-video-game/model"
	"rent-video-game/usecase"
	"rent-video-game/utils"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type RatingHandler struct {
	ratingUsecase *usecase.RatingUsecase
}

func NewRatingHandler(ratingUsecase *usecase.RatingUsecase) *RatingHandler {
	return &RatingHandler{ratingUsecase}
}

func (h *RatingHandler) RatingRoutes(e *echo.Echo) {
	e.POST("/user/rating", middleware.UserAuthMiddleware()(h.CreateRating))
	e.GET("/lessor/rating/product/:product_id", middleware.UserAuthMiddleware()(h.GetAllRatingByProduct))
}

func (h *RatingHandler) CreateRating(c echo.Context) error {
	var ratingReq *model.RatingRequest
	if err := c.Bind(&ratingReq); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	userID, err := UserToken(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	rating := &model.Ratings{
		ProductID: ratingReq.ProductID,
		Review:    ratingReq.Review,
		Stars:     ratingReq.Stars,
	}

	if rating.UserID != uuid.Nil && rating.UserID != userID {
		return echo.NewHTTPError(http.StatusForbidden, "forbidden access")
	}

	existingRating, err := h.ratingUsecase.GetRatingByUserAndProduct(userID, rating.ProductID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if existingRating != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "you have already rated this product")
	}

	rating.UserID = userID
	rating, err = h.ratingUsecase.CreateRating(rating)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	ratingData := model.RatingData{
		RatingID:  rating.RatingID,
		ProductID: rating.ProductID,
		Review:    rating.Review,
		Stars:     rating.Stars,
	}

	response := model.RatingResponse{
		Message: "success create rating",
		Data:    []model.RatingData{ratingData},
	}

	return c.JSON(http.StatusOK, response)
}

func (h *RatingHandler) GetAllRatingByProduct(c echo.Context) error {
	productID := c.Param("product_id")
	id := utils.StringToInt(productID)

	ratings, err := h.ratingUsecase.GetAllRatingByProduct(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	var ratingData []model.RatingData
	for _, rating := range ratings {
		ratingData = append(ratingData, model.RatingData{
			RatingID:  rating.RatingID,
			ProductID: rating.ProductID,
			Review:    rating.Review,
			Stars:     rating.Stars,
		})
	}

	response := model.RatingResponse{
		Message: "success get all rating by product",
		Data:    ratingData,
	}

	return c.JSON(http.StatusOK, response)
}
