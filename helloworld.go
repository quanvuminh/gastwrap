package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/quanvuminh/gastwrap/users"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/facebookgo/grace/gracehttp"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
)

const (
	jwtsecret string = "VeryStrongSecret"
)

// LoginUser define an user trying to login
type LoginUser struct {
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
}

// Login and receive JSON access token
func login(ctx echo.Context) error {
	var login LoginUser
	if err := ctx.Bind(&login); err != nil {
		return err
	}

	if login.Username == "admin" && login.Password == "VeryStrongPassword" {
		// Create token
		token := jwt.New(jwt.SigningMethodHS256)

		// Set claims
		claims := token.Claims.(jwt.MapClaims)
		claims["name"] = "Admin"
		claims["admin"] = true
		claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

		// Generate encoded token and send it as response.
		t, err := token.SignedString([]byte(jwtsecret))
		if err != nil {
			return err
		}
		return ctx.JSON(http.StatusOK, map[string]string{
			"token": t,
		})
	}

	return echo.ErrUnauthorized
}

func accessible(ctx echo.Context) error {
	return ctx.String(http.StatusOK, "Accessible")
}

func restricted(ctx echo.Context) error {
	user := ctx.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	name := claims["name"].(string)
	return ctx.String(http.StatusOK, "Welcome "+name+"!")
}

func main() {
	listenAddr := flag.String("listen", ":8080", "Listening IP Address") // Listen address for the server
	prefix := flag.String("prefix", "", "Required - Prefix for endpoints. E.g. http://domain.tld/prefix/endpoint")
	logdir := flag.String("logdir", "/tmp", "Log file directory")
	flag.Parse()

	// Endpoint prefix for multi instances
	if *prefix == "" {
		fmt.Println("Flag 'prefix' is required. Run with -h for details")
		os.Exit(2)
	}

	// Path of log file
	logfile := *logdir + "/gastwrap.log"
	f, err := os.OpenFile(logfile, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
	if err != nil {
		fmt.Println("Can not write log file")
		os.Exit(2)
	}
	defer f.Close()

	// New instance
	app := echo.New()

	// Customization
	app.Logger.SetOutput(f)
	app.Logger.SetLevel(log.DEBUG)

	// Middlewares
	app.Pre(middleware.RemoveTrailingSlash()) // Remove trailing slash to the request URI
	app.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE},
	})) // CORS
	app.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Output: f,
	})) // Logs the information about each HTTP request
	app.Use(middleware.Recover())   // Recovers from panics anywhere
	app.Use(middleware.RequestID()) // Unique id for each request
	app.Use(middleware.Secure())    // Protection against XSS attack and other code injection attacks.

	// Login route
	app.POST("/login", login)

	// Restricted group
	r := app.Group("/" + *prefix)
	r.Use(middleware.JWT([]byte(jwtsecret)))
	r.GET("", restricted)

	r.POST("/users/new", func(ctx echo.Context) error {

		var newuser users.User
		if err := ctx.Bind(&newuser); err != nil {
			return err
		}

		err := users.Create(newuser)
		if err != nil {
			return ctx.JSONPretty(http.StatusOK, echo.Map{
				"status": "Fail: " + err.Error(),
			}, "  ")
		}

		return ctx.JSONPretty(http.StatusOK, echo.Map{
			"status": "Success",
		}, "  ")
	})

	r.GET("/users/:id", func(ctx echo.Context) error {
		uid := ctx.Param("id")
		u, err := users.Get(uid)
		if err != nil {
			return ctx.HTML(http.StatusNotFound, "Error: "+err.Error())
		}

		return ctx.JSONPretty(http.StatusOK, u, "  ")
	})

	r.DELETE("/users/:id", func(ctx echo.Context) error {
		uid := ctx.Param("id")
		err := users.Delete(uid)
		if err != nil {
			return ctx.JSONPretty(http.StatusOK, echo.Map{
				"status": "Fail: " + err.Error(),
			}, "  ")
		}

		return ctx.JSONPretty(http.StatusOK, echo.Map{
			"status": "Success",
		}, "  ")
	})

	r.PUT("/users/:id", func(ctx echo.Context) error {
		uid := ctx.Param("id")
		var newuser users.User
		if err := ctx.Bind(&newuser); err != nil {
			return err
		}

		err := users.Update(uid, newuser)
		if err != nil {
			return ctx.JSONPretty(http.StatusOK, echo.Map{
				"status": "Fail: " + err.Error(),
			}, "  ")
		}

		return ctx.JSONPretty(http.StatusOK, echo.Map{
			"status": "Success",
		}, "  ")
	})

	app.Server.Addr = *listenAddr

	// Serve it like a boss
	app.Logger.Fatal(gracehttp.Serve(app.Server))
}
