package twilioHandlers

import (
	"github.com/elena-kolevska/velociraptor/clients"
	"github.com/elena-kolevska/velociraptor/requests"
	"github.com/gomodule/redigo/redis"
	"github.com/labstack/echo"
	"net/http"
)

func HandleWebhook(c echo.Context, redisConn redis.Conn, c clients.Client) error {
	t := requests.Webhook{}
	t.ParsePayload(c)
	t.HandleOnMessageSent(redisConn, &c)
	return c.String(http.StatusOK, "Thanks, Twilio!")
}
