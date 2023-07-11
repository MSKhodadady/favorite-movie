package myPac

import (
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type (
	Movie struct {
		Name string `json:"name" validate:"required" selector:"div>div>a"`
		Year string `json:"year" validate:"required" selector:"div>div>ul>li>span"`
	}
	User struct {
		Name     string `json:"username" validator:"required|ascii"`
		Password string `json:"password" validator:"required"`
	}
	UserRC struct {
		User
		jwt.RegisteredClaims
	}
	CustomValidator struct {
		/* Validator *validator.Validate */
	}
	AppConfig struct {
		JwtKey    string `json:"jwt-key"`
		MovieKey  string `json:"movie-hash-secret"`
		DbConnStr string `json:"db-connection-string"`
		TokenExp  int    `json:"token-exp-count"`
		// possible values: milisec sec min hour
		TokenExpUnit  string `json:"token-exp-unit"`
		ServerAddress string `json:"server-address"`
		TLSEnabled    bool   `json:"tls-enabled"`
		TLSCertFile   string `json:"tls-cert-file"`
		TLSKeyFile    string `json:"tls-key-file"`
	}
	SearchText struct {
		Text string `json:"text" validator:"required"`
	}
	Suggestion struct {
		Movie
		Hash string `json:"hash" validator:"required"`
	}
)

func (m Movie) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.Name, validation.Required),
		validation.Field(&m.Year, validation.Required),
	)
}

func (u User) Validate() error {
	return validation.ValidateStruct(&u,
		validation.Field(&u.Name, validation.Required, is.ASCII),
		validation.Field(&u.Password, validation.Required))
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

func (cv CustomValidator) Validate(i any) error {
	v, ok := i.(validation.Validatable)
	if ok {
		if err := v.Validate(); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
	} else {
		panic("value isn't validator")
	}
	return nil
}
