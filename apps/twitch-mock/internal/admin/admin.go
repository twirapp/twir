package admin

import (
	"embed"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/twirapp/twir/apps/twitch-mock/internal/state"
	twitchws "github.com/twirapp/twir/apps/twitch-mock/internal/websocket"
)

//go:embed templates/index.html
var templateFS embed.FS

type Server struct {
	state  *state.State
	logger *slog.Logger
	ws     *twitchws.Server
	engine *gin.Engine
}

func New(appState *state.State, logger *slog.Logger, ws *twitchws.Server) *Server {
	engine := gin.New()
	engine.Use(gin.Recovery())

	server := &Server{
		state:  appState,
		logger: logger,
		ws:     ws,
		engine: engine,
	}

	engine.GET("/admin", server.index)
	engine.GET("/admin/state/stream", server.getStreamState)
	engine.POST("/admin/state/stream", server.setStreamState)
	engine.GET("/admin/state/follows", server.getFollowState)
	engine.POST("/admin/state/follows", server.setFollowState)
	engine.GET("/admin/state/moderators", server.getModeratorState)
	engine.POST("/admin/state/moderators", server.setModeratorState)
	engine.POST("/admin/trigger/:event", server.trigger)

	return server
}

func (s *Server) Handler() http.Handler {
	return s.engine
}

func (s *Server) index(c *gin.Context) {
	content, err := templateFS.ReadFile("templates/index.html")
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.Data(http.StatusOK, "text/html; charset=utf-8", content)
}

func (s *Server) trigger(c *gin.Context) {
	eventType := c.Param("event")
	payload, err := decodeBody(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	switch eventType {
	case "stream.online":
		s.state.SetStreamOnline(true)
	case "stream.offline":
		s.state.SetStreamOnline(false)
	}

	s.logger.Info("admin triggered mock event", slog.String("event_type", eventType), slog.Any("payload", payload))
	s.ws.Broadcast(eventType, payload)

	c.JSON(http.StatusOK, gin.H{"status": "ok", "event_type": eventType})
}

func (s *Server) getStreamState(c *gin.Context) {
	title, gameID, gameName, viewerCount := s.state.StreamMeta()
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, gin.H{
		"online":       s.state.IsStreamOnline(),
		"title":        title,
		"game_id":      gameID,
		"game_name":    gameName,
		"viewer_count": viewerCount,
	})
}

func (s *Server) setStreamState(c *gin.Context) {
	var body struct {
		Online      bool   `json:"online"`
		Title       string `json:"title"`
		GameID      string `json:"game_id"`
		GameName    string `json:"game_name"`
		ViewerCount int    `json:"viewer_count"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.Header("Content-Type", "application/json")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	s.state.SetStreamMeta(body.Title, body.GameID, body.GameName, body.ViewerCount)
	s.state.SetStreamOnline(body.Online)
	s.getStreamState(c)
}

func (s *Server) getFollowState(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, gin.H{
		"total":             s.state.FollowersTotal(),
		"followed_user_ids": s.state.ListFollowedUserIDs(),
	})
}

func (s *Server) setFollowState(c *gin.Context) {
	var body struct {
		Total           int      `json:"total"`
		FollowedUserIDs []string `json:"followed_user_ids"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.Header("Content-Type", "application/json")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	existing := s.state.ListFollowedUserIDs()
	for _, userID := range existing {
		s.state.SetUserFollowed(userID, false)
	}

	for _, userID := range normalizeIDs(body.FollowedUserIDs) {
		s.state.SetUserFollowed(userID, true)
	}

	s.state.SetFollowersTotal(body.Total)
	s.getFollowState(c)
}

func (s *Server) getModeratorState(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, gin.H{
		"moderator_ids": s.state.ListModerators(),
	})
}

func (s *Server) setModeratorState(c *gin.Context) {
	var body struct {
		ModeratorIDs []string `json:"moderator_ids"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.Header("Content-Type", "application/json")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	existing := s.state.ListModerators()
	for _, userID := range existing {
		s.state.SetModerator(userID, false)
	}

	for _, userID := range normalizeIDs(body.ModeratorIDs) {
		s.state.SetModerator(userID, true)
	}

	s.getModeratorState(c)
}

func decodeBody(body io.Reader) (map[string]any, error) {
	var payload map[string]any
	if err := json.NewDecoder(body).Decode(&payload); err != nil {
		if err == io.EOF {
			return map[string]any{}, nil
		}

		return nil, err
	}

	return payload, nil
}

func normalizeIDs(values []string) []string {
	seen := map[string]time.Time{}
	result := make([]string, 0, len(values))
	for _, value := range values {
		id := strings.TrimSpace(value)
		if id == "" {
			continue
		}

		if _, exists := seen[id]; exists {
			continue
		}

		seen[id] = time.Now().UTC()
		result = append(result, id)
	}

	sort.Strings(result)
	return result
}
