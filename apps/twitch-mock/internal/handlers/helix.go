package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/twirapp/twir/apps/twitch-mock/internal/config"
	"github.com/twirapp/twir/apps/twitch-mock/internal/state"
)

type Server struct {
	config *config.Config
	state  *state.State
	logger *slog.Logger
	engine *gin.Engine
}

func New(cfg *config.Config, appState *state.State, logger *slog.Logger) *Server {
	engine := gin.New()
	engine.Use(gin.Recovery())

	server := &Server{
		config: cfg,
		state:  appState,
		logger: logger,
		engine: engine,
	}

	server.registerRoutes()

	return server
}

func (s *Server) Handler() http.Handler {
	return s.engine
}

func (s *Server) registerRoutes() {
	s.engine.GET("/health", s.health)
	s.engine.POST("/oauth2/token", s.token)
	s.engine.GET("/oauth2/authorize", s.authorize)
	s.engine.GET("/oauth2/validate", s.validate)

	helix := s.engine.Group("/helix")
	helix.GET("/users", s.users)
	helix.GET("/streams", s.streams)
	helix.GET("/channels", s.channels)
	helix.GET("/channels/followers", s.channelFollowers)
	helix.GET("/eventsub/conduits", s.listConduits)
	helix.POST("/eventsub/conduits", s.createConduit)
	helix.PATCH("/eventsub/conduits/shards", s.updateConduitShards)
	helix.POST("/eventsub/subscriptions", s.createSubscription)
	helix.DELETE("/eventsub/subscriptions", s.deleteSubscription)
	helix.DELETE("/eventsub/subscriptions/:id", s.deleteSubscription)
	helix.POST("/chat/messages", s.chatMessages)
	helix.GET("/moderation/moderators", s.moderators)
	helix.POST("/moderation/bans", s.moderationBans)

	s.engine.NoRoute(func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/helix/") {
			s.defaultHelix(c)
			return
		}

		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
	})
}

func (s *Server) users(c *gin.Context) {
	user := mockBroadcasterUser()
	if shouldUseBotUser(c) {
		user = mockBotUser()
	}

	c.JSON(http.StatusOK, gin.H{"data": []gin.H{user}})
}

func (s *Server) streams(c *gin.Context) {
	stream, online := s.state.StreamSnapshot()
	if !online {
		c.JSON(http.StatusOK, gin.H{"data": []any{}})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": []any{stream}})
}

func (s *Server) channels(c *gin.Context) {
	title, gameID, gameName, _ := s.state.StreamMeta()

	c.JSON(http.StatusOK, gin.H{
		"data": []gin.H{
			{
				"broadcaster_id":       config.MockBroadcasterID,
				"broadcaster_login":    config.MockBroadcasterLogin,
				"broadcaster_name":     config.MockBroadcasterName,
				"broadcaster_language": "en",
				"game_id":              gameID,
				"game_name":            gameName,
				"title":                title,
				"delay":                0,
			},
		},
	})
}

func (s *Server) channelFollowers(c *gin.Context) {
	userID := strings.TrimSpace(c.Query("user_id"))
	if userID != "" {
		followedAt, ok := s.state.IsUserFollowed(userID)
		if !ok {
			c.JSON(http.StatusOK, gin.H{"data": []any{}, "total": s.state.FollowersTotal(), "pagination": gin.H{}})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data":       []gin.H{buildFollowEntry(userID, followedAt)},
			"total":      s.state.FollowersTotal(),
			"pagination": gin.H{},
		})
		return
	}

	followedIDs := s.state.ListFollowedUserIDs()
	if first, err := strconv.Atoi(c.DefaultQuery("first", "0")); err == nil && first > 0 && first < len(followedIDs) {
		followedIDs = followedIDs[:first]
	}

	data := make([]gin.H, 0, len(followedIDs))
	for _, followedID := range followedIDs {
		followedAt, ok := s.state.IsUserFollowed(followedID)
		if !ok {
			continue
		}

		data = append(data, buildFollowEntry(followedID, followedAt))
	}

	c.JSON(http.StatusOK, gin.H{
		"data":       data,
		"total":      s.state.FollowersTotal(),
		"pagination": gin.H{},
	})
}

func (s *Server) moderators(c *gin.Context) {
	moderatorIDs := s.state.ListModerators()
	data := make([]gin.H, 0, len(moderatorIDs))

	for _, moderatorID := range moderatorIDs {
		if moderatorID == config.MockBotID {
			bot := mockBotUser()
			data = append(data, gin.H{
				"user_id":    bot["id"],
				"user_login": bot["login"],
				"user_name":  bot["display_name"],
			})
			continue
		}

		data = append(data, gin.H{
			"user_id":    moderatorID,
			"user_login": "mod_" + strings.ToLower(moderatorID),
			"user_name":  "Mod_" + titleCaseIdentifier(moderatorID),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"data":       data,
		"pagination": gin.H{},
	})
}

func (s *Server) chatMessages(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	s.logger.Info("chat message received", slog.String("body", string(body)))
	c.Request.Body = io.NopCloser(bytes.NewReader(body))

	c.JSON(http.StatusOK, gin.H{
		"data": []gin.H{{
			"message_id": "mock-msg-id",
			"is_sent":    true,
		}},
	})
}

func (s *Server) moderationBans(c *gin.Context) {
	var body struct {
		Data struct {
			UserID string `json:"user_id"`
		} `json:"data"`
	}

	_ = c.ShouldBindJSON(&body)
	if body.Data.UserID == "" {
		body.Data.UserID = "mock"
	}

	c.JSON(http.StatusOK, gin.H{
		"data": []gin.H{{
			"broadcaster_id": config.MockBroadcasterID,
			"moderator_id":   config.MockBotID,
			"user_id":        body.Data.UserID,
			"end_time":       nil,
		}},
	})
}

func (s *Server) defaultHelix(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"data":       []any{},
		"total":      0,
		"pagination": gin.H{},
	})
}

func shouldUseBotUser(c *gin.Context) bool {
	token := authToken(c.GetHeader("Authorization"))
	if token == config.MockBotToken {
		return true
	}

	if c.Query("id") == config.MockBotID || c.Query("login") == config.MockBotLogin {
		return true
	}

	return false
}

func authToken(header string) string {
	parts := strings.SplitN(header, " ", 2)
	if len(parts) != 2 {
		return strings.TrimSpace(header)
	}

	return strings.TrimSpace(parts[1])
}

func mockBroadcasterUser() gin.H {
	return gin.H{
		"id":                config.MockBroadcasterID,
		"login":             config.MockBroadcasterLogin,
		"display_name":      config.MockBroadcasterName,
		"type":              "",
		"broadcaster_type":  "partner",
		"description":       "Mock streamer",
		"profile_image_url": "",
		"offline_image_url": "",
		"view_count":        0,
		"created_at":        "2020-01-01T00:00:00Z",
	}
}

func mockBotUser() gin.H {
	return gin.H{
		"id":                config.MockBotID,
		"login":             config.MockBotLogin,
		"display_name":      config.MockBotName,
		"type":              "",
		"broadcaster_type":  "",
		"description":       "Mock bot",
		"profile_image_url": "",
		"offline_image_url": "",
		"view_count":        0,
		"created_at":        time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC).Format(time.RFC3339),
	}
}

func decodeBodyMap(c *gin.Context) (map[string]any, error) {
	var payload map[string]any
	if err := json.NewDecoder(c.Request.Body).Decode(&payload); err != nil {
		if err == io.EOF {
			return map[string]any{}, nil
		}

		return nil, err
	}

	return payload, nil
}

func buildFollowEntry(userID string, followedAt time.Time) gin.H {
	if followedAt.IsZero() {
		followedAt = time.Now().UTC()
	}

	return gin.H{
		"user_id":     userID,
		"user_name":   titleCaseIdentifier(userID),
		"user_login":  strings.ToLower(userID),
		"followed_at": followedAt.Format(time.RFC3339),
	}
}

func titleCaseIdentifier(value string) string {
	if value == "" {
		return ""
	}

	runes := []rune(strings.ToLower(value))
	capitalizeNext := true
	for i, r := range runes {
		if capitalizeNext && r >= 'a' && r <= 'z' {
			runes[i] = r - ('a' - 'A')
			capitalizeNext = false
			continue
		}

		capitalizeNext = r == '_' || r == '-' || r == ' '
	}

	return string(runes)
}
