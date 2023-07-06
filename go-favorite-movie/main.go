package main

import (
	"crypto/sha1"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"time"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
	lib "sadeq.go/favorite-movie/lib"
)

// : global vars ---------------------------------------------------------------
var appConf lib.AppConfig

// : functions -----------------------------------------------------------------
func hashPass(pass string) string {
	return base64.StdEncoding.EncodeToString(sha256.New().Sum([]byte(pass)))
}

func createToken(username string, key string) string {
	urc := lib.UserRC{
		User: lib.User{Name: username, Password: ""},
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt: jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(
				func() time.Duration {
					var unit time.Duration
					switch appConf.TokenExpUnit {
					case "milisec":
						unit = time.Millisecond
					case "sec":
						unit = time.Second
					case "min":
						unit = time.Minute
					case "hour":
						unit = time.Hour
					default:
						unit = time.Minute
					}
					return unit * time.Duration(appConf.TokenExp)
				}(),
			)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, urc)
	s, err := token.SignedString([]byte(key))
	if err != nil {
		panic(err)
	}

	return s
}

func extractUsernameToken(c echo.Context) string {
	username := c.
		Get("user").(*jwt.Token).
		Claims.(jwt.MapClaims)["username"].(string)

	return username
}

// : MAIN ----------------------------------------------------------------------
func main() {

	//: app config ------------------------------------------------------------
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

	//: routes ----------------------------------------------------------------
	e.GET("/api/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello World baby!")
	})

	e.GET("/api/:username", func(c echo.Context) error {
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

	//: add movie -------------------------------------------------------------
	e.POST("/api/movie", func(c echo.Context) error {
		username := extractUsernameToken(c)

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
			`INSERT INTO public.tbmovie (username, movie_name, "year")
			VALUES($1, $2, $3);`,
			username,
			sgst.Name, sgst.Year)

		if err != nil {
			return c.String(http.StatusConflict, "duplicate")
		}

		return c.String(200, "added")
	})

	e.DELETE("/api/movie", func(c echo.Context) error {
		username := c.
			Get("user").(*jwt.Token).
			Claims.(jwt.MapClaims)["username"].(string)

		m := new(lib.Movie)
		if err := c.Bind(m); err != nil {
			return c.String(http.StatusNotAcceptable, "not acceptable")
		} else if err := c.Validate(m); err != nil {
			return err
		}

		res, err := db.Exec(`DELETE FROM tbMovie
			WHERE username=$1 AND movie_name=$2 AND "year"=$3;`,
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

	//: sign-up sign-in -------------------------------------------------------
	e.POST("/api/sign-in", func(c echo.Context) error {
		u := new(lib.User)

		if err := c.Bind(u); err != nil {
			return c.String(http.StatusBadRequest, "bad request")
		} else if err := c.Validate(u); err != nil {
			return err
		}

		var dbUser string
		var dbPassHash string
		err := db.
			QueryRow(`SELECT username, passHash 
			FROM tbUser WHERE username = $1`, u.Name).
			Scan(&dbUser, &dbPassHash)

		switch err {
		case sql.ErrNoRows:
			return c.String(http.StatusNotFound, "username or password")
		case nil:
			passHash := hashPass(u.Password)

			if dbPassHash == passHash {

				return c.JSON(http.StatusOK, map[string]string{
					"token": createToken(dbUser, appConf.JwtKey),
				})
			} else {
				return c.String(http.StatusNotFound, "username or password")
			}
		default:
			panic(err)
		}
	})

	/* e.POST("/api/sign-up", func(c echo.Context) error {
		u := new(lib.User)

		if err := c.Bind(u); err != nil {
			return c.String(http.StatusBadRequest, "bad request")
		} else if err := c.Validate(u); err != nil {
			return err
		}

		var dbUser string
		err := db.
			QueryRow(
				"SELECT username FROM tbUser where username = $1", u.Name).
			Scan(&dbUser)

		switch err {
		case sql.ErrNoRows:
			db.Exec(`INSERT INTO tbUser (username, passhash)
			VALUES ($1, $2);`, u.Name, hashPass(u.Password))

			return c.JSON(http.StatusOK, map[string]string{
				"token": createToken(u.Name, appConf.JwtKey),
			})
		case nil:
			return c.String(http.StatusConflict, "another username")
		default:
			panic(err)
		}

	}) */

	//: suggest film ----------------------------------------------------------
	e.POST("/api/suggest", func(c echo.Context) error {
		sgst := new(lib.SearchText)

		if err := c.Bind(sgst); err != nil {
			fmt.Println(err)
			return c.String(http.StatusBadRequest, "bad request")
		} else if err := c.Validate(sgst); err != nil {
			return err
		}

		ms := lib.SearchForMovieName(sgst.Text)

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
	//: start server ----------------------------------------------------------
	e.Logger.Info(e.Start("localhost:2020"))
}
