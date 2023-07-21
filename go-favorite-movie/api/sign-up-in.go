package apis

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"gopkg.in/gomail.v2"

	lib "sadeq.go/favorite-movie/lib"
)

type SignUpUserRC struct {
	lib.SignUpUser
	*jwt.RegisteredClaims
}

func SignInUpApi(e *echo.Echo, db *sql.DB, appConf lib.AppConfig, emailDialer *gomail.Dialer) {
	e.POST("/api/sign-in", func(c echo.Context) error {
		u := new(lib.User)

		if err := c.Bind(u); err != nil {
			return c.String(http.StatusBadRequest, "bad request")
		} else if err := c.Validate(u); err != nil {
			return err
		}

		var dbUser string
		var dbPassHash string
		err := db.
			QueryRow(lib.QSignIn, u.Name).
			Scan(&dbUser, &dbPassHash)

		switch err {
		case sql.ErrNoRows:
			return c.String(http.StatusNotFound, "username or password")
		case nil:
			passHash := lib.HashPass(u.Password)

			if dbPassHash == passHash {

				return c.JSON(http.StatusOK, map[string]string{
					"token": lib.CreateToken(dbUser, appConf),
				})
			} else {
				return c.String(http.StatusNotFound, "username or password")
			}
		default:
			panic(err)
		}
	})

	e.POST("/api/sign-up", func(c echo.Context) error {
		u := new(lib.SignUpUser)

		if err := c.Bind(u); err != nil {
			return c.String(http.StatusBadRequest, "bad request")
		} else if err := c.Validate(u); err != nil {
			return err
		}

		var dbUser, dbEmail string
		errUsername := db.
			QueryRow(
				"SELECT username FROM tbUser WHERE username = $1", u.Username).
			Scan(&dbUser)

		errEmail := db.
			QueryRow(
				"SELECT email FROM tbUser WHERE email = $1", u.Email).
			Scan(&dbEmail)
		// means there some record exist with that username or email
		if errUsername == nil || errEmail == nil {
			return c.JSON(http.StatusConflict, map[string]bool{
				"username": errUsername == nil,
				"email":    errEmail == nil,
			})
		}
		// panic if there exists another error other than ErrNoRows
		if errUsername != sql.ErrNoRows {
			panic(errUsername)
		}
		if errEmail != sql.ErrNoRows {
			panic(errEmail)
		}

		// create token

		t := jwt.NewWithClaims(jwt.SigningMethodHS256, SignUpUserRC{
			SignUpUser: *u,
			RegisteredClaims: &jwt.RegisteredClaims{
				IssuedAt:  jwt.NewNumericDate(time.Now()),
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
			},
		})

		s, err := t.SignedString([]byte(appConf.SignUpJwtToken))

		if err != nil {
			panic(err)
		}
		// create email url
		// TODO must refer to fronted
		url := lib.FullServerAddress(appConf) + "/api/verify?token=" + s
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
		m.SetHeader("From", "one.worship@outlook.com")
		m.SetHeader("To", "worldofshie@gmail.com")
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

		token, err := jwt.ParseWithClaims(tokenString, &SignUpUserRC{}, func(t *jwt.Token) (interface{}, error) {
			return []byte(appConf.SignUpJwtToken), nil
		})

		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}

		claims := token.Claims.(*SignUpUserRC)

		fmt.Println(claims)
		// TODO add user to db

		return c.String(200, "token verified!")

	})

}
