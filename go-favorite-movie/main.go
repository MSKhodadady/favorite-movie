package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"time"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
	apis "sadeq.go/favorite-movie/api"
	lib "sadeq.go/favorite-movie/lib"
)

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

	//: migrations ------------------------------------------------------------
	var migrationVersion int
	var lastMigrations time.Time

	fmt.Println("-- checking db migrations")

	errDb := db.
		QueryRow("SELECT migration_version, last_migration FROM tb_config;").
		Scan(&migrationVersion, &lastMigrations)

	/* sort.Slice(lib.MigrationQueries, func(i, j int) bool {
		return lib.MigrationQueries[i].Version < lib.MigrationQueries[j].Version
	}) */

	if errDb != nil {
		if errDb.Error() == "pq: relation \"tb_config\" does not exist" {
			fmt.Println("-- there isn't any table, running migration 1:")
			mig := lib.MigrationQueries[0]
			fmt.Println("-- mig version: ", mig.Version, " - ", mig.Message)
			for _, v := range mig.Queries {
				fmt.Println(v)
				_, errQ := db.Exec(v)
				if errQ != nil {
					fmt.Println(errQ)
				}
			}
			migrationVersion = 1
		} else {
			panic(errDb)
		}
	}

	lastMigrationVersion := lib.MigrationQueries[len(lib.MigrationQueries)-1].Version
	fmt.Println("-- db migration version: ", migrationVersion)
	fmt.Println("-- last migration: ", lastMigrationVersion)

	if migrationVersion == lastMigrationVersion {
		fmt.Println("-- db is up to latest migration!")
	} else {
		fmt.Println("-- upgrading db to latest migration")

		for _, mig := range lib.MigrationQueries {
			if mig.Version > migrationVersion {
				fmt.Println("-- mig version: ", mig.Version, " -- ", mig.Message)
				for _, v := range mig.Queries {
					fmt.Println(v)
					_, errQ := db.Exec(v)
					if errQ != nil {
						fmt.Println(errQ)
					}
				}
			}
		}
	}

	//: APIS -------------------------------------------------------------
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
	e.File("/u/:username", "frontend/u/[username].html")
	//: start server
	if appConf.TLSEnabled {
		e.Logger.Info(e.StartTLS(appConf.ServerAddress, appConf.TLSCertFile, appConf.TLSKeyFile))
	} else {
		e.Logger.Info(e.Start(appConf.ServerAddress))
	}
}
