package manager

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strconv"

	"github.com/kvizyx/twitchy/eventsub"
	"github.com/twirapp/twir/libs/logger"
	twitchconduits "github.com/twirapp/twir/libs/repositories/twitch_conduits"
)

type conduitsResponse struct {
	Data []conduitsResponseConduit `json:"data"`
}

type conduitsResponseConduit struct {
	Id         string `json:"id"`
	ShardCount int    `json:"shard_count"`
}

type createConduitRequest struct {
	ShardCount int `json:"shard_count"`
}

type createConduitResponse struct {
	Data []conduitsResponseConduit `json:"data"`
}

func (c *Manager) createConduit() error {
	ctx := context.TODO()

	conduit, err := c.ensureConduit(ctx)
	if err != nil {
		return fmt.Errorf("failed to ensure conduit: %w", err)
	}

	c.logger.Info(
		"conduit ensured",
		slog.String("id", conduit.Id),
		slog.Int("shard_count", conduit.ShardCount), // hardcoded, should be from conduit
	)

	if err := c.SubscribeWithLimits(
		ctx,
		eventsub.EventTypeUserAuthorizationRevoke,
		eventsub.ConduitTransport{
			Method:    "conduit",
			ConduitId: conduit.Id,
		},
		"1",
		"", // channelId, not needed for this event
		"", // botId, not needed for this event
	); err != nil {
		c.logger.Error("Failed to subscribe to UserAuthorizationRevoke event", logger.Error(err))
	} else {
		c.logger.Info("Subscribed to UserAuthorizationRevoke event")
	}

	c.currentConduit = conduit

	return nil
}

func (c *Manager) ensureConduit(ctx context.Context) (*conduitsResponseConduit, error) {
	mu := c.redSync.NewMutex("eventsub:conduits")
	err := mu.Lock()
	if err != nil {
		return nil, fmt.Errorf("failed to lock conduits mutex: %w", err)
	}
	defer mu.Unlock()

	appToken, err := c.twirBus.Tokens.RequestAppToken.Request(ctx, struct{}{})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.twitch.tv/helix/eventsub/conduits", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Client-id", c.config.TwitchClientId)
	req.Header.Set("Authorization", "Bearer "+appToken.Data.AccessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("cannot get conduits: %s", string(body))
	}

	var conduits conduitsResponse
	if err := json.Unmarshal(body, &conduits); err != nil {
		return nil, err
	}

	var currentConduit *conduitsResponseConduit
	if len(conduits.Data) == 0 {
		if err := c.conduitsRepository.DeleteAll(ctx); err != nil {
			return nil, fmt.Errorf("failed to delete all conduits: %w", err)
		}

		newConduit, err := c.twitchCreateConduit(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to create conduit: %w", err)
		}

		currentConduit = newConduit

		if _, err := c.conduitsRepository.Create(
			ctx, twitchconduits.CreateInput{
				ID:         currentConduit.Id,
				ShardCount: int8(currentConduit.ShardCount),
			},
		); err != nil {
			return nil, fmt.Errorf("failed to create conduit in db: %w", err)
		}
	} else {
		currentConduit = &conduits.Data[0]

		dbCurrentConduit, err := c.conduitsRepository.GetOne(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get current conduit from db: %w", err)
		}

		if dbCurrentConduit.ID != currentConduit.Id {
			if err := c.conduitsRepository.DeleteAll(ctx); err != nil {
				return nil, fmt.Errorf("failed to delete all conduits: %w", err)
			}

			if _, err := c.conduitsRepository.Create(
				ctx, twitchconduits.CreateInput{
					ID:         currentConduit.Id,
					ShardCount: int8(currentConduit.ShardCount),
				},
			); err != nil {
				return nil, fmt.Errorf("failed to create conduit in db: %w", err)
			}
		}
	}

	return currentConduit, nil
}

func (c *Manager) twitchCreateConduit(ctx context.Context) (*conduitsResponseConduit, error) {
	appToken, err := c.twirBus.Tokens.RequestAppToken.Request(ctx, struct{}{})
	if err != nil {
		return nil, err
	}

	createReq := createConduitRequest{
		ShardCount: 3, // how many replicas of eventsub i runed in prod
	}

	bodyBytes, err := json.Marshal(createReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal create conduit request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.twitch.tv/helix/eventsub/conduits", bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Client-id", c.config.TwitchClientId)
	req.Header.Set("Authorization", "Bearer "+appToken.Data.AccessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to create conduit: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("cannot create conduit: %s", string(body))
	}

	var createResp createConduitResponse
	if err := json.Unmarshal(body, &createResp); err != nil {
		return nil, err
	}

	if len(createResp.Data) == 0 {
		return nil, fmt.Errorf("no conduit created")
	}

	return &createResp.Data[0], nil
}

func (c *Manager) twitchUpdateConduitShard(ctx context.Context) error {
	mu := c.redSync.NewMutex("eventsub:shard:update")
	err := mu.Lock()
	if err != nil {
		return fmt.Errorf("failed to lock conduits mutex: %w", err)
	}
	defer mu.Unlock()

	if c.wsCurrentSessionId == nil {
		return fmt.Errorf("websocket session id is not set, cannot update conduit shard")
	}

	var shardId int
	currentReplicaId := os.Getenv("REPLICA")
	if currentReplicaId != "" {
		parsed, err := strconv.Atoi(currentReplicaId)
		if err != nil {
			return fmt.Errorf("failed to parse REPLICA env var: %w", err)
		}

		shardId = parsed - 1 // REPLICA is 1-based, but shardId is 0-based
	}

	appToken, err := c.twirBus.Tokens.RequestAppToken.Request(ctx, struct{}{})
	if err != nil {
		return err
	}

	updateReq := map[string]any{
		"conduit_id": c.currentConduit.Id,
		"shards": []map[string]any{
			{
				"id": shardId,
				"transport": map[string]any{
					"method":     "websocket",
					"session_id": c.wsCurrentSessionId,
				},
			},
		},
	}

	bodyBytes, err := json.Marshal(updateReq)
	if err != nil {
		return fmt.Errorf("failed to marshal update conduit shard request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, "https://api.twitch.tv/helix/eventsub/conduits/shards", bytes.NewBuffer(bodyBytes))
	if err != nil {
		return err
	}
	req.Header.Set("Client-id", c.config.TwitchClientId)
	req.Header.Set("Authorization", "Bearer "+appToken.Data.AccessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to update conduit shard: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("cannot update conduit shard: %s", string(body))
	}

	c.logger.Info(
		"Updated conduit shard",
		slog.String("conduit_id", c.currentConduit.Id),
		slog.Int("shard_id", shardId),
		slog.String("session_id", *c.wsCurrentSessionId),
		slog.String("current_replica_id", currentReplicaId),
		slog.String("response", string(body)),
	)

	return nil
}
