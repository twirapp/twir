package twitch

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/goccy/go-json"
	"github.com/nicklaw5/helix/v2"
	"github.com/samber/lo"
	"golang.org/x/sync/errgroup"
)

const userIdCacheKey = "cache:twir:twitch:users:"
const userCacheDuration = 1 * time.Hour

func buildUserCacheKeyForId(userId string) string {
	return userIdCacheKey + userId
}

type TwitchUser struct {
	helix.User
	IsTwitchBanned bool `json:"isTwitchBanned"`
}

func (c *CachedTwitchClient) GetUserById(ctx context.Context, id string) (*TwitchUser, error) {
	if id == "" {
		return nil, nil
	}

	if bytes, _ := c.redis.Get(ctx, buildUserCacheKeyForId(id)).Bytes(); len(bytes) > 0 {
		var helixUser *TwitchUser
		if err := json.Unmarshal(bytes, helixUser); err != nil {
			return nil, err
		}

		return helixUser, nil
	}

	twitchReq, err := c.client.GetUsers(&helix.UsersParams{IDs: []string{id}})
	if err != nil {
		return nil, err
	}
	if twitchReq.ErrorMessage != "" {
		return nil, fmt.Errorf("cannot get twitch user: %s", twitchReq.ErrorMessage)
	}

	if len(twitchReq.Data.Users) == 0 {
		return nil, fmt.Errorf("user not found")
	}

	user := twitchReq.Data.Users[0]

	userBytes, err := json.Marshal(&user)
	if err != nil {
		return nil, err
	}

	if err := c.redis.Set(
		ctx,
		buildUserCacheKeyForId(id),
		userBytes,
		userCacheDuration,
	).Err(); err != nil {
		return nil, err
	}

	return &TwitchUser{
		User:           user,
		IsTwitchBanned: false,
	}, nil
}

func (c *CachedTwitchClient) GetUsersByIds(ctx context.Context, ids []string) (
	[]TwitchUser,
	error,
) {
	if len(ids) == 0 {
		return nil, nil
	}

	var resultedUsers []TwitchUser
	var resultedUsersMutex sync.Mutex

	var neededIdsForRequest []string

	for _, id := range ids {
		bytes, err := c.redis.Get(ctx, buildUserCacheKeyForId(id)).Bytes()
		if len(bytes) == 0 || err != nil {
			neededIdsForRequest = append(neededIdsForRequest, id)
			continue
		}

		var user TwitchUser
		if err := json.Unmarshal(bytes, &user); err != nil {
			return nil, err
		}

		resultedUsers = append(resultedUsers, user)
	}

	var twitchWg errgroup.Group
	for _, chunk := range lo.Chunk(neededIdsForRequest, 100) {
		chunk := chunk

		twitchWg.Go(
			func() error {
				twitchReq, err := c.client.GetUsers(&helix.UsersParams{IDs: chunk})
				if err != nil {
					return err
				}
				if twitchReq.ErrorMessage != "" {
					return fmt.Errorf("cannot get twitch user: %s", twitchReq.ErrorMessage)
				}

				resultedUsersMutex.Lock()

				for _, user := range twitchReq.Data.Users {
					resultedUsers = append(
						resultedUsers, TwitchUser{
							User:           user,
							IsTwitchBanned: false,
						},
					)
				}

				resultedUsersMutex.Unlock()

				return nil
			},
		)
	}

	if err := twitchWg.Wait(); err != nil {
		return nil, err
	}

	for _, id := range ids {
		user, _ := lo.Find(
			resultedUsers, func(item TwitchUser) bool {
				return item.ID == id
			},
		)
		if user.Login == "" {
			user.ID = id
			user.IsTwitchBanned = true
		}

		userBytes, err := json.Marshal(&user)
		if err != nil {
			return nil, err
		}

		if err := c.redis.Set(
			ctx,
			buildUserCacheKeyForId(user.ID),
			userBytes,
			userCacheDuration,
		).Err(); err != nil {
			return nil, err
		}
	}

	resultedUsers = lo.Filter(
		resultedUsers,
		func(item TwitchUser, _ int) bool {
			return !item.IsTwitchBanned
		},
	)

	return resultedUsers, nil
}
