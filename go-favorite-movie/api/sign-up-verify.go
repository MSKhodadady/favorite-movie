package apis

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"gopkg.in/gomail.v2"
	l "sadeq.go/favorite-movie/lib"
)

type SignUpUserRC struct {
	l.SignUpUser
	*jwt.RegisteredClaims
}

func SignUpVerifyApi(e *echo.Echo, db *sql.DB, appConf l.AppConfig, emailDialer *gomail.Dialer) {
	e.POST("/api/sign-up", func(c echo.Context) error {
		u := new(l.SignUpUser)

		if err := c.Bind(u); err != nil {
			return c.String(http.StatusBadRequest, "bad request")
		} else if err := c.Validate(u); err != nil {
			return err
		}

		userExist, errU := l.CheckDBExist("username", "tbUser", u.Username, db)
		emailExist, errE := l.CheckDBExist("email", "tbUser", u.Email, db)

		if errU != nil || errE != nil {
			fmt.Println("-- err: username: ", errU)
			fmt.Println("-- err: email: ", errE)
			return c.String(http.StatusInternalServerError, "")
		}

		if userExist || emailExist {
			return c.JSON(http.StatusConflict, map[string]bool{
				"username": userExist,
				"email":    emailExist,
			})
		}
		// hashing password
		u.Password = l.HashPass(u.Password)
		// create token
		t := jwt.NewWithClaims(jwt.SigningMethodHS256, SignUpUserRC{
			SignUpUser: *u,
			RegisteredClaims: &jwt.RegisteredClaims{
				IssuedAt:  jwt.NewNumericDate(time.Now()),
				ExpiresAt: jwt.NewNumericDate(l.CalculateTokenExp(appConf)),
			},
		})

		s, err := t.SignedString([]byte(appConf.SignUpJwtToken))

		if err != nil {
			fmt.Println("-- error: ", err)
			return c.String(http.StatusInternalServerError, "")
		}
		// create email url
		url := appConf.FrontendAddress + "/verify?token=" + s
		// create email body
		body := fmt.Sprintf(`
		<html>
			<body>
				<h1>Verification</h1>
				<a href="%s">Click on Me to verify</a>
			</body>
		</html>
		`, url)

		// send email
		m := gomail.NewMessage()
		m.SetHeader("From", appConf.Smtp.Username)
		m.SetHeader("To", u.Email)
		m.SetHeader("Subject", "FavMov Verification Email")
		m.SetBody("text/html", body)

		if errMail := emailDialer.DialAndSend(m); errMail != nil {
			fmt.Println("-- email err: ", errMail)
			return c.String(http.StatusInternalServerError, "")
		} else {
			return c.String(200, "verification email sent")
		}

	})

	e.GET("/api/verify", func(c echo.Context) error {
		tokenString := c.QueryParam("token")

		if tokenString == "" {
			return c.String(http.StatusBadRequest, "not any token!")
		}

		token, err := jwt.ParseWithClaims(
			tokenString,
			&SignUpUserRC{},
			func(t *jwt.Token) (interface{}, error) {
				return []byte(appConf.SignUpJwtToken), nil
			},
			// we hashed the password and it will fail in validation,
			// so we disable validation
			jwt.WithoutClaimsValidation())

		if err != nil {
			return c.String(http.StatusNotAcceptable, err.Error())
		}

		claims := token.Claims.(*SignUpUserRC)

		if time.Now().After(claims.ExpiresAt.Time) {
			return c.String(http.StatusUnauthorized, "token expired")
		}

		// check if token used to create a user
		// une, errUsername := l.CheckUsernameExist(claims.Username, db)
		une, errUsername := l.CheckDBExist("username", "tbUser", claims.Username, db)

		if errUsername != nil {
			fmt.Println("-- error: ", errUsername)
			return c.String(http.StatusInternalServerError, "")
		}

		if une {
			return c.String(http.StatusUnauthorized, "token expired")
		}

		// add user to db
		_, errDb := db.Exec(l.QAddUser, claims.Username, claims.Password, claims.Email)

		if errDb != nil {
			fmt.Println(errDb)
			return c.String(http.StatusInternalServerError, "")
		}

		loginToken, err := l.CreateToken(claims.Username, appConf)

		if err != nil {
			fmt.Println("-- error: ", err)
			return c.String(http.StatusInternalServerError, "")
		}

		return c.JSON(200, map[string]string{
			"token": loginToken,
		})

	})
}
