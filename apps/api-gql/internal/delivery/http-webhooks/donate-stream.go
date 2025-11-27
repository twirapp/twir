package http_webhooks

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	model "github.com/twirapp/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/logger"
)

type donateStreamIncomingData struct {
	Type     string `json:"type,omitempty"`
	Uid      string `json:"uid"`
	Message  string `json:"message"`
	Sum      string `json:"sum"`
	Nickname string `json:"nickname"`
}

func (c *Webhooks) donateStreamHandler(g *gin.Context) {
	integration := model.ChannelsIntegrations{}
	id := g.Param("id")

	if err := c.db.
		WithContext(g.Request.Context()).
		Where("id = ?", id).First(&integration).Error; err != nil {
		g.JSON(http.StatusNotFound, gin.H{"error": "Integration not found"})
		return
	}

	body := &donateStreamIncomingData{}
	if err := g.BindJSON(body); err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
		return
	}

	if body.Type == "confirm" {
		value, err := c.kv.Get(
			g.Request.Context(),
			"donate_stream_confirmation"+integration.ID,
		).String()
		if err != nil {
			c.logger.Error("cannot get confirmation from kv", logger.Error(err))
			g.JSON(http.StatusInternalServerError, gin.H{"error": "Internal error"})
			return
		}

		integrationsMessage := &pbMessage{
			TwitchUserId: integration.ChannelID,
			Amount:       body.Sum,
			Currency:     "RUB",
			Message:      body.Message,
			UserName:     body.Nickname,
		}
		integrationsNameBytes, err := json.Marshal(integrationsMessage)
		if err != nil {
			c.logger.Error("cannot marshal message", logger.Error(err))
		} else {
			c.pubSub.Publish("donations:new", integrationsNameBytes)
		}

		g.String(http.StatusOK, value)
		return
	}

	g.String(http.StatusOK, "ok")
}
