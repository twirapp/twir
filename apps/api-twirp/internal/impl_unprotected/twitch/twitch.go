package twitch

import (
	"context"
	"encoding/json"
	"github.com/nicklaw5/helix/v2"
	"github.com/samber/lo"
	"github.com/satont/tsuwari/apps/api-twirp/internal/impl_deps"
	generatedTwitch "github.com/satont/tsuwari/libs/grpc/generated/api/twitch"
	"github.com/satont/tsuwari/libs/twitch"
	"strings"
	"sync"
	"time"
)

type Twitch struct {
	*impl_deps.Deps
}

const redisLoginsPrefix = "api:cache:twitch:users:by:logins:"
const redisIdsPrefix = "api:cache:twitch:users:by:ids:"
const cacheDuration = 24 * time.Hour

func (c *Twitch) getUsersFromCache(ctx context.Context, keys []string) ([]helix.User, error) {
	if len(keys) == 0 {
		return nil, nil
	}

	cachedUsersByIds, err := c.Redis.MGet(ctx, keys...).Result()
	if err != nil {
		return nil, err
	}

	var users []helix.User
	for _, cachedUser := range cachedUsersByIds {
		if cachedUser == nil {
			continue
		}

		var user helix.User
		if err := json.Unmarshal([]byte(cachedUser.(string)), &user); err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

func (c *Twitch) getUsersFromTwitch(ctx context.Context, params *helix.UsersParams) ([]helix.User, error) {
	twitchClient, err := twitch.NewAppClientWithContext(ctx, *c.Config, c.Grpc.Tokens)
	if err != nil {
		return nil, err
	}

	twitchReq, twitchErr := twitchClient.GetUsers(params)
	if twitchErr != nil || twitchReq.ErrorMessage != "" || len(twitchReq.Data.Users) == 0 {
		return nil, twitchErr
	}

	defer func() {
		for _, user := range twitchReq.Data.Users {
			bytes, err := json.Marshal(user)
			if err == nil {
				c.Redis.Set(ctx, redisLoginsPrefix+strings.ToLower(user.Login), bytes, cacheDuration)
				c.Redis.Set(ctx, redisIdsPrefix+user.ID, bytes, cacheDuration)
			}
		}
	}()

	return twitchReq.Data.Users, nil
}

func (c *Twitch) TwitchSearchUsers(ctx context.Context, req *generatedTwitch.TwitchSearchUsersRequest) (*generatedTwitch.TwitchSearchUsersResponse, error) {
	var wg sync.WaitGroup
	var mu sync.Mutex
	twitchUsers := make([]helix.User, 0, len(req.Ids)+len(req.Names))

	cachedUsersByIds, err := c.getUsersFromCache(ctx, lo.Map(req.Ids, func(id string, _ int) string {
		return redisIdsPrefix + id
	}))
	if err != nil {
		return nil, err
	}
	req.Ids = lo.Filter(req.Ids, func(id string, _ int) bool {
		return !lo.ContainsBy(cachedUsersByIds, func(user helix.User) bool {
			return user.ID == id
		})
	})
	twitchUsers = append(twitchUsers, cachedUsersByIds...)

	cachedUsersByNames, err := c.getUsersFromCache(ctx, lo.Map(req.Names, func(name string, _ int) string {
		return redisLoginsPrefix + name
	}))
	if err != nil {
		return nil, err
	}
	req.Names = lo.Filter(req.Names, func(name string, _ int) bool {
		return !lo.ContainsBy(cachedUsersByNames, func(user helix.User) bool {
			return user.Login == strings.ToLower(name)
		})
	})
	twitchUsers = append(twitchUsers, cachedUsersByNames...)

	idsChunks := lo.Chunk(req.Ids, 100)
	namesChunks := lo.Chunk(req.Names, 100)

	for _, idsChunk := range idsChunks {
		wg.Add(1)
		go func(ids []string) {
			defer wg.Done()
			users, err := c.getUsersFromTwitch(ctx, &helix.UsersParams{
				IDs: ids,
			})
			if err != nil {
				return
			}
			mu.Lock()
			defer mu.Unlock()
			twitchUsers = append(twitchUsers, users...)
		}(idsChunk)
	}

	for _, namesChunk := range namesChunks {
		wg.Add(1)
		go func(ids []string) {
			defer wg.Done()
			users, err := c.getUsersFromTwitch(ctx, &helix.UsersParams{
				Logins: ids,
			})
			if err != nil {
				return
			}
			mu.Lock()
			defer mu.Unlock()
			twitchUsers = append(twitchUsers, users...)
		}(namesChunk)
	}

	wg.Wait()

	twitchUsers = lo.UniqBy(twitchUsers, func(user helix.User) string {
		return user.ID
	})

	convertedUsers := lo.Map(twitchUsers, func(user helix.User, _ int) *generatedTwitch.TwitchUser {
		return &generatedTwitch.TwitchUser{
			Id:              user.ID,
			Login:           user.Login,
			DisplayName:     user.DisplayName,
			Type:            user.Type,
			BroadcasterType: user.BroadcasterType,
			Description:     user.Description,
			ProfileImageUrl: user.ProfileImageURL,
			OfflineImageUrl: user.OfflineImageURL,
			CreatedAt:       uint64(user.CreatedAt.UnixMilli()),
		}
	})

	return &generatedTwitch.TwitchSearchUsersResponse{
		Users: convertedUsers,
	}, nil
}
