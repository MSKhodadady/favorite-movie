package apis

import (
	"crypto/sha1"
	"database/sql"
	"encoding/base64"
	"net/http"

	"github.com/labstack/echo/v4"
	lib "sadeq.go/favorite-movie/lib"
)

func AddDelMovieApi(e *echo.Echo, db *sql.DB, appConf lib.AppConfig) {
	e.POST("/api/movie", func(c echo.Context) error {
		username := lib.ExtractUsernameToken(c)

		sgst := new(lib.Suggestion)
		if err := c.Bind(sgst); err != nil {
			return c.String(http.StatusNotAcceptable, "not acceptable")
		} else if err := c.Validate(sgst); err != nil {
			return err
		}

		h := base64.
			StdEncoding.
			EncodeToString(
				sha1.
					New().
					Sum([]byte(sgst.Name + sgst.Year + appConf.MovieKey)),
			)

		if h != sgst.Hash {
			return c.String(http.StatusBadRequest, "not a suggested movie")
		}

		_, err := db.Exec(
			lib.QAddMovie,
			username,
			sgst.Name, sgst.Year)

		if err != nil {
			return c.String(http.StatusConflict, "duplicate")
		}

		return c.String(200, "added")
	})

	e.DELETE("/api/movie", func(c echo.Context) error {
		username := lib.ExtractUsernameToken(c)

		m := new(lib.Movie)
		if err := c.Bind(m); err != nil {
			return c.String(http.StatusNotAcceptable, "not acceptable")
		} else if err := c.Validate(m); err != nil {
			return err
		}

		res, err := db.Exec(lib.QDeleteMovie,
			username, m.Name, m.Year)

		if err != nil {
			panic(err)
		}

		if ra, err := res.RowsAffected(); err != nil {
			panic(err)
		} else if ra == 0 {
			return c.String(http.StatusBadGateway, "no such movie of user")
		}
		return c.String(200, "deleted")
	})
}
