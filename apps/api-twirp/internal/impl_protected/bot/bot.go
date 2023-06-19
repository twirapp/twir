package bot

import (
	"context"
	"github.com/nicklaw5/helix/v2"
	"github.com/samber/lo"
	"github.com/satont/tsuwari/apps/api-twirp/internal/impl_deps"
	model "github.com/satont/tsuwari/libs/gomodels"
	"github.com/satont/tsuwari/libs/grpc/generated/api/bots"
	"github.com/satont/tsuwari/libs/grpc/generated/api/meta"
	"github.com/satont/tsuwari/libs/twitch"
	"github.com/twitchtv/twirp"
	"google.golang.org/protobuf/types/known/emptypb"
	"net/http"
	"sync"
)

type Bot struct {
	*impl_deps.Deps
}

func (c *Bot) BotInfo(ctx context.Context, meta *meta.BaseRequestMeta) (*bots.BotInfo, error) {
	dashboardId, ok := ctx.Value("dashboardId").(string)
	if !ok || dashboardId == "" {
		return nil, twirp.NewError(twirp.ErrorCode(http.StatusBadRequest), "no dashboardId provided")
	}

	dbUser := &model.Users{}
	err := c.Db.WithContext(ctx).Where("id = ?", dashboardId).Preload("Channel").Find(dbUser).Error
	if err != nil {
		return nil, err
	}
	if dbUser.ID == "" || dbUser.Channel == nil {
		return nil, twirp.NotFoundError("user not found")
	}

	twitchClient, err := twitch.NewUserClient(dashboardId, *c.Config, c.Grpc.Tokens)
	if err != nil {
		return nil, err
	}

	var wg sync.WaitGroup
	result := &bots.BotInfo{}
	wg.Add(2)

	go func() {
		defer wg.Done()

		if dashboardId == dbUser.Channel.BotID {
			result.IsMod = true
			return
		}

		mods, err := twitchClient.GetModerators(&helix.GetModeratorsParams{
			BroadcasterID: dashboardId,
			UserIDs:       []string{dbUser.Channel.BotID},
		})
		if err != nil {
			return
		}

		result.IsMod = lo.If(len(mods.Data.Moderators) == 0, false).Else(true)
	}()

	go func() {
		defer wg.Done()
		infoReq, err := twitchClient.GetUsers(&helix.UsersParams{
			IDs: []string{dbUser.Channel.BotID},
		})
		if err != nil {
			return
		}

		if len(infoReq.Data.Users) == 0 {
			return
		}

		result.BotId = infoReq.Data.Users[0].ID
		result.BotName = infoReq.Data.Users[0].Login
	}()

	wg.Wait()

	return result, nil
}

func (c *Bot) BotJoinPart(ctx context.Context, request *bots.BotJoinPartRequest) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
}
