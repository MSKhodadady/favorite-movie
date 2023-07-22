package myPac

import (
	"fmt"
	"net/http"
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/labstack/echo/v4"
)

func (m Movie) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.Name, validation.Required),
		validation.Field(&m.Year, validation.Required),
	)
}

func (u User) Validate() error {
	return validation.ValidateStruct(&u,
		validation.Field(&u.Name,
			validation.Required,
			validation.Match(regexp.MustCompile("^[a-zA-Z_][a-zA-Z_0-9]{4,}$")),
		),
		validation.Field(&u.Password,
			validation.Required,
			// TODO not sure!
			// validation.Length(8, 0),
			// validation.Match(regexp.MustCompile("^[a-zA-Z0-9@$!%*?&_]*$")),
		),
	)
}

func (s SearchText) Validate() error {
	return validation.ValidateStruct(&s,
		validation.Field(&s.Text, validation.Required))
}

func (sgst Suggestion) Validate() error {
	return validation.ValidateStruct(&sgst,
		validation.Field(&sgst.Movie),
		validation.Field(&sgst.Hash, validation.Required))
}
func (u SignUpUser) Validate() error {
	return validation.ValidateStruct(&u,
		validation.Field(&u.Username,
			validation.Required,
			validation.Match(regexp.MustCompile("^[a-zA-Z_][a-zA-Z_0-9]{4,}$")),
		),
		validation.Field(&u.Password,
			validation.Required,
			validation.Length(8, 0),
			validation.Match(regexp.MustCompile("^[a-zA-Z0-9@$!%*?&_]*$")),
			validation.Match(regexp.MustCompile("^.*[0-9].*$")),
			validation.Match(regexp.MustCompile("^.*[a-z].*$")),
			validation.Match(regexp.MustCompile("^.*[A-Z].*$")),
			validation.Match(regexp.MustCompile("^.*[@$!%*?&_].*$")),
		),
		validation.Field(&u.Email,
			validation.Required,
			is.Email,
		),
	)
}

func (cv CustomValidator) Validate(i any) error {
	v, ok := i.(validation.Validatable)
	if ok {
		if err := v.Validate(); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
	} else {
		fmt.Println("value isn't validator")
		return echo.NewHTTPError(http.StatusInternalServerError, "")

	}
	return nil
}
