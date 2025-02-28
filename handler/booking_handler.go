package handler

import (
	"fmt"
	"net/http"
	"rent-video-game/middleware"
	"rent-video-game/model"
	"rent-video-game/usecase"
	"rent-video-game/utils"

	"github.com/labstack/echo/v4"
)

type BookingHandler struct {
	bookingUsecase *usecase.BookingUsecase
	userUsecase    *usecase.UserUsecase
	productUsecase *usecase.ProductUsecase
	lessorUsecase  *usecase.LessorUsecase
}

func NewBookingHandler(
	bookingUsecase *usecase.BookingUsecase,
	userUsecase *usecase.UserUsecase,
	productUsecase *usecase.ProductUsecase,
	lessorUsecase *usecase.LessorUsecase,
) *BookingHandler {
	return &BookingHandler{
		bookingUsecase: bookingUsecase,
		userUsecase:    userUsecase,
		productUsecase: productUsecase,
		lessorUsecase:  lessorUsecase,
	}
}

func (u *BookingHandler) BookingRoutes(e *echo.Echo) {
	e.POST("/user/booking", middleware.UserAuthMiddleware()(u.CreateBooking))
	e.GET("/user/booking/:booking_id", middleware.UserAuthMiddleware()(u.GetBookingByID))
	e.GET("/user/booking", middleware.UserAuthMiddleware()(u.GetAllBookingByUser))
}

func (u *BookingHandler) CreateBooking(c echo.Context) error {
	var bookingReq *model.BookingRequest
	if err := c.Bind(&bookingReq); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	userID, err := UserToken(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	booking := &model.Bookings{
		ProductID: bookingReq.ProductID,
		StartDate: bookingReq.StartDate,
		EndDate:   bookingReq.EndDate,
	}

	booking.UserID = userID        // set user id from token
	booking.Status = model.Pending // set status to pending

	isOwner, err := u.bookingUsecase.IsUserProductOwner(userID, booking.ProductID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if isOwner {
		return echo.NewHTTPError(http.StatusForbidden, "lessors cannot book their own products")
	}

	booking, err = u.bookingUsecase.CreateBooking(booking)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	lessor, err := u.lessorUsecase.GetLessorByProductID(booking.ProductID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	product, err := u.productUsecase.GetProductByID(booking.ProductID, lessor.LessorID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	bookingData := model.BookingData{
		BookingID:   booking.BookingID,
		ProductName: product.Name,
		StartDate:   booking.StartDate,
		EndDate:     booking.EndDate,
		Status:      string(booking.Status),
	}

	if err := u.productUsecase.DecrementStockAvailability(booking.ProductID); err != nil {
		fmt.Printf("failed to decrement stock: %v\n", err)
	}

	user, err := u.userUsecase.GetUserByID(userID)
	if err == nil {
		go func() {
			err := utils.SendBookingNotification(user.Email, user.Name, string(model.Pending), booking.BookingID, product.RentalCostPerMonth)
			if err != nil {
				fmt.Printf("failed to send topup notification: %v\n", err)
			}
		}()
	}

	response := model.BookingResponse{
		Message: "success create booking",
		Data:    []model.BookingData{bookingData},
	}

	return c.JSON(http.StatusOK, response)
}

func (u *BookingHandler) GetBookingByID(c echo.Context) error {
	bookingID := c.Param("booking_id")
	id := utils.StringToInt(bookingID)

	userID, err := UserToken(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	booking, err := u.bookingUsecase.GetBookingByID(id, userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	bookingData := model.BookingData{
		BookingID:   booking.BookingID,
		ProductName: booking.Products.Name,
		StartDate:   booking.StartDate,
		EndDate:     booking.EndDate,
		Status:      string(booking.Status),
	}

	response := model.BookingResponse{
		Message: "success get booking",
		Data:    []model.BookingData{bookingData},
	}

	return c.JSON(http.StatusOK, response)
}

func (u *BookingHandler) GetAllBookingByUser(c echo.Context) error {
	userID, err := UserToken(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	bookings, err := u.bookingUsecase.GetAllBookingByUser(userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	var bookingData []model.BookingData
	for _, booking := range bookings {
		bookingData = append(bookingData, model.BookingData{
			BookingID:   booking.BookingID,
			ProductName: booking.Products.Name,
			StartDate:   booking.StartDate,
			EndDate:     booking.EndDate,
			Status:      string(booking.Status),
		})
	}

	response := model.BookingResponse{
		Message: "success get all user booking",
		Data:    bookingData,
	}

	return c.JSON(http.StatusOK, response)
}
