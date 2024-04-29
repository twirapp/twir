package twitch

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/goccy/go-json"
	"github.com/nicklaw5/helix/v2"
	"github.com/samber/lo"
	"golang.org/x/sync/errgroup"
)

const userIdCacheKey = "cache:twir:twitch:users-by-id:"
const userNameCacheKey = "cache:twir:twitch:users-by-name:"
const userCacheDuration = 1 * time.Hour

func buildUserCacheKeyForId(userId string) string {
	return userIdCacheKey + userId
}

func buildUserCacheKeyForName(userName string) string {
	return userNameCacheKey + userName
}

type TwitchUser struct {
	helix.User
	// internal field, do not use outside
	NotFound bool `json:"isTwitchBanned"`
}

func (c *CachedTwitchClient) GetUserById(ctx context.Context, id string) (*TwitchUser, error) {
	if id == "" {
		return nil, nil
	}

	if bytes, _ := c.redis.Get(ctx, buildUserCacheKeyForId(id)).Bytes(); len(bytes) > 0 {
		var helixUser TwitchUser
		if err := json.Unmarshal(bytes, &helixUser); err != nil {
			return nil, err
		}

		return &helixUser, nil
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

	user := TwitchUser{
		User:     twitchReq.Data.Users[0],
		NotFound: false,
	}

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

	return &user, nil
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
							User:     user,
							NotFound: false,
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
			user.NotFound = true
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
			return !item.NotFound
		},
	)

	return resultedUsers, nil
}

func (c *CachedTwitchClient) GetUsersByNames(ctx context.Context, names []string) (
	[]TwitchUser,
	error,
) {
	if len(names) == 0 {
		return nil, nil
	}

	for i, name := range names {
		names[i] = strings.ToLower(name)
	}

	var resultedUsers []TwitchUser
	var resultedUsersMutex sync.Mutex

	var neededNamesForRequest []string

	for _, name := range names {
		bytes, err := c.redis.Get(ctx, buildUserCacheKeyForName(name)).Bytes()
		if len(bytes) == 0 || err != nil {
			neededNamesForRequest = append(neededNamesForRequest, name)
			continue
		}

		var user TwitchUser
		if err := json.Unmarshal(bytes, &user); err != nil {
			return nil, err
		}

		resultedUsers = append(resultedUsers, user)
	}

	var twitchWg errgroup.Group
	for _, chunk := range lo.Chunk(neededNamesForRequest, 100) {
		chunk := chunk

		twitchWg.Go(
			func() error {
				twitchReq, err := c.client.GetUsers(&helix.UsersParams{Logins: chunk})
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
							User:     user,
							NotFound: false,
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

	for _, name := range names {
		user, _ := lo.Find(
			resultedUsers, func(item TwitchUser) bool {
				return strings.ToLower(item.Login) == name
			},
		)
		if user.Login == "" {
			user.ID = name
			user.NotFound = true
		}

		userBytes, err := json.Marshal(&user)
		if err != nil {
			return nil, err
		}

		if err := c.redis.Set(
			ctx,
			buildUserCacheKeyForName(user.ID),
			userBytes,
			userCacheDuration,
		).Err(); err != nil {
			return nil, err
		}
	}

	resultedUsers = lo.Filter(
		resultedUsers,
		func(item TwitchUser, _ int) bool {
			return !item.NotFound
		},
	)

	return resultedUsers, nil
}
