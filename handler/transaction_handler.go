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

type TransactionHandler struct {
	transactionUsecase *usecase.TransactionUsecase
	bookingUsecase     *usecase.BookingUsecase
	userUsecase        *usecase.UserUsecase
	lessorUsecase      *usecase.LessorUsecase
}

func NewTransactionHandler(
	transactionUsecase *usecase.TransactionUsecase,
	bookingUsecase *usecase.BookingUsecase,
	userUsecase *usecase.UserUsecase,
	lessorUsecase *usecase.LessorUsecase,
) *TransactionHandler {
	return &TransactionHandler{
		transactionUsecase: transactionUsecase,
		bookingUsecase:     bookingUsecase,
		userUsecase:        userUsecase,
		lessorUsecase:      lessorUsecase,
	}
}

func (u *TransactionHandler) TransactionRoutes(e *echo.Echo) {
	e.POST("/user/transaction", middleware.UserAuthMiddleware()(u.CreateTransaction))
	e.GET("/user/transaction/:transaction_id", middleware.UserAuthMiddleware()(u.GetTransactionByID))
	e.GET("/user/transactions", middleware.UserAuthMiddleware()(u.GetAllTransactionByUser))
}

func (u *TransactionHandler) CreateTransaction(c echo.Context) error {
	var transaction *model.Transactions
	if err := c.Bind(&transaction); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	userID, err := UserToken(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	renter, err := u.userUsecase.GetUserByID(userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get user data: "+err.Error())
	}

	if renter.Amount < transaction.Amount {
		return echo.NewHTTPError(http.StatusBadRequest, "insufficient balance")
	}

	booking, err := u.bookingUsecase.GetBookingByID(transaction.BookingID, userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get booking data: "+err.Error())
	}

	product, err := u.bookingUsecase.GetProductByID(booking.ProductID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get product data: "+err.Error())
	}

	lessor, err := u.lessorUsecase.GetLessorByID(product.LessorID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get lessor data: "+err.Error())
	}

	transaction.LessorID = lessor.LessorID

	lessorUser, err := u.userUsecase.GetUserByID(lessor.UserID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get lessor's user data: "+err.Error())
	}

	transaction.UserID = userID

	transaction, err = u.transactionUsecase.CreateTransaction(transaction)
	if err != nil {
		renter.Amount += transaction.Amount
		_, revertErr := u.userUsecase.TransactionUser(userID, renter)
		if revertErr != nil {
			fmt.Printf("failed to revert user balance: %v\n", revertErr)
		}

		lessorUser.Amount -= transaction.Amount
		_, revertErr = u.userUsecase.TransactionUser(lessor.UserID, lessorUser)
		if revertErr != nil {
			fmt.Printf("failed to revert lessor balance: %v\n", revertErr)
		}

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	renter.Amount -= transaction.Amount
	_, err = u.userUsecase.TransactionUser(userID, renter)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to update user balance: "+err.Error())
	}

	lessorUser.Amount += transaction.Amount
	_, err = u.userUsecase.TransactionUser(lessor.UserID, lessorUser)
	if err != nil {
		renter.Amount += transaction.Amount
		_, revertErr := u.userUsecase.TransactionUser(userID, renter)
		if revertErr != nil {
			fmt.Printf("failed to revert user balance: %v\n", revertErr)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to update lessor balance: "+err.Error())
	}

	_, err = u.bookingUsecase.UpdateBooking(transaction.BookingID, model.Approved, booking)
	if err != nil {
		fmt.Printf("failed to update booking status: %v\n", err)
	}

	user, err := u.userUsecase.GetUserByID(userID)
	if err == nil {
		go func() {
			err := utils.SendBookingNotification(user.Email, user.Name, string(model.Approved), booking.BookingID, transaction.Amount)
			if err != nil {
				fmt.Printf("failed to send topup notification: %v\n", err)
			}

			err = utils.SendTransactionNotification(lessorUser.Email, lessorUser.Name, transaction.TransactionID, transaction.Amount, lessor.UserID, lessorUser.Amount)
			if err != nil {
				fmt.Printf("failed to send transaction notification: %v\n", err)
			}

			err = utils.SendTransactionNotification(user.Email, user.Name, transaction.TransactionID, transaction.Amount, userID, user.Amount)
			if err != nil {
				fmt.Printf("failed to send transaction notification: %v\n", err)
			}
		}()
	}

	transactionData := model.TransactionData{
		TransactionID: transaction.TransactionID,
		BookingID:     transaction.BookingID,
		ReceiveID:     transaction.LessorID,
		Amount:        transaction.Amount,
	}

	response := model.TransactionResponse{
		Message: "transaction created successfully",
		Data:    []model.TransactionData{transactionData},
	}

	return c.JSON(http.StatusOK, response)
}

func (u *TransactionHandler) GetTransactionByID(c echo.Context) error {
	transactionID := c.Param("transaction_id")
	id := utils.StringToInt(transactionID)

	transaction, err := u.transactionUsecase.GetTransactionByID(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, transaction)
}

func (u *TransactionHandler) GetAllTransactionByUser(c echo.Context) error {
	userID, err := UserToken(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	transactions, err := u.transactionUsecase.GetAllTransactionByUser(userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, transactions)
}
