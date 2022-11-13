package auth

import (
	"database/sql"
	"time"

	model "github.com/satont/tsuwari/libs/gomodels"

	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
	"github.com/satont/go-helix/v2"
	"github.com/satont/tsuwari/apps/api/internal/types"
	"github.com/satont/tsuwari/libs/nats/bots"
	"github.com/satont/tsuwari/libs/nats/eventsub"
	"github.com/satont/tsuwari/libs/nats/scheduler"
	uuid "github.com/satori/go.uuid"
	"google.golang.org/protobuf/proto"
	"gorm.io/gorm"
)

func checkUser(
	username, userId string,
	tokens helix.AccessCredentials,
	services types.Services,
) error {
	defaultBot := model.Bots{}
	err := services.DB.Where("type = ?", "DEFAULT").First(&defaultBot).Error
	if err != nil {
		return fiber.NewError(500, "bot not created, cannot create user")
	}

	tokenData := model.Tokens{
		AccessToken:         tokens.AccessToken,
		RefreshToken:        tokens.RefreshToken,
		ExpiresIn:           int32(tokens.ExpiresIn),
		ObtainmentTimestamp: time.Now(),
	}

	user := model.Users{}
	err = services.DB.
		Where(`"users"."id" = ?`, userId).
		Joins("Channel").
		Joins("Token").
		First(&user).Error

	if err != nil && err == gorm.ErrRecordNotFound {
		newToken := tokenData
		newToken.ID = uuid.NewV4().String()

		if err = services.DB.Create(&newToken).Error; err != nil {
			services.Logger.Sugar().Error(err)
			return err
		}

		user.ID = userId
		user.TokenID = sql.NullString{String: newToken.ID, Valid: true}
		user.ApiKey = uuid.NewV4().String()

		if err = services.DB.Create(&user).Error; err != nil {
			services.Logger.Sugar().Error(err)
			return err
		}

		channel := createChannelModel(user.ID, defaultBot.ID)
		if err = services.DB.Create(&channel).Error; err != nil {
			services.Logger.Sugar().Error(err)
			return err
		}
		user.Channel = &channel
	}
	if err != nil && err != gorm.ErrRecordNotFound {
		services.Logger.Sugar().Error(err)
		return fiber.NewError(500, "internal error")
	} else {
		if user.Channel == nil {
			channel := createChannelModel(user.ID, defaultBot.ID)
			if err = services.DB.Create(&channel).Error; err != nil {
				services.Logger.Sugar().Error(err)
				return err
			}
			user.Channel = &channel
		}

		if user.TokenID.Valid {
			tokenData.ID = user.TokenID.String
			if err = services.DB.Select("*").Updates(&tokenData).Error; err != nil {
				services.Logger.Sugar().Error(err)
				return err
			}
		} else {
			tokenData.ID = uuid.NewV4().String()
			if err = services.DB.Select("*").Create(&tokenData).Error; err != nil {
				services.Logger.Sugar().Error(err)
				return err
			}
		}
	}

	bytes, _ := proto.Marshal(&bots.JoinOrLeaveRequest{
		Action:   lo.If(user.Channel.IsEnabled, "join").Else("part"),
		BotId:    user.Channel.BotID,
		UserName: username,
	})
	services.Nats.Publish(bots.SUBJECTS_JOIN_OR_LEAVE, bytes)

	bytes, _ = proto.Marshal(&scheduler.CreateDefaultCommandsRequest{
		UserId: userId,
	})
	services.Nats.Publish(scheduler.SUBJECTS_CREATE_DEFAULT_COMMANDS, bytes)

	bytes, _ = proto.Marshal(&eventsub.SubscribeToEvents{
		ChannelId: userId,
	})
	services.Nats.Publish(eventsub.SUBJECTS_SUBSCTUBE_TO_EVENTS_BY_CHANNEL_ID, bytes)

	return nil
}

func createChannelModel(userId, botId string) model.Channels {
	return model.Channels{
		ID:             userId,
		IsEnabled:      true,
		IsTwitchBanned: false,
		IsBanned:       false,
		BotID:          botId,
	}
}
