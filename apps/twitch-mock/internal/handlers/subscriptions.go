package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/twirapp/twir/apps/twitch-mock/internal/state"
)

type createSubscriptionRequest struct {
	Type      string         `json:"type"`
	Version   string         `json:"version"`
	Condition map[string]any `json:"condition"`
	Transport struct {
		Method    string `json:"method"`
		SessionID string `json:"session_id"`
	} `json:"transport"`
}

func (s *Server) createSubscription(c *gin.Context) {
	var req createSubscriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	subscription := s.state.CreateSubscription(
		req.Type,
		req.Version,
		req.Condition,
		state.SubscriptionTransport{
			Method:    req.Transport.Method,
			SessionID: req.Transport.SessionID,
		},
	)

	c.JSON(http.StatusAccepted, gin.H{
		"data":           []any{subscription},
		"total":          1,
		"total_cost":     1,
		"max_total_cost": 10000,
	})
}

func (s *Server) deleteSubscription(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		id = c.Query("id")
	}

	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "subscription id is required"})
		return
	}

	s.state.DeleteSubscription(id)
	c.Status(http.StatusNoContent)
}
