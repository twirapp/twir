package twitch

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	json "github.com/bytedance/sonic"
	"github.com/nicklaw5/helix/v2"
	"github.com/samber/lo"
	"github.com/satont/twir/apps/api/internal/impl_deps"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/twitch"
	generatedTwitch "github.com/twirapp/twir/libs/api/messages/twitch"
)

type Twitch struct {
	*impl_deps.Deps
}

const redisLoginsPrefix = "twitch:user:by:login:"
const redisIdsPrefix = "twitch:user:by:id:"
const cacheDuration = 24 * time.Hour

func (c *Twitch) getUsersFromCache(ctx context.Context, keys []string) ([]helix.User, error) {
	if len(keys) == 0 {
		return nil, nil
	}

	cachedUsers, err := c.Redis.MGet(ctx, keys...).Result()
	if err != nil {
		return nil, err
	}

	var users []helix.User
	for _, cachedUser := range cachedUsers {
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

func (c *Twitch) getUsersFromTwitch(ctx context.Context, params *helix.UsersParams) (
	[]helix.User,
	error,
) {
	twitchClient, err := twitch.NewAppClientWithContext(ctx, c.Config, c.Grpc.Tokens)
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

func (c *Twitch) TwitchGetUsers(
	ctx context.Context,
	req *generatedTwitch.TwitchGetUsersRequest,
) (*generatedTwitch.TwitchGetUsersResponse, error) {
	var wg sync.WaitGroup
	var mu sync.Mutex
	twitchUsers := make([]helix.User, 0, len(req.Ids)+len(req.Names))

	req.Ids = lo.Filter(
		req.Ids, func(id string, _ int) bool {
			return id != ""
		},
	)
	req.Names = lo.Filter(
		req.Names, func(name string, _ int) bool {
			return name != ""
		},
	)

	if len(req.Ids) == 0 && len(req.Names) == 0 {
		return &generatedTwitch.TwitchGetUsersResponse{
			Users: nil,
		}, nil
	}

	cachedUsersByIds, err := c.getUsersFromCache(
		ctx, lo.Map(
			req.Ids, func(id string, _ int) string {
				return redisIdsPrefix + id
			},
		),
	)
	if err != nil {
		return nil, err
	}
	req.Ids = lo.Filter(
		req.Ids, func(id string, _ int) bool {
			return !lo.ContainsBy(
				cachedUsersByIds, func(user helix.User) bool {
					return user.ID == id
				},
			)
		},
	)
	twitchUsers = append(twitchUsers, cachedUsersByIds...)

	cachedUsersByNames, err := c.getUsersFromCache(
		ctx, lo.Map(
			req.Names, func(name string, _ int) string {
				return redisLoginsPrefix + strings.ToLower(name)
			},
		),
	)
	if err != nil {
		return nil, err
	}
	req.Names = lo.Filter(
		req.Names, func(name string, _ int) bool {
			return !lo.ContainsBy(
				cachedUsersByNames, func(user helix.User) bool {
					return user.Login == strings.ToLower(name)
				},
			)
		},
	)
	twitchUsers = append(twitchUsers, cachedUsersByNames...)

	idsChunks := lo.Chunk(req.Ids, 100)
	namesChunks := lo.Chunk(req.Names, 100)

	for _, idsChunk := range idsChunks {
		wg.Add(1)
		go func(ids []string) {
			defer wg.Done()
			users, err := c.getUsersFromTwitch(
				ctx, &helix.UsersParams{
					IDs: ids,
				},
			)
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
			users, err := c.getUsersFromTwitch(
				ctx, &helix.UsersParams{
					Logins: ids,
				},
			)
			if err != nil {
				return
			}
			mu.Lock()
			defer mu.Unlock()
			twitchUsers = append(twitchUsers, users...)
		}(namesChunk)
	}

	wg.Wait()

	twitchUsers = lo.UniqBy(
		twitchUsers, func(user helix.User) string {
			return user.ID
		},
	)

	convertedUsers := lo.Map(
		twitchUsers, func(user helix.User, _ int) *generatedTwitch.TwitchUser {
			return &generatedTwitch.TwitchUser{
				Id:              user.ID,
				Login:           user.Login,
				DisplayName:     user.DisplayName,
				Type:            user.Type,
				BroadcasterType: user.BroadcasterType,
				Description:     user.Description,
				ProfileImageUrl: user.ProfileImageURL,
				OfflineImageUrl: user.OfflineImageURL,
				CreatedAt:       fmt.Sprint(user.CreatedAt.UnixMilli()),
			}
		},
	)

	return &generatedTwitch.TwitchGetUsersResponse{
		Users: convertedUsers,
	}, nil
}

func (c *Twitch) TwitchSearchChannels(
	ctx context.Context,
	req *generatedTwitch.TwitchSearchChannelsRequest,
) (*generatedTwitch.TwitchSearchChannelsResponse, error) {
	if req.Query == "" {
		return nil, fmt.Errorf("query is empty")
	}

	twitchClient, err := twitch.NewAppClientWithContext(ctx, c.Config, c.Grpc.Tokens)
	if err != nil {
		return nil, err
	}

	twitchReq, twitchErr := twitchClient.SearchChannels(
		&helix.SearchChannelsParams{
			Channel: req.Query,
		},
	)
	if twitchErr != nil {
		return nil, twitchErr
	}
	if twitchReq.ErrorMessage != "" {
		return nil, fmt.Errorf(twitchReq.ErrorMessage)
	}

	channels := make([]*generatedTwitch.Channel, 0, len(twitchReq.Data.Channels))
	for _, channel := range twitchReq.Data.Channels {
		channels = append(
			channels, &generatedTwitch.Channel{
				Id:              channel.ID,
				Login:           channel.BroadcasterLogin,
				DisplayName:     channel.DisplayName,
				ProfileImageUrl: channel.ThumbnailURL,
				Title:           channel.Title,
				GameName:        channel.GameName,
				GameId:          channel.GameID,
				IsLive:          channel.IsLive,
			},
		)
	}

	sort.Slice(
		channels, func(i, j int) bool {
			name1 := channels[i].Login
			name2 := channels[j].Login

			containsName1 := strings.Contains(name1, req.Query)
			containsName2 := strings.Contains(name2, req.Query)

			if containsName1 && !containsName2 {
				return true
			} else if !containsName1 && containsName2 {
				return false
			} else {
				return name1 < name2
			}
		},
	)

	if req.TwirOnly {
		var dbUsers []model.Users
		channelsIds := make([]string, 0, len(channels))
		for _, channel := range channels {
			channelsIds = append(channelsIds, channel.Id)
		}

		if err := c.Db.
			WithContext(ctx).
			Select("id").
			Where("id IN ?", channelsIds).
			Find(&dbUsers).Error; err != nil {
			return nil, err
		}

		channels = lo.Filter(
			channels,
			func(channel *generatedTwitch.Channel, _ int) bool {
				return lo.ContainsBy(
					dbUsers, func(user model.Users) bool {
						return user.ID == channel.Id
					},
				)
			},
		)
	}

	return &generatedTwitch.TwitchSearchChannelsResponse{
		Channels: channels,
	}, nil
}
