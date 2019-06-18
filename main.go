package main

import (
	"net/http"

	"github.com/elena-kolevska/velociraptor/middleware/twilio"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func twilioAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if releaseStage != "localhost" {
			scheme := "http"
			if c.IsTLS() {
				scheme = "https"
			}

			formParameters, err := c.FormParams()
			if err != nil {
				return c.String(http.StatusUnauthorized, "Something's wrong with your request")
			}

			if twilio.IsValidTwilioSignature(scheme, c.Request().Host, twilioAPIToken, c.Request().RequestURI, formParameters, c.Request().Header.Get("X-Twilio-Signature")) == false {
				return c.String(http.StatusUnauthorized, "You're not Twilio")
			}
		}

		return next(c)
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
	e.Logger.Fatal(e.StartTLS(":"+port, certFile, keyFile))
}
