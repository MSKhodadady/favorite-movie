package apis

import (
	"database/sql"

	"github.com/labstack/echo/v4"
	myPac "sadeq.go/favorite-movie/lib"
)

func UsernameList(e *echo.Echo, db *sql.DB) {

	e.GET("/api/user-list", func(c echo.Context) error {
		rows, err := db.Query(myPac.QFirst20UserName)

		if err != nil {
			panic(err)
		}

		type UserMovieCount struct {
			Name       string `json:"name"`
			MovieCount int    `json:"movieCount"`
		}

		var res []UserMovieCount
		for rows.Next() {
			umc := new(UserMovieCount)
			rows.Scan(&umc.Name, &umc.MovieCount)

			res = append(res, *umc)
		}
		rows.Close()

		return c.JSON(200, res)

	})

}
