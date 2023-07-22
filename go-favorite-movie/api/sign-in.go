package apis

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	lib "sadeq.go/favorite-movie/lib"
)

func SignInApi(e *echo.Echo, db *sql.DB, appConf lib.AppConfig) {
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

				token, err := lib.CreateToken(dbUser, appConf)

				if err != nil {
					fmt.Println("-- error: ", err)
					return c.String(http.StatusInternalServerError, "")
				}

				return c.JSON(http.StatusOK, map[string]string{
					"token": token,
				})
			} else {
				return c.String(http.StatusNotFound, "username or password")
			}
		default:
			fmt.Println("-- error: ", err)
			return c.String(http.StatusInternalServerError, "")
		}
	})

}
