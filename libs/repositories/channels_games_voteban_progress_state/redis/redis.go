package redis

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	channelsgamesvotebanprogressstate "github.com/twirapp/twir/libs/repositories/channels_games_voteban_progress_state"
	"github.com/twirapp/twir/libs/repositories/channels_games_voteban_progress_state/model"
)

type Opts struct {
	Redis *redis.Client
}

func New(opts Opts) *Redis {
	return &Redis{
		redis: opts.Redis,
	}
}

func NewFx(redis *redis.Client) *Redis {
	return New(Opts{Redis: redis})
}

var _ channelsgamesvotebanprogressstate.Repository = (*Redis)(nil)

type Redis struct {
	redis *redis.Client
}

func (r *Redis) buildKey(channelID string) string {
	return fmt.Sprintf("channels:%s:games:voteban", channelID)
}

func (r *Redis) buildUserVoteKey(channelID, userID string) string {
	return fmt.Sprintf("channels:%s:games:voteban:votes:%s", channelID, userID)
}

func (r *Redis) buildUserVotesPattern(channelID string) string {
	return fmt.Sprintf("channels:%s:games:voteban:votes:*", channelID)
}

func (r *Redis) Get(ctx context.Context, channelID string) (model.VoteState, error) {
	key := r.buildKey(channelID)

	// Check if the hash exists first
	exists, err := r.redis.Exists(ctx, key).Result()
	if err != nil {
		return model.Nil, err
	}
	if exists == 0 {
		return model.Nil, channelsgamesvotebanprogressstate.ErrNotFound
	}

	var state model.VoteState
	if err := r.redis.HGetAll(ctx, key).Scan(&state); err != nil {
		if errors.Is(err, redis.Nil) {
			return model.Nil, channelsgamesvotebanprogressstate.ErrNotFound
		}
		return model.Nil, err
	}

	return state, nil
}

func (r *Redis) Exists(ctx context.Context, channelID string) (bool, error) {
	key := r.buildKey(channelID)

	exists, err := r.redis.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}

	return exists > 0, nil
}

func (r *Redis) Create(ctx context.Context, channelID string, state model.VoteState, ttl time.Duration) error {
	key := r.buildKey(channelID)

	if err := r.redis.HSet(ctx, key, state).Err(); err != nil {
		return err
	}

	if err := r.redis.Expire(ctx, key, ttl).Err(); err != nil {
		// Cleanup on error
		r.redis.Del(ctx, key)
		return err
	}

	return nil
}

func (r *Redis) Update(ctx context.Context, channelID string, state model.VoteState) error {
	key := r.buildKey(channelID)
	return r.redis.HSet(ctx, key, state).Err()
}

func (r *Redis) Delete(ctx context.Context, channelID string) error {
	key := r.buildKey(channelID)
	return r.redis.Del(ctx, key).Err()
}

func (r *Redis) UserHasVoted(ctx context.Context, channelID, userID string) (bool, error) {
	key := r.buildUserVoteKey(channelID, userID)

	exists, err := r.redis.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}

	return exists > 0, nil
}

func (r *Redis) MarkUserVoted(ctx context.Context, channelID, userID string, ttl time.Duration) error {
	key := r.buildUserVoteKey(channelID, userID)
	return r.redis.Set(ctx, key, 1, ttl).Err()
}

func (r *Redis) ClearUserVotes(ctx context.Context, channelID string) error {
	pattern := r.buildUserVotesPattern(channelID)

	var cursor uint64
	for {
		keys, nextCursor, err := r.redis.Scan(ctx, cursor, pattern, 100).Result()
		if err != nil {
			return err
		}

		if len(keys) > 0 {
			if err := r.redis.Del(ctx, keys...).Err(); err != nil {
				return err
			}
		}

		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}

	return nil
}
