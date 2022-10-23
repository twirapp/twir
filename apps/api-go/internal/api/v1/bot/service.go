package bot

import (
	model "tsuwari/models"
	"tsuwari/twitch"

	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
	"github.com/satont/go-helix/v2"
	"github.com/satont/tsuwari/apps/api-go/internal/types"
	"github.com/satont/tsuwari/libs/nats/bots"
	"google.golang.org/protobuf/proto"
)

func handleGet(channelId string, services types.Services) (*bool, error) {
	client, err := twitch.NewUserClient(twitch.UsersServiceOpts{
		Db:           services.DB,
		ClientId:     services.Cfg.TwitchClientId,
		ClientSecret: services.Cfg.TwitchClientSecret,
	}).Create(channelId)
	if client == nil || err != nil {
		return nil, fiber.NewError(
			500,
			"cannot create twitch client from your tokens. Please try to reauthorize",
		)
	}

	channel := &model.Channels{}
	err = services.DB.Where("id = ?", channelId).First(channel).Error
	if err != nil || channel == nil {
		return nil, fiber.NewError(404, "cannot find channel in db")
	}

	mods, err := client.GetChannelMods(&helix.GetChannelModsParams{
		BroadcasterID: channelId,
		UserID:        channel.BotID,
	})
	if err != nil {
		services.Logger.Sugar().Error(err)
		return nil, fiber.NewError(500, "cannot get mods of channel")
	}

	if len(mods.Data.Mods) == 0 {
		return lo.ToPtr(false), nil
	} else {
		return lo.ToPtr(true), nil
	}
}

func handlePatch(channelId string, dto *connectionDto, services types.Services) error {
	twitchUsers, err := services.Twitch.Client.GetUsers(
		&helix.UsersParams{IDs: []string{channelId}},
	)
	if err != nil {
		services.Logger.Sugar().Error(err)
		return fiber.NewError(500, "cannot get twitch user")
	}

	if len(twitchUsers.Data.Users) == 0 {
		return fiber.NewError(404, "user not found on twitch")
	}

	user := twitchUsers.Data.Users[0]

	dbUser := model.Channels{}
	err = services.DB.Where(`"id" = ?`, channelId).First(&dbUser).Error

	if err != nil {
		services.Logger.Sugar().Error(err)
		return fiber.NewError(500, "cannot get user from database")
	}

	bytes, _ := proto.Marshal(&bots.JoinOrLeaveRequest{
		Action:   dto.Action,
		BotId:    dbUser.BotID,
		UserName: user.Login,
	})

	if dto.Action == "part" {
		dbUser.IsEnabled = false
	} else {
		dbUser.IsEnabled = true
	}

	services.DB.Where(`"id" = ?`, channelId).Select("*").Updates(&dbUser)
	services.Nats.Publish(bots.SUBJECTS_JOIN_OR_LEAVE, bytes)

	return nil
}
