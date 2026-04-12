package handlers

import (
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/twirapp/twir/apps/twitch-mock/internal/config"
)

func (s *Server) token(c *gin.Context) {
	grantType := c.DefaultQuery("grant_type", c.PostForm("grant_type"))

	switch grantType {
	case "client_credentials":
		c.JSON(http.StatusOK, gin.H{
			"access_token": config.MockAppToken,
			"token_type":   "bearer",
			"expires_in":   99999999,
			"scope":        []string{},
		})
	case "authorization_code", "refresh_token":
		c.JSON(http.StatusOK, gin.H{
			"access_token":  config.MockUserToken,
			"token_type":    "bearer",
			"expires_in":    99999999,
			"scope":         []string{},
			"refresh_token": "mock-user-refresh",
		})
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "unsupported grant_type"})
	}
}

func (s *Server) authorize(c *gin.Context) {
	redirect, err := url.Parse(s.config.SiteBaseURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	redirect = redirect.JoinPath("login")
	query := redirect.Query()
	query.Set("code", "mock_code_1234567890")
	if stateValue := c.Query("state"); stateValue != "" {
		query.Set("state", stateValue)
	}
	redirect.RawQuery = query.Encode()

	c.Redirect(http.StatusFound, redirect.String())
}

func (s *Server) validate(c *gin.Context) {
	token := authToken(c.GetHeader("Authorization"))

	response := gin.H{
		"client_id":  "mock",
		"login":      config.MockBroadcasterLogin,
		"user_id":    config.MockBroadcasterID,
		"expires_in": 99999999,
		"scopes":     []string{},
	}

	if token == config.MockBotToken {
		response["login"] = config.MockBotLogin
		response["user_id"] = config.MockBotID
	}

	c.JSON(http.StatusOK, response)
}
