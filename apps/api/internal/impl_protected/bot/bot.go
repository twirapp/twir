package bot

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/nicklaw5/helix/v2"
	"github.com/samber/lo"
	"github.com/satont/twir/apps/api/internal/impl_deps"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/twitch"
	"github.com/twirapp/twir/libs/api/messages/bots"
	"github.com/twirapp/twir/libs/api/messages/meta"
	"github.com/twirapp/twir/libs/bus-core/eventsub"
	"github.com/twitchtv/twirp"
	"golang.org/x/sync/errgroup"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Bot struct {
	*impl_deps.Deps
}

func (c *Bot) BotInfo(ctx context.Context, _ *meta.BaseRequestMeta) (*bots.BotInfo, error) {
	dashboardId, ok := ctx.Value("dashboardId").(string)
	if !ok || dashboardId == "" {
		return nil, twirp.NewError(twirp.Internal, "no dashboardId provided")
	}

	dbUser := &model.Users{}
	err := c.Db.WithContext(ctx).Where("id = ?", dashboardId).Preload("Channel").Find(dbUser).Error
	if err != nil {
		return nil, err
	}

	if dbUser.ID == "" || dbUser.Channel == nil {
		return nil, twirp.NotFoundError("user not found")
	}

	twitchClient, err := twitch.NewUserClient(dashboardId, c.Config, c.Grpc.Tokens)
	if err != nil {
		return nil, err
	}

	result := &bots.BotInfo{
		Enabled: dbUser.Channel.IsEnabled,
	}

	g, ctx := errgroup.WithContext(ctx)

	g.Go(
		func() error {
			if dashboardId == dbUser.Channel.BotID {
				result.IsMod = true
				return nil
			}

			mods, err := twitchClient.GetModerators(
				&helix.GetModeratorsParams{
					BroadcasterID: dashboardId,
					UserIDs:       []string{dbUser.Channel.BotID},
				},
			)
			if err != nil {
				return err
			}
			if mods.ErrorMessage != "" {
				return fmt.Errorf("cannot get moderators: %s", mods.ErrorMessage)
			}

			result.IsMod = lo.If(len(mods.Data.Moderators) == 0, false).Else(true)
			return nil
		},
	)

	g.Go(
		func() error {
			infoReq, err := twitchClient.GetUsers(
				&helix.UsersParams{
					IDs: []string{dbUser.Channel.BotID},
				},
			)
			if err != nil {
				return err
			}
			if len(infoReq.Data.Users) == 0 {
				return fmt.Errorf("cannot get user info: %s", infoReq.ErrorMessage)
			}

			result.BotId = infoReq.Data.Users[0].ID
			result.BotName = infoReq.Data.Users[0].Login
			return nil
		},
	)

	if err := g.Wait(); err != nil {
		return nil, fmt.Errorf("cannot get bot info: %w", err)
	}

	go func() {
		err := c.Db.Model(&model.Channels{}).Where("id = ?", dbUser.ID).Update(
			`"isBotMod"`,
			result.IsMod,
		).Error
		if err != nil {
			c.Logger.Error("cannot update channel", slog.String("channelId", dbUser.ID))
		}
	}()

	return result, nil
}

func (c *Bot) BotJoinPart(ctx context.Context, request *bots.BotJoinPartRequest) (
	*emptypb.Empty,
	error,
) {
	dashboardId, ok := ctx.Value("dashboardId").(string)
	if !ok || dashboardId == "" {
		return nil, twirp.Internal.Error("no dashboardId provided")
	}

	dbChannel := &model.Channels{}
	err := c.Db.Where(`"id" = ?`, dashboardId).Find(dbChannel).Error
	if err != nil {
		return nil, err
	}
	if dbChannel.ID == "" {
		return nil, twirp.NotFoundError("channel not found")
	}

	twitchClient, err := twitch.NewAppClientWithContext(ctx, c.Config, c.Grpc.Tokens)
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
		c.Bus.EventSub.Subscribe.Publish(
			eventsub.EventsubSubscribeRequest{ChannelID: dashboardId},
		)
	}

	return &emptypb.Empty{}, nil
}
