package webhooks

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	model "github.com/satont/twir/libs/gomodels"
)

type donatelloBody struct {
	PubId       string `json:"pubId"`
	ClientName  string `json:"clientName"`
	Message     string `json:"message"`
	Amount      string `json:"amount"`
	Currency    string `json:"currency"`
	Source      string `json:"source"`
	Goal        string `json:"goal"`
	IsPublished bool   `json:"isPublished"`
	CreatedAt   string `json:"createdAt"`
}

func (c *Webhooks) donatelloHandler(g *gin.Context) {
	apiKey := g.GetHeader("X-Key")
	if apiKey == "" {
		g.JSON(http.StatusBadRequest, gin.H{"error": "X-Key header is required"})
		return
	}

	integration := &model.ChannelsIntegrations{}
	if err := c.db.
		WithContext(g.Request.Context()).
		Where(`"id" = ?`, apiKey).
		First(integration).
		Error; err != nil {
		g.JSON(http.StatusNotFound, gin.H{"error": "Integration not found"})
		return
	}

	body := &donatelloBody{}
	if err := g.BindJSON(body); err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	integrationsMessage := &pbMessage{
		TwitchUserId: integration.ChannelID,
		Amount:       body.Amount,
		Currency:     body.Currency,
		Message:      body.Message,
		UserName:     body.ClientName,
	}
	integrationsNameBytes, err := json.Marshal(integrationsMessage)
	if err != nil {
		c.logger.Error("cannot marshal message", slog.Any("err", err))
	} else {
		c.pubSub.Publish("donations:new", integrationsNameBytes)
	}

	g.JSON(http.StatusOK, "ok")
}
