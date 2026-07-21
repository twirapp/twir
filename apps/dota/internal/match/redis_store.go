package match

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

const (
	lifecycleActionStreamKey    = "stream:twir:dota:{dota}:lifecycle-actions"
	lifecycleActionStreamMaxLen = int64(100_000)
	updateStatsMaxRetries       = 3

	absentSnapshotSentinel = ""

	// Append intents before state so failed XADDs cannot advance state; retries may
	// duplicate intents, which workers safely reject through snapshot validation.
	compareAndSwapSnapshotScript = `
local current = redis.call("GET", KEYS[1])
if current == false then current = "" end
if current == ARGV[2] then return 2 end
if current ~= ARGV[1] then return 0 end
for i = 5, #ARGV do
  redis.call("XADD", KEYS[2], "MAXLEN", "~", ARGV[4], "*", "action", ARGV[i])
end
redis.call("SET", KEYS[1], ARGV[2], "PX", ARGV[3])
return 1`
)

type redisStateClient interface {
	Get(context.Context, string) *redis.StringCmd
	Eval(context.Context, string, []string, ...interface{}) *redis.Cmd
}

type RedisStateStore struct {
	client redisStateClient
}

var _ StateStore = (*RedisStateStore)(nil)

func NewRedisStateStore(client *redis.Client) *RedisStateStore {
	return &RedisStateStore{client: client}
}

func (s *RedisStateStore) Load(ctx context.Context, channelID uuid.UUID) (Snapshot, error) {
	raw, err := s.client.Get(ctx, snapshotKey(channelID)).Result()
	if errors.Is(err, redis.Nil) {
		return Snapshot{
			ChannelID: channelID,
			State:     StateIdle,
		}, nil
	}
	if err != nil {
		return Snapshot{}, fmt.Errorf("load dota match snapshot: %w", err)
	}

	var snapshot Snapshot
	if err := json.Unmarshal([]byte(raw), &snapshot); err != nil {
		return Snapshot{}, fmt.Errorf("decode dota match snapshot: %w", err)
	}
	if snapshot.ChannelID == uuid.Nil {
		snapshot.ChannelID = channelID
	} else if snapshot.ChannelID != channelID {
		return Snapshot{}, fmt.Errorf(
			"stored snapshot channel ID %s does not match requested channel %s",
			snapshot.ChannelID,
			channelID,
		)
	}

	snapshot.rawState = raw
	return snapshot, nil
}

func (s *RedisStateStore) CompareAndSwap(
	ctx context.Context,
	current Snapshot,
	next Snapshot,
	actions []LifecycleAction,
) (bool, error) {
	if err := validateCompareAndSwapInput(current, next, actions); err != nil {
		return false, err
	}
	mutationID, err := uuid.NewRandom()
	if err != nil {
		return false, fmt.Errorf("generate dota match mutation ID: %w", err)
	}
	next.MutationID = mutationID.String()

	expectedRaw, err := expectedSnapshotRaw(current)
	if err != nil {
		return false, fmt.Errorf("encode expected dota match snapshot: %w", err)
	}
	nextRaw, err := json.Marshal(next)
	if err != nil {
		return false, fmt.Errorf("encode next dota match snapshot: %w", err)
	}

	args := make([]interface{}, 0, 4+len(actions))
	args = append(
		args,
		expectedRaw,
		string(nextRaw),
		snapshotTTL.Milliseconds(),
		lifecycleActionStreamMaxLen,
	)
	for _, action := range actions {
		actionJSON, err := json.Marshal(action)
		if err != nil {
			return false, fmt.Errorf("encode lifecycle action: %w", err)
		}
		args = append(args, string(actionJSON))
	}

	updated, err := s.client.Eval(
		ctx,
		compareAndSwapSnapshotScript,
		[]string{snapshotKey(current.ChannelID), lifecycleActionStreamKey},
		args...,
	).Int64()
	if err != nil {
		return false, fmt.Errorf("compare and swap dota match snapshot: %w", err)
	}

	switch updated {
	case 0:
		return false, nil
	case 1, 2:
		return true, nil
	default:
		return false, fmt.Errorf("unexpected dota match compare-and-swap result: %d", updated)
	}
}

func (s *RedisStateStore) UpdateStats(
	ctx context.Context,
	channelID uuid.UUID,
	mmr int,
	sessionWins int,
	sessionLosses int,
) error {
	for attempt := 0; attempt < updateStatsMaxRetries; attempt++ {
		current, err := s.Load(ctx, channelID)
		if err != nil {
			return err
		}

		next := current
		next.Revision++
		next.Mmr = mmr
		next.SessionWins = sessionWins
		next.SessionLosses = sessionLosses

		updated, err := s.CompareAndSwap(ctx, current, next, nil)
		if err != nil {
			return err
		}
		if updated {
			return nil
		}
	}

	return fmt.Errorf("update dota match stats: compare-and-swap retry limit reached")
}

func validateCompareAndSwapInput(current Snapshot, next Snapshot, actions []LifecycleAction) error {
	if current.ChannelID == uuid.Nil {
		return errors.New("current snapshot channel ID is required")
	}
	if next.ChannelID != current.ChannelID {
		return errors.New("next snapshot channel ID must match current snapshot")
	}
	if current.Revision == math.MaxUint64 {
		return errors.New("current snapshot revision cannot be incremented")
	}
	if next.Revision != current.Revision+1 {
		return fmt.Errorf(
			"next snapshot revision must increment from %d to %d",
			current.Revision,
			current.Revision+1,
		)
	}

	for _, action := range actions {
		switch action.Kind {
		case ActionCreate, ActionResolve, ActionCancel:
		default:
			return errors.New("lifecycle action kind is invalid")
		}
		if action.ChannelID == uuid.Nil {
			return errors.New("lifecycle action channel ID is required")
		}
		if action.ChannelID != current.ChannelID {
			return errors.New("lifecycle action channel ID must match snapshot")
		}
		if action.MatchID <= 0 {
			return errors.New("lifecycle action match ID must be positive")
		}
		if action.Revision != next.Revision {
			return errors.New("lifecycle action revision must match next snapshot")
		}
		switch action.Kind {
		case ActionCreate:
			if next.MatchID <= 0 {
				return errors.New("next snapshot match ID must be positive for a create action")
			}
			if action.MatchID != next.MatchID {
				return errors.New("create action match ID must match next snapshot")
			}
		case ActionResolve, ActionCancel:
			if current.MatchID <= 0 {
				return errors.New("current snapshot match ID must be positive for a terminal action")
			}
			if action.MatchID != current.MatchID {
				return errors.New("terminal action match ID must match current snapshot")
			}
		}
	}

	return nil
}

func expectedSnapshotRaw(snapshot Snapshot) (string, error) {
	if snapshot.rawState != "" {
		return snapshot.rawState, nil
	}
	if snapshot.Revision == 0 {
		return absentSnapshotSentinel, nil
	}

	raw, err := json.Marshal(snapshot)
	if err != nil {
		return "", err
	}
	return string(raw), nil
}
