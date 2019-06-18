package main

import (
	"github.com/elena-kolevska/velociraptor/middleware/twilio"
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	//"github.com/gomodule/redigo/redis"
)

//var (
//	pool      = newPool()
//	redisConn = pool.Get()
//)

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
	twilioGroup.Use(twilio.Auth)

	twilioGroup.POST("/", func(c echo.Context) error {
		// Parse request
		//
		return c.String(http.StatusOK, "Hello, Twilio!")
	})
	e.Logger.Fatal(e.StartTLS(":"+port, certFile, keyFile))
}


//func newPool() *redis.Pool {
//	return &redis.Pool{
//		MaxIdle:   20,
//		MaxActive: 1000, // max number of connections
//		Dial: func() (redis.Conn, error) {
//			c, err := redis.Dial("tcp", redisHost+":"+redisPort)
//			if err != nil {
//				panic(err.Error())
//			}
//			return c, err
//		},
//	}
//}