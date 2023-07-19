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

		return c.String(200, "OK")

		/* switch err {
		case sql.ErrNoRows:
			db.Exec(lib.QSignUp, u.Name, lib.HashPass(u.Password))

			return c.JSON(http.StatusOK, map[string]string{
				"token": lib.CreateToken(u.Name, appConf),
			})
		case nil:
			return c.String(http.StatusConflict, "another username")
		default:
			panic(err)
		} */

	})
}
