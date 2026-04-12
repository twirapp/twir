package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/twirapp/twir/apps/twitch-mock/internal/state"
)

type createConduitRequest struct {
	ShardCount int `json:"shard_count"`
}

type updateConduitShardsRequest struct {
	ConduitID string `json:"conduit_id"`
	Shards    []struct {
		ID        int `json:"id"`
		Transport struct {
			Method    string `json:"method"`
			SessionID string `json:"session_id"`
		} `json:"transport"`
	} `json:"shards"`
}

func (s *Server) listConduits(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"data": s.state.ListConduits()})
}

func (s *Server) createConduit(c *gin.Context) {
	var req createConduitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	conduit := s.state.CreateConduit(req.ShardCount)
	c.JSON(http.StatusOK, gin.H{"data": []any{conduit}})
}

func (s *Server) updateConduitShards(c *gin.Context) {
	var req updateConduitShardsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	shards := make([]state.ConduitShard, 0, len(req.Shards))
	for _, shard := range req.Shards {
		shards = append(shards, state.ConduitShard{
			ID: shard.ID,
			Transport: state.ConduitShardTransport{
				Method:    shard.Transport.Method,
				SessionID: shard.Transport.SessionID,
			},
		})
	}

	if err := s.state.UpdateConduitShards(req.ConduitID, shards); err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, state.ErrConduitNotFound) {
			status = http.StatusNotFound
		}

		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}
