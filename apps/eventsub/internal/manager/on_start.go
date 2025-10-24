package manager

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strconv"

	"github.com/imroc/req/v3"
	twitchconduits "github.com/twirapp/twir/libs/repositories/twitch_conduits"
	"github.com/twirapp/twitchy/eventsub"
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
		c.logger.Error("Failed to subscribe to UserAuthorizationRevoke event", slog.Any("err", err))
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

	var conduits conduitsResponse
	resp, err := req.R().
		SetContext(ctx).
		SetHeader("Client-id", c.config.TwitchClientId).
		SetBearerAuthToken(appToken.Data.AccessToken).
		SetContentType("application/json").
		SetSuccessResult(&conduits).
		Get("https://api.twitch.tv/helix/eventsub/conduits")
	if err != nil {
		return nil, err
	}

	if !resp.IsSuccessState() {
		return nil, fmt.Errorf("cannot get conduits: %s", resp.String())
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
	var createResp *createConduitResponse

	resp, err := req.R().
		SetContext(ctx).
		SetHeader("Client-id", c.config.TwitchClientId).
		SetBearerAuthToken(appToken.Data.AccessToken).
		SetContentType("application/json").
		SetBody(createReq).
		SetSuccessResult(&createResp).
		Post("https://api.twitch.tv/helix/eventsub/conduits")
	if err != nil {
		return nil, fmt.Errorf("failed to create conduit: %w", err)
	}
	if !resp.IsSuccessState() {
		return nil, fmt.Errorf("cannot create conduit: %s", resp.String())
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

	resp, err := req.R().
		SetContext(ctx).
		SetHeader("Client-id", c.config.TwitchClientId).
		SetBearerAuthToken(appToken.Data.AccessToken).
		SetContentType("application/json").
		SetBody(updateReq).
		Patch("https://api.twitch.tv/helix/eventsub/conduits/shards")
	if err != nil {
		return fmt.Errorf("failed to update conduit shard: %w", err)
	}
	if !resp.IsSuccessState() {
		return fmt.Errorf("cannot update conduit shard: %s", resp.String())
	}

	c.logger.Info(
		"Updated conduit shard",
		slog.String("conduit_id", c.currentConduit.Id),
		slog.Int("shard_id", shardId),
		slog.String("session_id", *c.wsCurrentSessionId),
		slog.String("current_replica_id", currentReplicaId),
	)

	return nil
}
