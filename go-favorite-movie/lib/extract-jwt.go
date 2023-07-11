package myPac

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func ExtractUsernameToken(c echo.Context) string {
	username := c.
		Get("user").(*jwt.Token).
		Claims.(jwt.MapClaims)["username"].(string)

	return username
}
