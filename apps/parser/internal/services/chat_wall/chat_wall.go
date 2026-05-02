package chat_wall

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/parser/locales"
	buscore "github.com/twirapp/twir/libs/bus-core"
	botsservice "github.com/twirapp/twir/libs/bus-core/bots"
	generic_cacher "github.com/twirapp/twir/libs/cache/generic-cacher"
	config "github.com/twirapp/twir/libs/config"
	platformentity "github.com/twirapp/twir/libs/entities/platform"
	deprecatedgormmodel "github.com/twirapp/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/i18n"
	"github.com/twirapp/twir/libs/redis_keys"
	"github.com/twirapp/twir/libs/repositories/chat_messages"
	chatmessagesrepository "github.com/twirapp/twir/libs/repositories/chat_messages"
	chatmessagemodel "github.com/twirapp/twir/libs/repositories/chat_messages/model"
	chatwallrepository "github.com/twirapp/twir/libs/repositories/chat_wall"
	"github.com/twirapp/twir/libs/repositories/chat_wall/model"
	"gorm.io/gorm"
)

type Opts struct {
	ChatWallRepository chatwallrepository.Repository
	ChatMessagesRepo   chat_messages.Repository

	Gorm          *gorm.DB
	ChatWallCache *generic_cacher.GenericCacher[[]model.ChatWall]
	Redis         *redis.Client
	Config        config.Config
	TwirBus       *buscore.Bus
}

func New(opts Opts) *Service {
	return &Service{
		repo:             opts.ChatWallRepository,
		chatMessagesRepo: opts.ChatMessagesRepo,
		gorm:             opts.Gorm,
		chatWallCache:    opts.ChatWallCache,
		redis:            opts.Redis,
		config:           opts.Config,
		twirBus:          opts.TwirBus,
	}
}

type Service struct {
	repo             chatwallrepository.Repository
	chatMessagesRepo chat_messages.Repository
	gorm             *gorm.DB
	chatWallCache    *generic_cacher.GenericCacher[[]model.ChatWall]
	redis            *redis.Client
	config           config.Config
	twirBus          *buscore.Bus
}

type CreateInput struct {
	DBChannelID     string
	Phrase          string
	Enabled         bool
	Action          model.ChatWallAction
	Duration        time.Duration
	TimeoutDuration *time.Duration
}

func (c *Service) Create(ctx context.Context, input CreateInput) (model.ChatWall, error) {
	parsedChannelID, err := uuid.Parse(input.DBChannelID)
	if err != nil {
		return model.ChatWall{}, err
	}

	currentChatWallsEnabledParam := true
	currentChatWalls, err := c.repo.GetMany(
		ctx,
		chatwallrepository.GetManyInput{
			ChannelID: parsedChannelID,
			Enabled:   &currentChatWallsEnabledParam,
		},
	)
	if err != nil {
		return model.ChatWall{}, fmt.Errorf(
			i18n.GetCtx(
				ctx,
				locales.Translations.Services.ChatWall.Errors.GetCurrentChatWalls.
					SetVars(locales.KeysServicesChatWallErrorsGetCurrentChatWallsVars{Reason: err.Error()}),
			),
		)
	}

	for _, chatWall := range currentChatWalls {
		if chatWall.Phrase == input.Phrase {
			return model.ChatWall{}, fmt.Errorf(
				i18n.GetCtx(
					ctx,
					locales.Translations.Services.ChatWall.Errors.CreateChatWallWithPhrase,
				),
			)
		}
	}

	wall, err := c.repo.Create(
		ctx,
		chatwallrepository.CreateInput{
			ChannelID:       parsedChannelID,
			Phrase:          input.Phrase,
			Enabled:         true,
			Action:          input.Action,
			Duration:        input.Duration,
			TimeoutDuration: input.TimeoutDuration,
		},
	)
	if err != nil {
		return model.ChatWall{}, fmt.Errorf(
			i18n.GetCtx(
				ctx,
				locales.Translations.Services.ChatWall.Errors.CreateChatWall.
					SetVars(locales.KeysServicesChatWallErrorsCreateChatWallVars{Reason: err.Error()}),
			),
		)
	}

	c.chatWallCache.Invalidate(ctx, input.DBChannelID)

	return wall, nil
}

type HandlePastMessagesInput struct {
	DBChannelID     string
	PlatformChannelID string
	Platform        platformentity.Platform
	Phrase          string
	Action          model.ChatWallAction
	TimeoutDuration *time.Duration
}

func (c *Service) HandlePastMessages(
	ctx context.Context,
	wall model.ChatWall,
	input HandlePastMessagesInput,
) error {
	parsedChannelID, err := uuid.Parse(input.DBChannelID)
	if err != nil {
		return err
	}

	chatWallSettings, err := c.repo.GetChannelSettings(ctx, parsedChannelID)
	if err != nil && !errors.Is(err, chatwallrepository.ErrSettingsNotFound) {
		return fmt.Errorf(
			i18n.GetCtx(
				ctx,
				locales.Translations.Services.ChatWall.Errors.GetChatWallSettings.
					SetVars(locales.KeysServicesChatWallErrorsGetChatWallSettingsVars{Reason: err.Error()}),
			),
		)
	}

	if chatWallSettings.IsNil() {
		return nil
	}

	timeGte := time.Now().Add(-10 * time.Minute)

	messages, err := c.chatMessagesRepo.GetMany(
		ctx,
		chatmessagesrepository.GetManyInput{
			ChannelID: &input.PlatformChannelID,
			TextLike:  &input.Phrase,
			Page:      0,
			PerPage:   1000,
			TimeGte:   &timeGte,
		},
	)
	if err != nil {
		return err
	}
	if len(messages) == 0 {
		return nil
	}

	var isSubscribersMuted bool
	if chatWallSettings != model.ChatWallSettingsNil {
		isSubscribersMuted = chatWallSettings.MuteSubscribers
	} else {
		isSubscribersMuted = true
	}

	var isVipsMuted bool
	if chatWallSettings != model.ChatWallSettingsNil {
		isVipsMuted = chatWallSettings.MuteVips
	} else {
		isVipsMuted = false
	}

	type userPlatformRef struct {
		ID         string `gorm:"column:id"`
		PlatformID string `gorm:"column:platform_id"`
	}

	messagePlatformUserIDs := lo.Map(messages, func(item chatmessagemodel.ChatMessage, _ int) string {
		return item.UserID
	})

	var userRefs []userPlatformRef
	if err := c.gorm.
		WithContext(ctx).
		Table("users").
		Select("id", "platform_id").
		Where("platform = ? AND platform_id IN ?", input.Platform, messagePlatformUserIDs).
		Find(&userRefs).Error; err != nil {
		return fmt.Errorf(
			i18n.GetCtx(
				ctx,
				locales.Translations.Services.ChatWall.Errors.GetUsersStats.
					SetVars(locales.KeysServicesChatWallErrorsGetUsersStatsVars{Reason: err.Error()}),
			),
		)
	}

	messageUserIDs := make([]uuid.UUID, 0, len(userRefs))
	platformUserIDByInternalID := make(map[string]string, len(userRefs))
	internalUserIDByPlatformID := make(map[string]uuid.UUID, len(userRefs))
	for _, ref := range userRefs {
		parsedUserID, err := uuid.Parse(ref.ID)
		if err != nil {
			continue
		}

		messageUserIDs = append(messageUserIDs, parsedUserID)
		platformUserIDByInternalID[ref.ID] = ref.PlatformID
		internalUserIDByPlatformID[ref.PlatformID] = parsedUserID
	}

	var usersStats []deprecatedgormmodel.UsersStats
	if err := c.gorm.
		WithContext(ctx).
		Where(
			`"userId" IN ? AND "channelId" = ?::uuid`,
			messageUserIDs,
			parsedChannelID,
		).
		Find(&usersStats).Error; err != nil {
		return fmt.Errorf(
			i18n.GetCtx(
				ctx,
				locales.Translations.Services.ChatWall.Errors.GetUsersStats.
					SetVars(locales.KeysServicesChatWallErrorsGetUsersStatsVars{Reason: err.Error()}),
			),
		)
	}

	messages = lo.Filter(
		messages,
		func(item chatmessagemodel.ChatMessage, _ int) bool {
			var foundUserStats *deprecatedgormmodel.UsersStats
			for _, userStats := range usersStats {
				if platformUserIDByInternalID[userStats.UserID] == item.UserID {
					foundUserStats = &userStats
					break
				}
			}

			if foundUserStats == nil {
				return true
			}

			if !isSubscribersMuted && foundUserStats.IsSubscriber {
				return false
			}

			if !isVipsMuted && foundUserStats.IsVip {
				return false
			}

			return true
		},
	)

	alreadyHandledMessagesIds, err := c.redis.SMembers(
		ctx,
		fmt.Sprintf(redis_keys.NukeRedisPrefix, input.DBChannelID),
	).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return fmt.Errorf(
			i18n.GetCtx(
				ctx,
				locales.Translations.Services.ChatWall.Errors.GetAlreadyHandled.
					SetVars(locales.KeysServicesChatWallErrorsGetAlreadyHandledVars{Reason: err.Error()}),
			),
		)
	}

	if input.Action == model.ChatWallActionDelete {
		mappedMessagesIDs := make([]string, 0, len(messages))
		for _, m := range messages {
			if !slices.Contains(alreadyHandledMessagesIds, m.ID.String()) {
				mappedMessagesIDs = append(mappedMessagesIDs, m.ID.String())
			}
		}

		if len(mappedMessagesIDs) == 0 {
			return nil
		}

			err = c.twirBus.Bots.DeleteMessage.Publish(
				ctx,
				botsservice.DeleteMessageRequest{
					ChannelId:  input.PlatformChannelID,
					MessageIds: mappedMessagesIDs,
			},
		)
		if err != nil {
			return fmt.Errorf(
				i18n.GetCtx(
					ctx,
					locales.Translations.Services.ChatWall.Errors.PublishDeletedMessages.
						SetVars(locales.KeysServicesChatWallErrorsPublishDeletedMessagesVars{Reason: err.Error()}),
				),
			)
		}
	} else if input.Action == model.ChatWallActionBan || input.Action == model.ChatWallActionTimeout {
		request := make([]botsservice.BanRequest, 0, len(messages))

		var banTime int
		if input.Action == model.ChatWallActionTimeout {
			banTime = int(input.TimeoutDuration.Seconds())
		}

		for _, m := range messages {
			if slices.Contains(alreadyHandledMessagesIds, m.ID.String()) {
				continue
			}

			request = append(
				request,
				botsservice.BanRequest{
					ChannelID: input.PlatformChannelID,
					UserID:    m.UserID,
					Reason: i18n.GetCtx(
						ctx,
						locales.Translations.Services.ChatWall.Info.BannedByTwir.SetVars(locales.KeysServicesChatWallInfoBannedByTwirVars{BanPhrase: input.Phrase}),
					),
					BanTime:        banTime,
					IsModerator:    false,
					AddModAfterBan: false,
				},
			)
		}

		err = c.twirBus.Bots.BanUsers.Publish(ctx, request)
		if err != nil {
			return fmt.Errorf(
				i18n.GetCtx(
					ctx,
					locales.Translations.Services.ChatWall.Errors.PublishBanUsers.SetVars(locales.KeysServicesChatWallErrorsPublishBanUsersVars{Reason: err.Error()}),
				),
			)
		}
	}

	newHandledMessagesIds := make([]string, 0, len(messages))
	for _, m := range messages {
		if !slices.Contains(alreadyHandledMessagesIds, m.ID.String()) {
			newHandledMessagesIds = append(newHandledMessagesIds, m.ID.String())
		}
	}

	if len(newHandledMessagesIds) == 0 {
		return nil
	}

	logs := make([]chatwallrepository.CreateLogInput, 0, len(newHandledMessagesIds))
	for _, m := range messages {
		if slices.Contains(alreadyHandledMessagesIds, m.ID.String()) {
			continue
		}

		internalUserID, ok := internalUserIDByPlatformID[m.UserID]
		if !ok || internalUserID == uuid.Nil {
			continue
		}
		logs = append(
			logs,
			chatwallrepository.CreateLogInput{
				WallID: wall.ID,
				UserID: internalUserID,
				Text:   m.Text,
			},
		)
	}

	if err := c.repo.CreateManyLogs(ctx, logs); err != nil {
		return fmt.Errorf(
			i18n.GetCtx(
				ctx,
				locales.Translations.Services.ChatWall.Errors.CreateChatLogsInDb.
					SetVars(locales.KeysServicesChatWallErrorsCreateChatLogsInDbVars{Reason: err.Error()}),
			),
		)
	}

	_, err = c.redis.Pipelined(
		ctx, func(p redis.Pipeliner) error {
			if err := p.SAdd(
				ctx,
				fmt.Sprintf(redis_keys.NukeRedisPrefix, input.DBChannelID),
				newHandledMessagesIds,
			).Err(); err != nil {
				return err
			}
			if err := p.Expire(
				ctx,
				fmt.Sprintf(redis_keys.NukeRedisPrefix, input.DBChannelID),
				20*time.Minute,
			).Err(); err != nil {
				return err
			}

			return nil
		},
	)
	if err != nil {
		return fmt.Errorf(
			i18n.GetCtx(
				ctx,
				locales.Translations.Services.ChatWall.Errors.HandledMessagesToRedis.
					SetVars(locales.KeysServicesChatWallErrorsHandledMessagesToRedisVars{Reason: err.Error()}),
			),
		)
	}

	return nil
}

type StopInput struct {
	DBChannelID string
	Phrase      string
}

var ErrChatWallNotFound = errors.New("chat wall not found")

func (c *Service) Stop(ctx context.Context, input StopInput) error {
	parsedChannelID, err := uuid.Parse(input.DBChannelID)
	if err != nil {
		return err
	}

	enabled := true

	walls, err := c.repo.GetMany(
		ctx,
		chatwallrepository.GetManyInput{
			ChannelID: parsedChannelID,
			Enabled:   &enabled,
		},
	)
	if err != nil {
		return fmt.Errorf(
			i18n.GetCtx(
				ctx,
				locales.Translations.Services.ChatWall.Errors.GetChatWalls.
					SetVars(locales.KeysServicesChatWallErrorsGetChatWallsVars{Reason: err.Error()}),
			),
		)
	}

	for _, wall := range walls {
		if wall.Phrase == input.Phrase {
			_, err = c.repo.Update(
				ctx,
				wall.ID,
				chatwallrepository.UpdateInput{
					Enabled: lo.ToPtr(false),
				},
			)
			if err != nil {
				return fmt.Errorf(
					i18n.GetCtx(
						ctx,
						locales.Translations.Services.ChatWall.Errors.UpdateChatWalls.
							SetVars(locales.KeysServicesChatWallErrorsUpdateChatWallsVars{Reason: err.Error()}),
					),
				)
			}

			c.chatWallCache.Invalidate(ctx, input.DBChannelID)
			return nil
		}
	}

	return ErrChatWallNotFound
}
