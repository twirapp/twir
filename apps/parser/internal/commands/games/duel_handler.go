package games

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"slices"
	"strings"
	"time"

	"github.com/goccy/go-json"
	"github.com/google/uuid"
	"github.com/guregu/null"
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

func (c *duelHandler) createHelixClient() (*helix.Client, error) {
	client, err := twitch.NewUserClient(
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

func (c *duelHandler) getSenderCurrentDuel(
	ctx context.Context,
) (*duelRedisCachedData, error) {
	keysOfAllDuelsOfChannel, err := c.parseCtx.Services.Redis.Keys(
		ctx,
		generateDuelRedisKey(c.parseCtx.Channel.ID, "*", "*"),
	).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, fmt.Errorf("cannot get keys of all duels of channel: %w", err)
	}

	if len(keysOfAllDuelsOfChannel) == 0 {
		return nil, nil
	}

	dataAsTarget := &duelRedisCachedData{}

	re := regexp.MustCompile(`commands:duel:(\d+):(\d+)`)

	var neededKey string

	for _, key := range keysOfAllDuelsOfChannel {
		matches := re.FindStringSubmatch(key)

		if len(matches) < 3 {
			continue
		}

		initiatorId := matches[1]
		targetId := matches[2]

		if initiatorId == c.parseCtx.Sender.ID {
			neededKey = key
			break
		}

		if targetId == c.parseCtx.Sender.ID {
			neededKey = key
			break
		}
	}

	if neededKey == "" {
		return nil, nil
	}

	dataAsTargetJson, err := c.parseCtx.Services.Redis.Get(
		ctx,
		neededKey,
	).Bytes()
	if err != nil {
		return nil, fmt.Errorf("cannot get duel data from cache: %w", err)
	}

	if err := json.Unmarshal(dataAsTargetJson, dataAsTarget); err != nil {
		return nil, fmt.Errorf("cannot unmarshal duel data from cache: %w", err)
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

func (c *duelHandler) saveResult(
	ctx context.Context,
	data duelRedisCachedData,
	dbChannel model.Channels,
	settings model.ChannelModulesSettingsDuel,
	loserId string,
) error {
	_, err := c.parseCtx.Services.Redis.Del(
		ctx,
		generateDuelRedisKey(c.parseCtx.Channel.ID, data.SenderID, data.TargetID),
	).Result()
	if err != nil {
		return fmt.Errorf("cannot delete duel data from cache: %w", err)
	}

	duel := model.ChannelDuel{
		ID:        uuid.New(),
		ChannelID: dbChannel.ID,
		SenderID:  null.StringFrom(data.SenderID),
		TargetID:  null.StringFrom(data.TargetID),
		LoserID:   null.StringFrom(loserId),
		CreatedAt: time.Now(),
	}

	if err := c.parseCtx.Services.Gorm.WithContext(ctx).Create(&duel).Error; err != nil {
		return fmt.Errorf("cannot save duel result: %w", err)
	}

	if settings.UserCooldown != 0 {
		_, err = c.parseCtx.Services.Redis.Set(
			ctx,
			"duels:cooldown:"+data.SenderID,
			"",
			time.Duration(settings.UserCooldown)*time.Second,
		).Result()

		if err != nil {
			return fmt.Errorf("cannot set user cooldown: %w", err)
		}
	}

	if settings.GlobalCooldown != 0 {
		_, err = c.parseCtx.Services.Redis.Set(
			ctx,
			"duels:cooldown:global",
			"",
			time.Duration(settings.GlobalCooldown)*time.Second,
		).Result()

		if err != nil {
			return fmt.Errorf("cannot set global cooldown: %w", err)
		}
	}

	return nil
}

func (c *duelHandler) isCooldown(ctx context.Context, userID string) (bool, error) {
	isUserCooldown, err := c.parseCtx.Services.Redis.Exists(
		ctx,
		"duels:cooldown:"+userID,
	).Result()
	if err != nil {
		return false, fmt.Errorf("cannot check cooldown: %w", err)
	}

	if isUserCooldown == 1 {
		return true, nil
	}

	isGlobalCooldown, err := c.parseCtx.Services.Redis.Exists(
		ctx,
		"duels:cooldown:global",
	).Result()
	if err != nil {
		return false, fmt.Errorf("cannot check cooldown: %w", err)
	}

	if isGlobalCooldown == 1 {
		return true, nil
	}

	return false, nil
}
