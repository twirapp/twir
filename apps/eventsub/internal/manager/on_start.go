package manager

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strconv"

	"github.com/kvizyx/twitchy/eventsub"
	"github.com/kvizyx/twitchy/helix"
	"github.com/twirapp/twir/libs/logger"
	twitchconduits "github.com/twirapp/twir/libs/repositories/twitch_conduits"
)

type conduitsResponseConduit struct {
	Id         string `json:"id"`
	ShardCount int    `json:"shard_count"`
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
		"", // broadcasterId, not needed for this event
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

	twitchClient, err := c.newAppTwitchClient(ctx)
	if err != nil {
		return nil, err
	}

	conduits, err := twitchClient.Conduits.GetConduits(ctx, helix.GetConduitsRequest{})
	if err != nil {
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
		currentConduit = &conduitsResponseConduit{
			Id:         conduits.Data[0].ID,
			ShardCount: conduits.Data[0].ShardCount,
		}

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
	twitchClient, err := c.newAppTwitchClient(ctx)
	if err != nil {
		return nil, err
	}

	shardCount := 3 // how many replicas of eventsub i runed in prod
	if c.config.TwitchMockEnabled {
		shardCount = 1
	}

	createResp, err := twitchClient.Conduits.CreateConduits(ctx, helix.CreateConduitsRequest{ShardCount: shardCount})
	if err != nil {
		return nil, fmt.Errorf("create conduit: %w", err)
	}
	if len(createResp.Data) == 0 {
		return nil, fmt.Errorf("no conduit created")
	}

	return &conduitsResponseConduit{
		Id:         createResp.Data[0].ID,
		ShardCount: createResp.Data[0].ShardCount,
	}, nil
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

	twitchClient, err := c.newAppTwitchClient(ctx)
	if err != nil {
		return err
	}

	_, err = twitchClient.Conduits.UpdateConduitShards(ctx, helix.UpdateConduitShardsRequest{
		ConduitID: c.currentConduit.Id,
		Shards: []helix.UpdateConduitShard{{
			ID: strconv.Itoa(shardId),
			Transport: helix.ConduitShardTransport{
				Method:    "websocket",
				SessionID: *c.wsCurrentSessionId,
			},
		}},
	})
	if err != nil {
		return fmt.Errorf("update conduit shard: %w", err)
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
