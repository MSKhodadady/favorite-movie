package apis

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	lib "sadeq.go/favorite-movie/lib"
)

func SuggestFilmApi(e *echo.Echo, appConf lib.AppConfig) {
	e.POST("/api/suggest", func(c echo.Context) error {
		sgst := new(lib.SearchText)

		if err := c.Bind(sgst); err != nil {
			fmt.Println(err)
			return c.String(http.StatusBadRequest, "bad request")
		} else if err := c.Validate(sgst); err != nil {
			return err
		}

		ms, errs := lib.SearchForMovieName(sgst.Text)

		if len(ms) == 0 && len(errs) > 0 {
			fmt.Println("-- error: ", errs)
			return c.String(http.StatusInternalServerError, "")
		} else if len(errs) > 0 {
			fmt.Println("-- error: ", errs)
		}

		msWithHas := make([]lib.Suggestion, len(ms))

		for i, m := range ms {
			msWithHas[i] = lib.Suggestion{
				Movie: *m,
				Hash: base64.
					StdEncoding.
					EncodeToString(
						sha1.
							New().
							Sum([]byte(m.Name + m.Year + appConf.MovieKey)),
					),
			}
		}

		return c.JSON(200, msWithHas)

	})
}
