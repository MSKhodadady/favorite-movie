package main

import (
	"database/sql"
	"encoding/json"
	"os"
	"regexp"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
	apis "sadeq.go/favorite-movie/api"
	lib "sadeq.go/favorite-movie/lib"
)

// : MAIN ----------------------------------------------------------------------
func main() {

	//: app config ------------------------------------------------------------
	var appConf lib.AppConfig
	confByte, err := os.ReadFile("env.json")
	if err != nil {
		panic(err)
	}
	json.Unmarshal(confByte, &appConf)

	//: config echo -----------------------------------------------------------
	e := echo.New()
	e.Validator = &lib.CustomValidator{}

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{"Content-Type", "Authorization"},
	}))

	e.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(appConf.JwtKey),
		Skipper: func(c echo.Context) bool {
			isSignIn, _ := regexp.MatchString(
				"/api/sign-(in|up)",
				c.Request().URL.Path,
			)
			isGET := c.Request().Method == "GET"

			return isSignIn || isGET
		},
	}))

	//: config db -------------------------------------------------------------
	db, err := sql.Open("postgres", appConf.DbConnStr)

	if err != nil {
		panic(err)
	}

	defer db.Close()

	//: user movie list
	apis.UserMovieList(e, db, appConf)
	//: add movie
	apis.AddDelMovieApi(e, db, appConf)
	//: sign-up sign-in
	apis.SignInUpApi(e, db, appConf)
	//: suggest film
	apis.SuggestFilmApi(e, appConf)
	//: user list
	apis.UsernameList(e, db)
	//: frontend
	e.Static("/", "frontend")
	//: start server
	if appConf.TLSEnabled {
		e.Logger.Info(e.StartTLS(appConf.ServerAddress, appConf.TLSCertFile, appConf.TLSKeyFile))
	} else {
		e.Logger.Info(e.Start(appConf.ServerAddress))
	}
}
