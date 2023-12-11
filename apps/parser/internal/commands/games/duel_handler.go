package games

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/goccy/go-json"
	"github.com/nicklaw5/helix/v2"
	"github.com/redis/go-redis/v9"
	"github.com/samber/lo"
	"github.com/satont/twir/apps/parser/internal/types"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/twitch"
)

type duelHandler struct {
	parseCtx    *types.ParseContext
	helixClient *helix.Client
}

type duelRedisCachedData struct {
	TargetUserLogin   string `json:"targetUserLogin"`
	TargetID          string `json:"targetId"`
	IsTargetModerator bool   `json:"isTargetModerator"`

	SenderUserLogin   string `json:"senderUserLogin"`
	SenderID          string `json:"senderId"`
	IsSenderModerator bool   `json:"isSenderModerator"`
}

func (c *duelHandler) getChannelSettings(ctx context.Context) (
	model.ChannelModulesSettingsDuel,
	error,
) {
	entity := model.ChannelModulesSettings{}
	var parsedSettings model.ChannelModulesSettingsDuel

	if err := c.parseCtx.Services.Gorm.WithContext(ctx).Where(
		`"channelId" = ? and "userId" is null and "type" = 'duel'`,
		c.parseCtx.Channel.ID,
	).First(&entity).Error; err != nil {
		return parsedSettings, err
	}

	if err := json.Unmarshal(entity.Settings, &parsedSettings); err != nil {
		return parsedSettings, err
	}

	return parsedSettings, nil
}

func (c *duelHandler) createHelixClient(ctx context.Context) (*helix.Client, error) {
	client, err := twitch.NewUserClientWithContext(
		ctx,
		c.parseCtx.Channel.ID,
		*c.parseCtx.Services.Config,
		c.parseCtx.Services.GrpcClients.Tokens,
	)
	if err != nil {
		return nil, err
	}

	c.helixClient = client

	return client, nil
}

func (c *duelHandler) getTwitchTargetUser() (helix.User, error) {
	targetUserName := strings.Replace(*c.parseCtx.Text, "@", "", 1)

	userRequest, err := c.helixClient.GetUsers(&helix.UsersParams{Logins: []string{targetUserName}})
	if err != nil {
		return helix.User{}, fmt.Errorf("cannot get user: %w", err)
	}
	if userRequest.ErrorMessage != "" {
		return helix.User{}, errors.New(userRequest.ErrorMessage)
	}

	if len(userRequest.Data.Users) == 0 {
		return helix.User{}, errors.New("user not found")
	}

	return userRequest.Data.Users[0], nil
}

func (c *duelHandler) getDbChannel(ctx context.Context) (model.Channels, error) {
	channel := model.Channels{}
	if err := c.parseCtx.Services.Gorm.WithContext(ctx).Where(
		`"id" = ?`,
		c.parseCtx.Channel.ID,
	).First(&channel).Error; err != nil {
		return model.Channels{}, err
	}

	return channel, nil
}

func (c *duelHandler) validateTarget(
	ctx context.Context,
	targetUser helix.User,
	dbChannel model.Channels,
) error {
	if targetUser.ID == c.parseCtx.Sender.ID {
		return errors.New("you cannot duel with yourself")
	}

	if targetUser.ID == c.parseCtx.Channel.ID {
		return errors.New("you cannot duel with streamer")
	}

	if dbChannel.BotID == targetUser.ID {
		return errors.New("you cannot duel with bot")
	}

	isAlreadyParticipant, err := c.parseCtx.Services.Redis.Exists(
		ctx,
		generateDuelRedisKey(c.parseCtx.Channel.ID, c.parseCtx.Sender.ID, "*"),
	).Result()
	if err != nil {
		return fmt.Errorf("cannot check target user: %w", err)
	}

	if isAlreadyParticipant == 1 {
		return errors.New("you are already participating in the duel")
	}

	isAlreadyParticipant, err = c.parseCtx.Services.Redis.Exists(
		ctx,
		generateDuelRedisKey(c.parseCtx.Channel.ID, targetUser.ID, "*"),
	).Result()
	if err != nil {
		return fmt.Errorf("cannot check target user: %w", err)
	}

	if isAlreadyParticipant == 1 {
		return errors.New("target user is already participating in the duel")
	}

	return nil
}

func (c *duelHandler) getChannelModerators() ([]helix.Moderator, error) {
	moderatorsRequest, err := c.helixClient.GetModerators(
		&helix.GetModeratorsParams{
			BroadcasterID: c.parseCtx.Channel.ID,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("cannot get moderators: %w", err)
	}
	if moderatorsRequest.ErrorMessage != "" {
		return nil, errors.New(moderatorsRequest.ErrorMessage)
	}

	return moderatorsRequest.Data.Moderators, nil
}

func (c *duelHandler) saveDuelDataToCache(
	ctx context.Context,
	targetUser helix.User,
	moderators []helix.Moderator,
	settings model.ChannelModulesSettingsDuel,
) error {
	redisCachedData := duelRedisCachedData{
		TargetUserLogin: targetUser.Login,
		TargetID:        targetUser.ID,
		IsTargetModerator: lo.SomeBy(
			moderators,
			func(item helix.Moderator) bool {
				return item.UserID == targetUser.ID
			},
		),

		SenderUserLogin:   c.parseCtx.Sender.Name,
		SenderID:          c.parseCtx.Sender.ID,
		IsSenderModerator: slices.Contains(c.parseCtx.Sender.Badges, "moderator"),
	}

	redisCachedDataJson, err := json.Marshal(redisCachedData)
	if err != nil {
		return fmt.Errorf("internal error when trying to find cached data: %w", err)
	}

	err = c.parseCtx.Services.Redis.Set(
		ctx,
		generateDuelRedisKey(c.parseCtx.Channel.ID, c.parseCtx.Sender.ID, targetUser.ID),
		redisCachedDataJson,
		time.Duration(settings.SecondsToAccept)*time.Second,
	).Err()
	if err != nil {
		return fmt.Errorf("cannot save duel data to cache: %w", err)
	}

	return nil
}

func (c *duelHandler) getCurrentDuelOfUserByKey(
	ctx context.Context,
	key string,
) (duelRedisCachedData, error) {
	cachedData, err := c.parseCtx.Services.Redis.Get(
		ctx,
		key,
	).Bytes()
	if err != nil && errors.Is(err, redis.Nil) {
		return duelRedisCachedData{}, fmt.Errorf(
			"internal error when trying to find cached data: %w",
			err,
		)
	}

	data := duelRedisCachedData{}

	if cachedData != nil {
		if err := json.Unmarshal(cachedData, &data); err != nil {
			return data, fmt.Errorf("internal error when trying to find cached data: %w", err)
		}

		return data, nil
	}

	return data, nil
}

func (c *duelHandler) getSenderCurrentDuel(
	ctx context.Context,
) (duelRedisCachedData, error) {
	dataAsInitiator, err := c.getCurrentDuelOfUserByKey(
		ctx,
		generateDuelRedisKey(c.parseCtx.Channel.ID, c.parseCtx.Sender.ID, "*"),
	)
	if err != nil {
		return duelRedisCachedData{}, err
	}
	if dataAsInitiator.SenderID != "" {
		return dataAsInitiator, nil
	}

	dataAsTarget, err := c.getCurrentDuelOfUserByKey(
		ctx,
		generateDuelRedisKey(c.parseCtx.Channel.ID, "*", c.parseCtx.Sender.ID),
	)
	if err != nil {
		return duelRedisCachedData{}, err
	}
	if dataAsTarget.TargetID != "" {
		return dataAsTarget, nil
	}

	return dataAsTarget, nil
}

func (c *duelHandler) timeoutUser(
	data duelRedisCachedData,
	dbChannel model.Channels,
	settings model.ChannelModulesSettingsDuel,
	userID string,
	isMod bool,
) error {
	if data.TargetID != userID && data.SenderID != userID {
		return errors.New("user not cached")
	}

	if isMod {
		_, err := c.helixClient.RemoveChannelModerator(
			&helix.RemoveChannelModeratorParams{
				BroadcasterID: c.parseCtx.Channel.ID,
				UserID:        userID,
			},
		)
		if err != nil {
			return fmt.Errorf("cannot remove moderator")
		}
	}

	_, err := c.helixClient.BanUser(
		&helix.BanUserParams{
			BroadcasterID: dbChannel.ID,
			ModeratorId:   dbChannel.ID,
			Body: helix.BanUserRequestBody{
				Duration: int(settings.TimeoutSeconds),
				Reason:   "lost in duel",
				UserId:   userID,
			},
		},
	)
	if err != nil {
		return fmt.Errorf("cannot ban user")
	}

	return nil
}
