package requests

import (
	"github.com/elena-kolevska/velociraptor/clients"
	"github.com/gomodule/redigo/redis"
	"github.com/labstack/echo"
)

type Webhook struct {
	eventType  string
	from       int64
	timestamp  int64
	channelSid string
	source     string
}

type WebhookInterface interface {
	ParsePayload(echo.Context) Webhook
	HandleOnMessageSent(redis.Conn, clients.Client)
}
