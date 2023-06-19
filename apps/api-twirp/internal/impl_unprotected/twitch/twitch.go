package twitch

import (
	"context"
	"github.com/nicklaw5/helix/v2"
	"github.com/samber/lo"
	"github.com/satont/tsuwari/apps/api-twirp/internal/impl_deps"
	generatedTwitch "github.com/satont/tsuwari/libs/grpc/generated/api/twitch"
	"github.com/satont/tsuwari/libs/twitch"
	"github.com/twitchtv/twirp"
	"sync"
)

type Twitch struct {
	*impl_deps.Deps
}

func (c *Twitch) TwitchSearchUsers(ctx context.Context, req *generatedTwitch.TwitchSearchUsersRequest) (*generatedTwitch.TwitchSearchUsersResponse, error) {
	twitchClient, err := twitch.NewAppClientWithContext(ctx, *c.Config, c.Grpc.Tokens)
	if err != nil {
		return nil, twirp.Internal.Error(err.Error())
	}

	var wg sync.WaitGroup
	var mu sync.Mutex
	twitchUsers := make([]helix.User, 0, len(req.Ids)+len(req.Names))

	idsChunks := lo.Chunk(req.Ids, 100)
	namesChunks := lo.Chunk(req.Names, 100)
	wg.Add(len(idsChunks) + len(namesChunks))

	for _, idsChunk := range idsChunks {
		go func(ids []string) {
			defer wg.Done()
			twitchReq, twitchErr := twitchClient.GetUsers(&helix.UsersParams{
				IDs: ids,
			})
			if twitchErr != nil || twitchReq.ErrorMessage != "" || len(twitchReq.Data.Users) == 0 {
				return
			}
			mu.Lock()
			defer mu.Unlock()
			twitchUsers = append(twitchUsers, twitchReq.Data.Users...)
		}(idsChunk)
	}

	for _, namesChunk := range namesChunks {
		go func(names []string) {
			defer wg.Done()
			twitchReq, twitchErr := twitchClient.GetUsers(&helix.UsersParams{
				Logins: names,
			})
			if twitchErr != nil || twitchReq.ErrorMessage != "" || len(twitchReq.Data.Users) == 0 {
				return
			}
			mu.Lock()
			defer mu.Unlock()
			twitchUsers = append(twitchUsers, twitchReq.Data.Users...)
		}(namesChunk)
	}

	wg.Wait()

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
			ViewCount:       int64(user.ViewCount),
			Email:           user.Email,
			CreatedAt:       uint64(user.CreatedAt.UnixMilli()),
		}
	})

	return &generatedTwitch.TwitchSearchUsersResponse{
		Users: convertedUsers,
	}, nil
}
