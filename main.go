package main

import (
	"github.com/elena-kolevska/velociraptor/clients"
	"github.com/elena-kolevska/velociraptor/config"
	"github.com/elena-kolevska/velociraptor/middleware/twilio"
	"log"
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

	log.Println("Acquiring access token from the destination API...")
	client := clients.ModelLifeClient{
		BaseURL:               config.ApiBaseURL,
		GetTokenURL:           config.ApiTokenPath,
		UpdateConversationURL: config.ApiUpdateConversationPath,
		ClientID:              config.ApiClientID,
		ClientSecret:          config.ApiClientSecret,
		HttpClient:            http.Client{},
	}
	err := client.GetAccessToken()
	if err != nil {
		log.Fatal("We couldn't get an access token from the remote API")
	}
	log.Println("Done.")
	log.Println("Starting server...")

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
		return c.String(http.StatusOK, "Hello, Twilio!")
	})

	e.Logger.Fatal(e.StartTLS(":"+config.EnvPort, config.EnvCertFile, config.EnvKeyFile))
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
