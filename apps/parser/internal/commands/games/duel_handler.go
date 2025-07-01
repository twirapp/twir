package games

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/guregu/null"
	"github.com/nicklaw5/helix/v2"
	"github.com/satont/twir/apps/parser/internal/types"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/twitch"
	"github.com/twirapp/twir/libs/bus-core/bots"
	"gorm.io/gorm"
)

type duelHandler struct {
	parseCtx    *types.ParseContext
	helixClient *helix.Client
}

func (c *duelHandler) getChannelSettings(ctx context.Context) (
	model.ChannelGamesDuel,
	error,
) {
	entity := model.ChannelGamesDuel{}

	if err := c.parseCtx.Services.Gorm.WithContext(ctx).Where(
		`"channel_id" = ?`,
		c.parseCtx.Channel.ID,
	).First(&entity).Error; err != nil {
		return entity, err
	}

	return entity, nil
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
	targetUserName := strings.Replace(
		c.parseCtx.ArgsParser.Get(duelTargetArgName).String(),
		"@",
		"",
		1,
	)

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

func (c *duelHandler) getUserCurrentDuel(ctx context.Context, userId string) (
	*model.ChannelDuel,
	error,
) {
	duel := model.ChannelDuel{}

	err := c.parseCtx.Services.Gorm.
		WithContext(ctx).
		Debug().
		Where(`channel_id = ?`, c.parseCtx.Channel.ID).
		Where("finished_at is null").
		Where("available_until >= now()").
		Where("sender_id = ? OR target_id = ?", userId, userId).
		First(&duel).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("cannot check user in duel: %w", err)
	}

	return &duel, nil
}

type targetValidateError struct {
	message string
}

func (e *targetValidateError) Error() string {
	return e.message
}

func (c *duelHandler) validateParticipants(
	ctx context.Context,
	senderUserId string,
	targetUserId string,
	dbChannel model.Channels,
) error {
	if targetUserId == c.parseCtx.Sender.ID {
		return &targetValidateError{message: "you cannot duel with yourself"}
	}
	if targetUserId == c.parseCtx.Channel.ID {
		return &targetValidateError{message: "you cannot duel with streamer"}
	}
	if dbChannel.BotID == targetUserId {
		return &targetValidateError{message: "you cannot duel with bot"}
	}

	targetDuelUser, err := c.getUserCurrentDuel(ctx, targetUserId)
	if err != nil {
		return fmt.Errorf("cannot check user in duel: %w", err)
	}
	if targetDuelUser != nil {
		return &targetValidateError{message: "target user already in duel"}
	}

	senderDuel, err := c.getUserCurrentDuel(ctx, senderUserId)
	if err != nil {
		return fmt.Errorf("cannot check user in duel: %w", err)
	}
	if senderDuel != nil {
		return &targetValidateError{message: "you already in duel"}
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

func (c *duelHandler) saveDuelData(
	ctx context.Context,
	targetUser helix.User,
	moderators []helix.Moderator,
	settings model.ChannelGamesDuel,
) error {
	var senderModerator bool
	var targetModerator bool
	for _, moderator := range moderators {
		if moderator.UserID == c.parseCtx.Sender.ID {
			senderModerator = true
		}
		if moderator.UserID == targetUser.ID {
			targetModerator = true
		}
	}

	entity := model.ChannelDuel{
		ID:              uuid.New(),
		ChannelID:       c.parseCtx.Channel.ID,
		SenderID:        null.StringFrom(c.parseCtx.Sender.ID),
		SenderModerator: senderModerator,
		SenderLogin:     c.parseCtx.Sender.Name,
		TargetID:        null.StringFrom(targetUser.ID),
		TargetModerator: targetModerator,
		TargetLogin:     targetUser.Login,
		LoserID:         null.String{},
		CreatedAt:       time.Now(),
		FinishedAt:      null.Time{},
		AvailableUntil:  time.Now().Add(time.Duration(settings.SecondsToAccept) * time.Second),
	}

	if err := c.parseCtx.Services.Gorm.WithContext(ctx).Create(&entity).Error; err != nil {
		return fmt.Errorf("cannot save duel data: %w", err)
	}

	return nil
}

func (c *duelHandler) timeoutUser(
	ctx context.Context,
	settings model.ChannelGamesDuel,
	userID string,
	isMod bool,
) error {
	return c.parseCtx.Services.Bus.Bots.BanUser.Publish(
		ctx,
		bots.BanRequest{
			ChannelID:      c.parseCtx.Channel.ID,
			UserID:         userID,
			Reason:         "lost in duel",
			BanTime:        int(settings.TimeoutSeconds),
			IsModerator:    isMod,
			AddModAfterBan: true,
		},
	)
}

func (c *duelHandler) saveResult(
	ctx context.Context,
	data model.ChannelDuel,
	settings model.ChannelGamesDuel,
	loserId null.String,
) error {

	data.LoserID = loserId
	data.FinishedAt = null.TimeFrom(time.Now())

	if err := c.parseCtx.Services.Gorm.WithContext(ctx).Save(&data).Error; err != nil {
		return fmt.Errorf("cannot save duel result: %w", err)
	}

	if settings.UserCooldown != 0 {
		_, err := c.parseCtx.Services.Redis.Set(
			ctx,
			"duels:cooldown:"+data.SenderID.String,
			"",
			time.Duration(settings.UserCooldown)*time.Second,
		).Result()

		if err != nil {
			return fmt.Errorf("cannot set user cooldown: %w", err)
		}
	}

	if settings.GlobalCooldown != 0 {
		_, err := c.parseCtx.Services.Redis.Set(
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
