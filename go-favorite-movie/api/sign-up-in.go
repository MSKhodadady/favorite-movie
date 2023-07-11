package apis

import (
	"database/sql"
	"net/http"

	"github.com/labstack/echo/v4"

	lib "sadeq.go/favorite-movie/lib"
)

func SignInUpApi(e *echo.Echo, db *sql.DB, appConf lib.AppConfig) {
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

	/* e.POST("/api/sign-up", func(c echo.Context) error {
		u := new(lib.User)

		if err := c.Bind(u); err != nil {
			return c.String(http.StatusBadRequest, "bad request")
		} else if err := c.Validate(u); err != nil {
			return err
		}

		var dbUser string
		err := db.
			QueryRow(
				"SELECT username FROM tbUser where username = $1", u.Name).
			Scan(&dbUser)

		switch err {
		case sql.ErrNoRows:
			db.Exec(lib.QSignUp, u.Name, hashPass(u.Password))

			return c.JSON(http.StatusOK, map[string]string{
				"token": createToken(u.Name, appConf.JwtKey),
			})
		case nil:
			return c.String(http.StatusConflict, "another username")
		default:
			panic(err)
		}

	}) */
}
