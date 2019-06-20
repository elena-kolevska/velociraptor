package main

import (
	"github.com/elena-kolevska/velociraptor/clients"
	"github.com/elena-kolevska/velociraptor/config"
	generalMiddleware "github.com/elena-kolevska/velociraptor/middleware"
	"github.com/elena-kolevska/velociraptor/middleware/twilio"
	"github.com/elena-kolevska/velociraptor/requests"
	"log"
	"net/http"

	"github.com/gomodule/redigo/redis"
	"github.com/labstack/echo"
)

var (
	pool      = newPool()
	redisConn = pool.Get()
)

func main() {

	// Redis Setup
	defer redisConn.Close()
	authenticateWithRedis()

	log.Println("Acquiring access token from the destination API...")
	client, err := clients.InstantiateClient()
	if err != nil {
		panic("We couldn't get an access token from the remote API")
	}
	log.Println("Starting server...")

	//// SERVER ////
	e := echo.New()
	e.HideBanner = true

	generalMiddleware.General(e)

	twilioGroup := e.Group("/twilio")
	twilioGroup.Use(twilio.Auth)
	twilioGroup.POST("/", func(c echo.Context) error {
		t := requests.Webhook{}
		t.ParsePayload(c)
		switch t.EventType {
		case "onMessageSent":
			go t.HandleOnMessageSent(redisConn, client) // Execute call to the remote API in a go routine (async)
			return c.String(http.StatusOK, "Thanks, Twilio!")
		default:
			return c.String(http.StatusOK, "Unsupported webhook type")
		}
	})

	e.Logger.Fatal(e.StartTLS(":"+config.EnvPort, config.EnvCertFile, config.EnvKeyFile))
}

func newPool() *redis.Pool {
	log.Println("Connecting to Redis...")
	return &redis.Pool{
		MaxIdle:   20,
		MaxActive: 1000, // max number of connections
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", config.RedisHost+":"+config.RedisPort)
			if err != nil {
				panic(err.Error())
			}
			return c, err
		},
	}
}

func authenticateWithRedis () {
	if config.RedisAuth != "" {
		log.Println("Authenticating with Redis...")

		_, err := redisConn.Do("AUTH", config.RedisAuth)
		if err != nil {
			panic(err.Error())
		}
	}

}