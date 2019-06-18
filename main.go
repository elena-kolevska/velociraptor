package main

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	// "github.com/gomodule/redigo/redis"
)

func twilioAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		return next(c)
		//return c.String(http.StatusUnauthorized, "You're not Twilio")
	}
}

func main() {
	e := echo.New()

	e.Pre(middleware.AddTrailingSlash())
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `[${time_rfc3339}] ${status} ${method} ${host}${path} ${latency_human}` + "\n",
	}))
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	twilioGroup := e.Group("/twilio")
	twilioGroup.Use(twilioAuth)

	twilioGroup.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, Twilio!")
	})
	e.Logger.Fatal(e.Start(":1323"))
}
