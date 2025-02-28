package handler

import (
	"fmt"
	"net/http"
	"os"
	"rent-video-game/middleware"
	"rent-video-game/model"
	"rent-video-game/usecase"
	"rent-video-game/utils"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/paymentintent"
	"golang.org/x/crypto/bcrypt"
)

var (
	Authorization = "Authorization"
	Bearer        = "Bearer "
)

type IUserHandler interface {
	UserRoutes(e *echo.Echo)
	RegisterUser(c echo.Context) error
	LoginUser(c echo.Context) error
	TopupUser(c echo.Context) error
}

type UserHandler struct {
	userUsecase         *usecase.UserUsecase
	topupHistoryUsecase *usecase.TopupHistoryUsecase
}

type UserHandlerInterface struct {
	userUsecase         usecase.IUserUsecase
	topupHistoryUsecase usecase.ITopupHistoryUsecase
}

func NewUserHandler(userUsecase *usecase.UserUsecase, topupHistoryUsecase *usecase.TopupHistoryUsecase) *UserHandler {
	return &UserHandler{
		userUsecase:         userUsecase,
		topupHistoryUsecase: topupHistoryUsecase,
	}
}

func NewUserHandlerWithInterface(userUsecase usecase.IUserUsecase, topupHistoryUsecase usecase.ITopupHistoryUsecase) *UserHandlerInterface {
	return &UserHandlerInterface{
		userUsecase:         userUsecase,
		topupHistoryUsecase: topupHistoryUsecase,
	}
}

func (u *UserHandler) UserRoutes(e *echo.Echo) {
	e.POST("/user/register", u.RegisterUser)
	e.POST("/user/login", u.LoginUser)
	e.POST("/user/topup", middleware.UserAuthMiddleware()(u.TopupUser))
}

func (u *UserHandler) RegisterUser(c echo.Context) error {
	var userRegister *model.RegisterRequest
	if err := c.Bind(&userRegister); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userRegister.Password), bcrypt.DefaultCost)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	userRegister.Password = string(hashedPassword) // save hashed password

	user := &model.Users{
		Name:     userRegister.Name,
		Email:    userRegister.Email,
		Password: userRegister.Password,
		Address:  userRegister.Address,
	}

	_, err = u.userUsecase.GetUserByEmail(user.Email)
	if err == nil {
		return echo.NewHTTPError(http.StatusBadRequest, "user already exists")
	}

	user, err = u.userUsecase.RegisterUser(user)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	response := model.RegisterResponse{
		Message: "success register user",
	}
	response.Data.UserID = user.UserID
	response.Data.Email = user.Email
	response.Data.Name = user.Name

	return c.JSON(http.StatusOK, response)
}

func (ui *UserHandlerInterface) RegisterUserInterface(c echo.Context) error {
	var userRegister *model.RegisterRequest
	if err := c.Bind(&userRegister); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userRegister.Password), bcrypt.DefaultCost)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	userRegister.Password = string(hashedPassword) // save hashed password

	user := &model.Users{
		Name:     userRegister.Name,
		Email:    userRegister.Email,
		Password: userRegister.Password,
		Address:  userRegister.Address,
	}

	_, err = ui.userUsecase.GetUserByEmail(user.Email)
	if err == nil {
		return echo.NewHTTPError(http.StatusBadRequest, "user already exists")
	}

	user, err = ui.userUsecase.RegisterUser(user)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	response := model.RegisterResponse{
		Message: "success register user",
	}
	response.Data.UserID = user.UserID
	response.Data.Email = user.Email
	response.Data.Name = user.Name

	return c.JSON(http.StatusOK, response)
}

func (u *UserHandler) LoginUser(c echo.Context) error {
	var loginReq *model.LoginRequest
	if err := c.Bind(&loginReq); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	user, err := u.userUsecase.GetUserByEmail(loginReq.Email)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid email")
	}

	if err := user.CompareHashAndPassword(loginReq.Password); err != nil {
		if loginReq.Password != user.Password {
			return echo.NewHTTPError(http.StatusUnauthorized, "invalid password")
		}
	}

	token, err := utils.GenerateUserToken(user.UserID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to generate token")
	}

	response := model.LoginResponse{
		Message: "success login user",
	}
	response.Data.Token = token

	return c.JSON(http.StatusOK, response)
}

func (u *UserHandler) TopupUser(c echo.Context) error {
	var topupReq model.TopupRequest

	userID, err := UserToken(c)
	if err != nil {
		return err
	}

	if err := c.Bind(&topupReq); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request body: " + err.Error(),
		})
	}

	if topupReq.Amount <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "amount must be greater than 0",
		})
	}

	// stripe payment
	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")

	// create payment intent
	params := &stripe.PaymentIntentParams{
		Amount:             stripe.Int64(int64(topupReq.Amount * 100)), // convert to cents
		Currency:           stripe.String("usd"),
		PaymentMethod:      stripe.String("pm_card_visa"),
		PaymentMethodTypes: []*string{stripe.String("card")},
		ConfirmationMethod: stripe.String("automatic"),
		Confirm:            stripe.Bool(true),
	}

	pi, err := paymentintent.New(params)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "payment processing failed: " + err.Error(),
		})
	}

	paymentID := pi.ID // save payment id

	if pi.Status == "succeeded" {
		user := &model.Users{
			Amount: topupReq.Amount,
		}

		updatedUser, err := u.userUsecase.TopupUser(userID, user)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": "failed to update user balance: " + err.Error(),
			})
		}

		user, err = u.userUsecase.GetUserByID(userID)
		if err == nil {
			go func() {
				err := utils.SendTopupNotification(user.Email, user.Name, topupReq.Amount, updatedUser.Amount, pi.ID)
				if err != nil {
					fmt.Printf("failed to send topup notification: %v\n", err)
				}
			}()
		}

		topupHistory := &model.TopupHistory{
			UserID:    userID,
			PaymentID: paymentID,
			Amount:    topupReq.Amount,
		}

		_, err = u.topupHistoryUsecase.CreateTopupHistory(topupHistory)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": "failed to create topup history: " + err.Error(),
			})
		}

		response := model.TopupResponse{
			Message: "topup successful",
		}

		response.Data.UserID = updatedUser.UserID
		response.Data.TopupAmount = topupReq.Amount
		response.Data.NewBalance = updatedUser.Amount
		response.Data.PaymentID = paymentID

		return c.JSON(http.StatusOK, response)
	} else {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "payment failed: " + string(pi.Status),
		})
	}
}

func UserToken(c echo.Context) (uuid.UUID, error) {
	tokenString := c.Request().Header.Get(Authorization)
	if tokenString == "" {
		return uuid.Nil, echo.NewHTTPError(http.StatusUnauthorized, "token required")
	}

	tokenString = strings.TrimPrefix(tokenString, Bearer)

	claims, err := utils.VerifyUserToken(tokenString)
	if err != nil {
		return uuid.Nil, echo.NewHTTPError(http.StatusUnauthorized, "invalid token")
	}

	userIDString := claims["user_id"].(string)
	userID, err := uuid.Parse(userIDString)
	if err != nil {
		return uuid.Nil, echo.NewHTTPError(http.StatusUnauthorized, "invalid user id")
	}

	return userID, nil
}
