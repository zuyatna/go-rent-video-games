package middleware

import (
	"net/http"
	"rent-video-game/utils"

	"strings"

	"github.com/labstack/echo/v4"
)

const (
	Authorization = "Authorization"
	Bearer        = "Bearer"
)

func UserAuthMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token := c.Request().Header.Get(Authorization)
			if token == "" {
				return c.JSON(http.StatusBadRequest, map[string]interface{}{
					"message": "token is required!",
				})
			}

			tokenParts := strings.Split(token, " ")
			if len(tokenParts) != 2 || tokenParts[0] != Bearer {
				return c.JSON(http.StatusBadRequest, map[string]interface{}{
					"message": "invalid token format! please using Bearer [token]",
				})
			}

			claims, err := utils.VerifyUserToken(tokenParts[1])
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]interface{}{
					"message": "invalid token!",
				})
			}

			c.Set("user_id", claims["user_id"])
			return next(c)
		}
	}
}
