package twitchactions

import (
	"context"
	"net/http"
	"time"

	"github.com/nicklaw5/helix/v2"
	"github.com/twirapp/twir/libs/twitch"
)

// botClientCacheTTL bounds how long a cached bot Helix client is reused.
// Caching keeps moderation actions working through tokens-service hiccups
// (no bus roundtrip per action) and removes per-call latency; a stale client
// that starts failing with 401 is evicted and recreated once per call.
const botClientCacheTTL = 15 * time.Minute

type cachedBotClient struct {
	client    *helix.Client
	expiresAt time.Time
}

func (c *TwitchActions) getBotClient(ctx context.Context, botID string) (*helix.Client, error) {
	if c.newBotClient != nil {
		return c.newBotClient(ctx, botID)
	}

	c.botClientsMu.Lock()
	cached, ok := c.botClients[botID]
	c.botClientsMu.Unlock()

	if ok && time.Now().Before(cached.expiresAt) {
		return cached.client, nil
	}

	client, err := twitch.NewBotClientWithContext(ctx, botID, c.config, c.twirBus)
	if err != nil {
		return nil, err
	}

	c.botClientsMu.Lock()
	if c.botClients == nil {
		c.botClients = make(map[string]cachedBotClient)
	}
	c.botClients[botID] = cachedBotClient{
		client:    client,
		expiresAt: time.Now().Add(botClientCacheTTL),
	}
	c.botClientsMu.Unlock()

	return client, nil
}

func (c *TwitchActions) evictBotClient(botID string) {
	c.botClientsMu.Lock()
	delete(c.botClients, botID)
	c.botClientsMu.Unlock()
}

func (c *TwitchActions) withBotClient(
	ctx context.Context,
	botID string,
	call func(client *helix.Client) (statusCode int, err error),
) error {
	client, err := c.getBotClient(ctx, botID)
	if err != nil {
		return err
	}

	statusCode, err := call(client)
	if err == nil || statusCode != http.StatusUnauthorized {
		return err
	}

	c.evictBotClient(botID)

	client, err = c.getBotClient(ctx, botID)
	if err != nil {
		return err
	}

	_, err = call(client)

	return err
}
