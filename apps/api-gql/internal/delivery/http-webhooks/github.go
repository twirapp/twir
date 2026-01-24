package http_webhooks

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/twirapp/twir/apps/api-gql/internal/services/webhook_notifications"
	"github.com/twirapp/twir/libs/logger"
)

func (c *Webhooks) githubWebhookHandler(g *gin.Context) {
	event := g.GetHeader("X-GitHub-Event")
	if event == "" {
		g.JSON(http.StatusBadRequest, gin.H{"error": "Missing X-GitHub-Event header"})
		return
	}

	payload, err := g.GetRawData()
	if err != nil {
		c.logger.Error("cannot read github webhook payload", logger.Error(err))
		g.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
		return
	}

	if err := c.verifyGithubSignature(g, payload); err != nil {
		c.logger.Warn("invalid github webhook signature", logger.Error(err))
		g.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid signature"})
		return
	}

	payload = normalizeGithubPayload(g, payload)

	if err := c.webhookNotificationsService.HandleGithubWebhook(
		g.Request.Context(),
		event,
		payload,
	); err != nil {
		switch {
		case errors.Is(err, webhook_notifications.ErrInvalidPayload):
			g.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
		default:
			c.logger.Error("failed to handle github webhook", logger.Error(err))
			g.JSON(http.StatusInternalServerError, gin.H{"error": "Internal error"})
		}
		return
	}

	g.Status(http.StatusOK)
}

func (c *Webhooks) verifyGithubSignature(g *gin.Context, payload []byte) error {
	secret := strings.TrimSpace(c.config.GithubWebhookSecret)
	if secret == "" {
		return nil
	}

	signature := strings.TrimSpace(g.GetHeader("X-Hub-Signature-256"))
	if signature != "" {
		expected := githubSignatureSHA256(secret, payload)
		if hmac.Equal([]byte(expected), []byte(signature)) {
			return nil
		}
		return errors.New("sha256 signature mismatch")
	}

	legacy := strings.TrimSpace(g.GetHeader("X-Hub-Signature"))
	if legacy == "" {
		return errors.New("missing signature headers")
	}

	expectedLegacy := githubSignatureSHA1(secret, payload)
	if hmac.Equal([]byte(expectedLegacy), []byte(legacy)) {
		return nil
	}

	return errors.New("sha1 signature mismatch")
}

func githubSignatureSHA256(secret string, payload []byte) string {
	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = mac.Write(payload)
	return "sha256=" + hex.EncodeToString(mac.Sum(nil))
}

func githubSignatureSHA1(secret string, payload []byte) string {
	mac := hmac.New(sha1.New, []byte(secret))
	_, _ = mac.Write(payload)
	return "sha1=" + hex.EncodeToString(mac.Sum(nil))
}

func normalizeGithubPayload(g *gin.Context, payload []byte) []byte {
	contentType := g.GetHeader("Content-Type")
	if !strings.Contains(contentType, "application/x-www-form-urlencoded") &&
		!bytes.HasPrefix(payload, []byte("payload=")) {
		return payload
	}

	parsed, err := url.ParseQuery(string(payload))
	if err != nil {
		return payload
	}

	raw := parsed.Get("payload")
	if raw == "" {
		return payload
	}

	return []byte(raw)
}
