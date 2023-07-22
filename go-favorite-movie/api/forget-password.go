package apis

/* import (
	"database/sql"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"gopkg.in/gomail.v2"
)

type EmailRC struct {
	Email string
	jwt.RegisteredClaims
}

func ForgetPassApi(e *echo.Echo, db *sql.DB, emailDialer *gomail.Dialer) {
	e.POST("/api/forget-pass", func(c echo.Context) error {
		e := new(struct {
			email string
		})

		if err := c.Bind(e); err != nil {
			return c.String(echo.ErrBadRequest.Code, "")
		} else if err := validation.Validate(e.email,
			validation.Required,
			is.Email); err != nil {
			return c.String(echo.ErrBadRequest.Code, err.Error())
		}

		var dbEmail string

		errEmail := db.
			QueryRow(
				"SELECT email FROM tbUser WHERE email = $1", e.email).
			Scan(&dbEmail)

		switch errEmail {
		case nil:



		case sql.ErrNoRows:
			return c.String(404, "not any user with this email")
		default:

		}

	})
}
*/
