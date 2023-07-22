package myPac

import (
	"crypto/sha256"
	"encoding/base64"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func CreateToken(username string, appConf AppConfig) (string, error) {
	urc := UserRC{
		User: User{Name: username, Password: ""},
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(CalculateTokenExp(appConf)),
		},
	}
	tokenBuilder := jwt.NewWithClaims(jwt.SigningMethodHS256, urc)
	s, err := tokenBuilder.SignedString([]byte(appConf.JwtKey))
	if err != nil {
		return "", err
	}

	return s, nil
}

func CalculateTokenExp(appConf AppConfig) time.Time {
	return time.Now().Add(
		func() time.Duration {
			var unit time.Duration
			switch appConf.TokenExpUnit {
			case "milisec":
				unit = time.Millisecond
			case "sec":
				unit = time.Second
			case "min":
				unit = time.Minute
			case "hour":
				unit = time.Hour
			default:
				unit = time.Minute
			}
			return unit * time.Duration(appConf.TokenExp)
		}(),
	)
}

func ExtractUsernameToken(c echo.Context) string {
	username := c.
		Get("user").(*jwt.Token).
		Claims.(jwt.MapClaims)["username"].(string)

	return username
}

func HashPass(pass string) string {
	return base64.StdEncoding.EncodeToString(
		sha256.New().Sum([]byte(pass)),
	)
}

func FullServerAddress(c AppConfig) string {
	if c.TLS.Enabled {
		return "https://" + c.ServerAddress
	} else {
		return "http://" + c.ServerAddress
	}
}
