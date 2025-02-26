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

type UserHandler struct {
	userUsecase *usecase.UserUsecase
}

func NewUserHandler(userUsecase *usecase.UserUsecase) *UserHandler {
	return &UserHandler{userUsecase: userUsecase}
}

func (u *UserHandler) UserRoutes(e *echo.Echo) {
	e.POST("/user/register", u.RegisterUser)
	e.POST("/user/login", u.LoginUser)
	e.POST("/user/topup", middleware.UserAuthMiddleware()(u.TopupUser))
}

func (u *UserHandler) RegisterUser(c echo.Context) error {
	var user *model.Users
	if err := c.Bind(&user); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	user.Password = string(hashedPassword)

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
	userID, err := UserToken(c)
	if err != nil {
		return err
	}

	var topupReq model.TopupRequest

	if err := c.Bind(&topupReq); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body: " + err.Error(),
		})
	}

	if topupReq.Amount <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Amount must be greater than 0",
		})
	}

	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")

	params := &stripe.PaymentIntentParams{
		Amount:             stripe.Int64(int64(topupReq.Amount * 100)),
		Currency:           stripe.String("usd"),
		PaymentMethod:      stripe.String("pm_card_visa"),
		PaymentMethodTypes: []*string{stripe.String("card")},
		ConfirmationMethod: stripe.String("automatic"),
		Confirm:            stripe.Bool(true),
	}

	pi, err := paymentintent.New(params)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Payment processing failed: " + err.Error(),
		})
	}

	paymentID := pi.ID

	if pi.Status == "succeeded" {
		updatedUser, err := u.userUsecase.UpdateUserAmount(userID, topupReq.Amount)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": "Failed to update user balance: " + err.Error(),
			})
		}

		user, err := u.userUsecase.GetUserByID(userID)
		if err == nil {
			go func() {
				err := utils.SendTopupNotification(user.Email, user.Name, topupReq.Amount, updatedUser.Amount, pi.ID)
				if err != nil {
					fmt.Printf("Failed to send topup notification: %v\n", err)
				}
			}()
		}

		response := model.TopupResponse{
			Message: "Topup successful",
		}

		response.Data.UserID = updatedUser.UserID
		response.Data.TopupAmount = topupReq.Amount
		response.Data.NewBalance = updatedUser.Amount
		response.Data.PaymentID = paymentID

		return c.JSON(http.StatusOK, response)
	} else {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Payment failed: " + string(pi.Status),
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
