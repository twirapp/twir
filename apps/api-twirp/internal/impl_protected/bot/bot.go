package bot

import (
	"context"
	"net/http"
	"sync"

	"github.com/nicklaw5/helix/v2"
	"github.com/samber/lo"
	"github.com/satont/twir/apps/api-twirp/internal/impl_deps"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/grpc/generated/api/bots"
	"github.com/satont/twir/libs/grpc/generated/api/meta"
	botsGrtpc "github.com/satont/twir/libs/grpc/generated/bots"
	"github.com/satont/twir/libs/twitch"
	"github.com/twitchtv/twirp"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Bot struct {
	*impl_deps.Deps
}

func (c *Bot) BotInfo(ctx context.Context, _ *meta.BaseRequestMeta) (*bots.BotInfo, error) {
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
	result := &bots.BotInfo{
		Enabled: dbUser.Channel.IsEnabled,
	}
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
	dashboardId, ok := ctx.Value("dashboardId").(string)
	if !ok || dashboardId == "" {
		return nil, twirp.NewError(twirp.ErrorCode(twirp.Internal), "no dashboardId provided")
	}

	dbChannel := &model.Channels{}
	err := c.Db.Where(`"id" = ?`, dashboardId).Find(dbChannel).Error
	if err != nil {
		return nil, err
	}
	if dbChannel.ID == "" {
		return nil, twirp.NotFoundError("channel not found")
	}

	twitchClient, err := twitch.NewAppClientWithContext(ctx, *c.Config, c.Grpc.Tokens)
	if err != nil {
		return nil, err
	}

	twitchUsers, err := twitchClient.GetUsers(
		&helix.UsersParams{IDs: []string{dashboardId}},
	)
	if err != nil || twitchUsers.ErrorMessage != "" || len(twitchUsers.Data.Users) == 0 {
		return nil, twirp.Internal.Error("user not found on twitch")
	}

	if request.Action == bots.BotJoinPartRequest_JOIN {
		dbChannel.IsEnabled = true
	} else {
		dbChannel.IsEnabled = false
	}

	if err := c.Db.Where(`"id" = ?`, dashboardId).Select("*").Save(dbChannel).Error; err != nil {
		return nil, err
	}

	if dbChannel.IsEnabled {
		c.Grpc.Bots.Join(context.Background(), &botsGrtpc.JoinOrLeaveRequest{
			BotId:    dbChannel.BotID,
			UserName: twitchUsers.Data.Users[0].Login,
		})
	} else {
		c.Grpc.Bots.Leave(context.Background(), &botsGrtpc.JoinOrLeaveRequest{
			BotId:    dbChannel.BotID,
			UserName: twitchUsers.Data.Users[0].Login,
		})
	}

	return &emptypb.Empty{}, nil
}
