package requests

import (
	"github.com/elena-kolevska/velociraptor/clients"
	"github.com/elena-kolevska/velociraptor/config"
	"github.com/gomodule/redigo/redis"
	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
	"strconv"
	"time"
)

//WebhookPayload
func (w *Webhook) ParsePayload(c echo.Context) {
	timeFormat := "2006-01-02T15:04:05.000Z"

	var eventTypeMap = map[string]string{
		"onMessageSend":      "onMessageSend",
		"onMessageSent":      "onMessageSent",
		"onMessageUpdate":    "onMessageUpdate",
		"onMessageUpdated":   "onMessageUpdated",
		"onMessageRemove":    "onMessageRemove",
		"onMessageRemoved":   "onMessageRemoved",
		"onChannelAdd":       "onChannelAdd",
		"onChannelAdded":     "onChannelAdded",
		"onChannelDestroy":   "onChannelDestroy",
		"onChannelDestroyed": "onChannelDestroyed",
		"onChannelUpdate":    "onChannelUpdate",
		"onChannelUpdated":   "onChannelUpdated",
		"onMemberAdd":        "onMemberAdd",
		"onMemberAdded":      "onMemberAdded",
		"onMemberRemove":     "onMemberRemove",
		"onMemberRemoved":    "onMemberRemoved",
		"onUserUpdate":       "onUserUpdate",
		"onUserUpdated":      "onUserUpdated",
	}

	eventType, ok := eventTypeMap[c.FormValue("EventType")]
	if !ok {
		eventType = "unknown"
	}

	from, _ := strconv.ParseInt(c.FormValue("From"), 0, 64) // Convert the service-provided timestamp to a unix timestamp

	w.from = from
	w.EventType = eventType
	w.source = "twilio"
	w.channelSid = c.FormValue("ChannelSid")

	w.setTimestamp(timeFormat, c.FormValue("DateCreated"))
}

func (w *Webhook) HandleOnMessageSent(r redis.Conn, c clients.Client) {

	key := "last-activity:" + w.channelSid

	val, _ := r.Do("SET", key, 1, "NX", "EX", config.ApiRefreshRate)

	if val == "OK" {
		err := c.UpdateLastActivity(&w.channelSid, &w.timestamp)
		if err != nil {
			log.Print(err)
		}
	}
}

func (w *Webhook) setTimestamp(format string, value string) {
	t, _ := time.Parse(format, value)
	w.timestamp = t.Unix()
}
