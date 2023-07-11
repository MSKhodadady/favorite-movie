package myPac

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func CreateToken(username string, appConf AppConfig) string {
	urc := UserRC{
		User: User{Name: username, Password: ""},
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt: jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(
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
			)),
		},
	}
	tokenBuilder := jwt.NewWithClaims(jwt.SigningMethodHS256, urc)
	s, err := tokenBuilder.SignedString([]byte(appConf.JwtKey))
	if err != nil {
		panic(err)
	}

	return s
}
