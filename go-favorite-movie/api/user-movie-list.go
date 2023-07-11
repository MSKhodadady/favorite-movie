package apis

import (
	"database/sql"

	"github.com/labstack/echo/v4"
	lib "sadeq.go/favorite-movie/lib"
)

func UserMovieList(e *echo.Echo, db *sql.DB, appConf lib.AppConfig) {
	e.GET("/api/u/:username", func(c echo.Context) error {
		username := c.Param("username")

		var dbUser string
		err := db.
			QueryRow(`SELECT username FROM tbUser
			WHERE username = $1`, username).Scan(&dbUser)

		switch err {
		case sql.ErrNoRows:
			return c.String(404, "no such user")
		case nil:
			rows, err := db.
				Query(`SELECT movie_name, "year"
				FROM tbMovie WHERE username = $1`, username)
			if err != nil {
				panic(err)
			}
			var res []lib.Movie
			for rows.Next() {
				m := new(lib.Movie)
				rows.Scan(&m.Name, &m.Year)

				res = append(res, *m)
			}
			rows.Close()

			return c.JSON(200, map[string]any{
				"username": username,
				"movies":   res,
			})
		default:
			panic(err)
		}
	})
}
