package generalMiddleware

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func General(e *echo.Echo) {
	e.Pre(middleware.AddTrailingSlash())
	e.Pre(middleware.Recover())
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `[${time_rfc3339}] ${status} ${method} ${host}${path} ${latency_human}` + "\n",
	}))
}
