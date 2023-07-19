package myPac

import (
	"github.com/golang-jwt/jwt/v5"
)

type (
	Movie struct {
		Name string `json:"name" selector:"div>div>a"`
		Year string `json:"year" selector:"div>div>ul>li>span"`
	}
	User struct {
		Name     string `json:"username"`
		Password string `json:"password"`
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
		Text string `json:"text"`
	}
	Suggestion struct {
		Movie
		Hash string `json:"hash"`
	}
	SignUpUser struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Email    string `json:"email"`
	}
)
