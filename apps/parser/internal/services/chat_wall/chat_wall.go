package chat_wall

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/samber/lo"
	config "github.com/satont/twir/libs/config"
	deprecatedgormmodel "github.com/satont/twir/libs/gomodels"
	buscore "github.com/twirapp/twir/libs/bus-core"
	botsservice "github.com/twirapp/twir/libs/bus-core/bots"
	generic_cacher "github.com/twirapp/twir/libs/cache/generic-cacher"
	"github.com/twirapp/twir/libs/grpc/tokens"
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
	TokensClient  tokens.TokensClient
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
		tokens:           opts.TokensClient,
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
	tokens           tokens.TokensClient
	twirBus          *buscore.Bus
}

type CreateInput struct {
	ChannelID       string
	Phrase          string
	Enabled         bool
	Action          model.ChatWallAction
	Duration        time.Duration
	TimeoutDuration *time.Duration
}

func (c *Service) Create(ctx context.Context, input CreateInput) (model.ChatWall, error) {
	currentChatWallsEnabledParam := true
	currentChatWalls, err := c.repo.GetMany(
		ctx,
		chatwallrepository.GetManyInput{
			ChannelID: input.ChannelID,
			Enabled:   &currentChatWallsEnabledParam,
		},
	)
	if err != nil {
		return model.ChatWall{}, fmt.Errorf("cannot get current chat walls: %s", err)
	}

	for _, chatWall := range currentChatWalls {
		if chatWall.Phrase == input.Phrase {
			return model.ChatWall{}, fmt.Errorf("cannot create chat wall with phrase that already exists")
		}
	}

	wall, err := c.repo.Create(
		ctx,
		chatwallrepository.CreateInput{
			ChannelID:       input.ChannelID,
			Phrase:          input.Phrase,
			Enabled:         true,
			Action:          input.Action,
			Duration:        10 * time.Minute,
			TimeoutDuration: input.TimeoutDuration,
		},
	)
	if err != nil {
		return model.ChatWall{}, fmt.Errorf("cannot create chat wall: %s", err)
	}

	c.chatWallCache.Invalidate(ctx, input.ChannelID)

	return wall, nil
}

type HandlePastMessagesInput struct {
	ChannelID       string
	Phrase          string
	Action          model.ChatWallAction
	TimeoutDuration *time.Duration
}

func (c *Service) HandlePastMessages(
	ctx context.Context,
	wall model.ChatWall,
	input HandlePastMessagesInput,
) error {
	chatWallSettings, err := c.repo.GetChannelSettings(ctx, input.ChannelID)
	if err != nil && !errors.Is(err, chatwallrepository.ErrSettingsNotFound) {
		return fmt.Errorf("cannot get chat wall settings: %s", err)
	}

	timeGte := time.Now().Add(-10 * time.Minute)

	messages, err := c.chatMessagesRepo.GetMany(
		ctx,
		chatmessagesrepository.GetManyInput{
			ChannelID: &input.ChannelID,
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

	var usersStats []deprecatedgormmodel.UsersStats
	if err := c.gorm.
		WithContext(ctx).
		Where(
			`"userId" IN ? AND "channelId" = ?`,
			lo.Map(
				messages, func(item chatmessagemodel.ChatMessage, _ int) string {
					return item.UserID
				},
			), input.ChannelID,
		).
		Find(&usersStats).Error; err != nil {
		return fmt.Errorf("cannot get users stats: %s", err)
	}

	messages = lo.Filter(
		messages,
		func(item chatmessagemodel.ChatMessage, _ int) bool {
			var foundUserStats *deprecatedgormmodel.UsersStats
			for _, userStats := range usersStats {
				if userStats.UserID == item.UserID {
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
		fmt.Sprintf(redis_keys.NukeRedisPrefix, input.ChannelID),
	).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return fmt.Errorf("cannot get already handled messages: %s", err)
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
			botsservice.DeleteMessageRequest{
				ChannelId:  input.ChannelID,
				MessageIds: mappedMessagesIDs,
			},
		)
		if err != nil {
			return fmt.Errorf("cannot publish delete messages: %s", err)
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
					ChannelID:      input.ChannelID,
					UserID:         m.UserID,
					Reason:         fmt.Sprintf("banned by twir for chat wall phrase: %s", input.Phrase),
					BanTime:        banTime,
					IsModerator:    false,
					AddModAfterBan: false,
				},
			)
		}

		err = c.twirBus.Bots.BanUsers.Publish(request)
		if err != nil {
			return fmt.Errorf("cannot publish ban users: %s", err)
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

	logs := make([]chatwallrepository.CreateLogInput, 0, len(messages))
	for _, m := range messages {
		logs = append(
			logs,
			chatwallrepository.CreateLogInput{
				WallID: wall.ID,
				UserID: m.UserID,
				Text:   m.Text,
			},
		)
	}

	if err := c.repo.CreateManyLogs(ctx, logs); err != nil {
		return fmt.Errorf("cannot create chat wall logs in db: %s", err)
	}

	_, err = c.redis.Pipelined(
		ctx, func(p redis.Pipeliner) error {
			if err := p.SAdd(
				ctx,
				fmt.Sprintf(redis_keys.NukeRedisPrefix, input.ChannelID),
				newHandledMessagesIds,
			).Err(); err != nil {
				return err
			}
			if err := p.Expire(
				ctx,
				fmt.Sprintf(redis_keys.NukeRedisPrefix, input.ChannelID),
				20*time.Minute,
			).Err(); err != nil {
				return err
			}

			return nil
		},
	)
	if err != nil {
		return fmt.Errorf("cannot add handled messages to redis: %s", err)
	}

	return nil
}

type StopInput struct {
	ChannelID string
	Phrase    string
}

var ErrChatWallNotFound = errors.New("chat wall not found")

func (c *Service) Stop(ctx context.Context, input StopInput) error {
	enabled := true

	walls, err := c.repo.GetMany(
		ctx,
		chatwallrepository.GetManyInput{
			ChannelID: input.ChannelID,
			Enabled:   &enabled,
		},
	)
	if err != nil {
		return fmt.Errorf("cannot get chat walls: %s", err)
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
				return fmt.Errorf("cannot update chat wall: %s", err)
			}

			c.chatWallCache.Invalidate(ctx, input.ChannelID)
			return nil
		}
	}

	return ErrChatWallNotFound
}
