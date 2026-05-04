package giveaways

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/guregu/null"
	buscore "github.com/twirapp/twir/libs/bus-core"
	botsbus "github.com/twirapp/twir/libs/bus-core/bots"
	giveawaysbus "github.com/twirapp/twir/libs/bus-core/giveaways"
	generic_cacher "github.com/twirapp/twir/libs/cache/generic-cacher"
	twitchcache "github.com/twirapp/twir/libs/cache/twitch"
	channels_giveaways "github.com/twirapp/twir/libs/entities/channels_giveaways"
	"github.com/twirapp/twir/libs/errors"
	"github.com/twirapp/twir/libs/logger"
	channelsrepository "github.com/twirapp/twir/libs/repositories/channels"
	"github.com/twirapp/twir/libs/repositories/channels_giveaways_settings"
	"github.com/twirapp/twir/libs/repositories/giveaways"
	"github.com/twirapp/twir/libs/repositories/giveaways_participants"
	giveawaysparticipantsmodel "github.com/twirapp/twir/libs/repositories/giveaways_participants/model"
	"github.com/twirapp/twir/libs/wsrouter"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In
	LC fx.Lifecycle

	GiveawaysRepository             giveaways.Repository
	GiveawaysParticipantsRepository giveaways_participants.Repository
	GiveawaysSettingsRepository     channels_giveaways_settings.Repository
	ChannelsRepository              channelsrepository.Repository
	GiveawaysCacher                 *generic_cacher.GenericCacher[[]channels_giveaways.Giveaway]
	TwirBus                         *buscore.Bus
	Logger                          *slog.Logger
	WsRouter                        wsrouter.WsRouter
	TwitchCache                     *twitchcache.CachedTwitchClient
}

func New(opts Opts) *Service {
	s := &Service{
		giveawaysRepository:             opts.GiveawaysRepository,
		giveawaysParticipantsRepository: opts.GiveawaysParticipantsRepository,
		giveawaysSettingsRepository:     opts.GiveawaysSettingsRepository,
		channelsRepository:              opts.ChannelsRepository,
		giveawaysCacher:                 opts.GiveawaysCacher,
		twirBus:                         opts.TwirBus,
		logger:                          opts.Logger,
		wsRouter:                        opts.WsRouter,
		twitchCache:                     opts.TwitchCache,
	}

	opts.LC.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				return s.twirBus.Giveaways.NewParticipants.SubscribeGroup(
					"api-gql",
					func(ctx context.Context, data giveawaysbus.NewParticipant) (struct{}, error) {
						err := s.handleNewParticipants(ctx, data)

						return struct{}{}, err
					},
				)
			},
			OnStop: func(ctx context.Context) error {
				s.twirBus.Giveaways.NewParticipants.Unsubscribe()
				return nil
			},
		},
	)

	return s
}

type Service struct {
	giveawaysRepository             giveaways.Repository
	giveawaysParticipantsRepository giveaways_participants.Repository
	giveawaysSettingsRepository     channels_giveaways_settings.Repository
	channelsRepository              channelsrepository.Repository
	giveawaysCacher                 *generic_cacher.GenericCacher[[]channels_giveaways.Giveaway]
	twirBus                         *buscore.Bus
	logger                          *slog.Logger
	wsRouter                        wsrouter.WsRouter
	twitchCache                     *twitchcache.CachedTwitchClient
}

type CreateInput struct {
	ChannelID            string
	Type                 channels_giveaways.GiveawayType
	Keyword              *string
	MinWatchedTime       *int64
	MinMessages          *int32
	MinUsedChannelPoints *int64
	MinFollowDuration    *int64
	RequireSubscription  bool
	CreatedByUserID      string
}

func (c *Service) handleNewParticipants(
	ctx context.Context,
	participant giveawaysbus.NewParticipant,
) error {
	if err := c.wsRouter.Publish(
		CreateNewParticipantSubscriptionKeyByGiveawayID(participant.GiveawayID),
		participant,
	); err != nil {
		c.logger.Error("cannot publish new participant", logger.Error(err))
		return errors.NewInternalError("Failed to publish new participant", err)
	}

	return nil
}

func (c *Service) Start(
	ctx context.Context,
	giveawayID uuid.UUID,
	channelID string,
) (channels_giveaways.Giveaway, error) {
	dbGiveaway, err := c.GiveawayGet(ctx, giveawayID, channelID)
	if err != nil {
		return channels_giveaways.GiveawayNil, err
	}

	if dbGiveaway == channels_giveaways.GiveawayNil {
		return channels_giveaways.GiveawayNil, errors.NewNotFoundError("Giveaway with this ID was not found")
	}

	if dbGiveaway.StoppedAt != nil {
		dbGiveaway, err = c.GiveawayUpdateStatus(
			ctx,
			giveawayID,
			channelID,
			UpdateStatusInput{
				StartedAt: null.NewTime(time.Now(), true),
				StoppedAt: null.Time{},
			},
		)
		if err != nil {
			return channels_giveaways.GiveawayNil, err
		}

		return dbGiveaway, nil
	}

	if dbGiveaway.StartedAt != nil {
		return channels_giveaways.GiveawayNil, errors.NewBadRequestError(
			"This giveaway has already been started and is not stopped yet",
		)
	}

	if dbGiveaway.StartedAt == nil {
		dbGiveaway, err = c.GiveawayUpdateStatus(
			ctx,
			giveawayID,
			channelID,
			UpdateStatusInput{
				StartedAt: null.NewTime(time.Now(), true),
			},
		)
		if err != nil {
			return channels_giveaways.GiveawayNil, err
		}

		return dbGiveaway, nil
	}

	return dbGiveaway, nil
}

func (c *Service) Stop(
	ctx context.Context,
	giveawayID uuid.UUID,
	channelID string,
) (channels_giveaways.Giveaway, error) {
	dbGiveaway, err := c.GiveawayGet(ctx, giveawayID, channelID)
	if err != nil {
		return channels_giveaways.GiveawayNil, err
	}

	if dbGiveaway == channels_giveaways.GiveawayNil {
		return channels_giveaways.GiveawayNil, errors.NewNotFoundError("Giveaway with this ID was not found")
	}

	if dbGiveaway.StartedAt == nil {
		return channels_giveaways.GiveawayNil, errors.NewBadRequestError(
			"Cannot stop a giveaway that has not been started yet",
		)
	}

	if dbGiveaway.StoppedAt != nil {
		return channels_giveaways.GiveawayNil, errors.NewBadRequestError("This giveaway has already been stopped")
	}

	dbGiveaway, err = c.GiveawayUpdateStatus(
		ctx,
		giveawayID,
		channelID,
		UpdateStatusInput{
			StoppedAt: null.NewTime(time.Now(), true),
		},
	)
	if err != nil {
		return channels_giveaways.GiveawayNil, err
	}

	return dbGiveaway, nil
}

func (c *Service) ChooseWinners(
	ctx context.Context,
	giveawayID uuid.UUID,
	channelID string,
) ([]channels_giveaways.GiveawayWinner, error) {
	dbGiveaway, err := c.GiveawayGet(ctx, giveawayID, channelID)
	if err != nil {
		return nil, err
	}

	if dbGiveaway == channels_giveaways.GiveawayNil {
		return nil, fmt.Errorf("Giveaway doesnt exists")
	}

	winners, err := c.twirBus.Giveaways.ChooseWinner.Request(
		ctx, giveawaysbus.ChooseWinnerRequest{
			GiveawayID: giveawayID.String(),
		},
	)
	if err != nil {
		c.logger.Error("Cannot choose winners", logger.Error(err))
		return nil, err
	}

	if len(winners.Data.Winners) == 0 {
		return nil, fmt.Errorf("Cannot choose winner, probably there is no more participants?")
	}

	var result []channels_giveaways.GiveawayWinner
	for _, winner := range winners.Data.Winners {
		result = append(result, c.giveawayWinnerBusModelToEntity(winner))
	}

	err = c.updateGiveawaysCacheForChannel(ctx, channelID)
	if err != nil {
		c.logger.Error("Cannot update winners cache", logger.Error(err))
		return nil, err
	}

	go func() {
		if err := c.sendWinnerMessage(context.Background(), channelID, result); err != nil {
			c.logger.Error("Cannot send winner message", logger.Error(err))
		}
	}()

	return result, nil
}

func (c *Service) Create(ctx context.Context, input CreateInput) (channels_giveaways.Giveaway, error) {
	parsedChannelID, err := uuid.Parse(input.ChannelID)
	if err != nil {
		return channels_giveaways.GiveawayNil, errors.NewInternalError("Failed to parse channel id", err)
	}

	parsedCreatedByUserID, err := uuid.Parse(input.CreatedByUserID)
	if err != nil {
		return channels_giveaways.GiveawayNil, errors.NewInternalError("Failed to parse created by user id", err)
	}

	// Validate that keyword is present for KEYWORD type
	if input.Type == channels_giveaways.GiveawayTypeKeyword && (input.Keyword == nil || *input.Keyword == "") {
		return channels_giveaways.GiveawayNil, errors.NewBadRequestError(
			"Keyword is required when creating a KEYWORD type giveaway",
		)
	}

	// Check for duplicate keyword giveaway
	if input.Type == channels_giveaways.GiveawayTypeKeyword && input.Keyword != nil {
		dbGiveaway, err := c.giveawaysRepository.GetByChannelIDAndKeyword(
			ctx,
			parsedChannelID,
			*input.Keyword,
		)
		if err != nil && err != giveaways.ErrNotFound {
			return channels_giveaways.GiveawayNil, errors.NewInternalError(
				"Failed to check for duplicate giveaway",
				err,
			)
		}

		if !dbGiveaway.IsNil() {
			return channels_giveaways.GiveawayNil, errors.NewConflictError(
				"A giveaway with this keyword already exists and runed on your channel",
			)
		}
	}

	giveaway, err := c.giveawaysRepository.Create(
		ctx,
		giveaways.CreateInput{
			ChannelID:            parsedChannelID,
			Type:                 channels_giveaways.GiveawayType(input.Type),
			Keyword:              input.Keyword,
			MinWatchedTime:       input.MinWatchedTime,
			MinMessages:          input.MinMessages,
			MinUsedChannelPoints: input.MinUsedChannelPoints,
			MinFollowDuration:    input.MinFollowDuration,
			RequireSubscription:  input.RequireSubscription,
			CreatedByUserID:      parsedCreatedByUserID,
		},
	)
	if err != nil {
		return channels_giveaways.GiveawayNil, errors.NewInternalError("Failed to create giveaway", err)
	}

	err = c.updateGiveawaysCacheForChannel(ctx, input.ChannelID)
	if err != nil {
		return channels_giveaways.GiveawayNil, err
	}

	return giveaway, nil
}

type GetParticipantsInput struct {
	OnlyWinners bool
}

func (c *Service) GetParticipantsForGiveaway(
	ctx context.Context,
	giveawayID uuid.UUID,
	input GetParticipantsInput,
) ([]channels_giveaways.GiveawayParticipant, error) {
	participants, err := c.giveawaysParticipantsRepository.GetManyByGiveawayID(
		ctx,
		giveawayID,
		giveaways_participants.GetManyInput{
			OnlyWinners: input.OnlyWinners,
		},
	)
	if err != nil {
		return nil, err
	}

	mappedParticipants := make([]channels_giveaways.GiveawayParticipant, 0, len(participants))
	for _, participant := range participants {
		mappedParticipants = append(mappedParticipants, c.giveawayParticipantModelToEntity(participant))
	}

	return mappedParticipants, nil
}

func (c *Service) GiveawayGet(
	ctx context.Context,
	giveawayID uuid.UUID,
	channelID string,
) (channels_giveaways.Giveaway, error) {
	parsedChannelID, err := uuid.Parse(channelID)
	if err != nil {
		return channels_giveaways.GiveawayNil, err
	}

	giveaway, err := c.giveawaysRepository.GetByID(ctx, giveawayID)
	if err != nil && err != giveaways.ErrNotFound {
		return channels_giveaways.GiveawayNil, errors.NewInternalError("Failed to fetch giveaway", err)
	}

	if giveaway.IsNil() {
		return channels_giveaways.GiveawayNil, nil
	}

	if giveaway.ChannelID != parsedChannelID {
		return channels_giveaways.GiveawayNil, nil
	}

	err = c.updateGiveawaysCacheForChannel(ctx, channelID)
	if err != nil {
		return channels_giveaways.GiveawayNil, err
	}

	return giveaway, nil
}

func (c *Service) GiveawaysGetMany(
	ctx context.Context,
	channelID string,
) ([]channels_giveaways.Giveaway, error) {
	parsedChannelID, err := uuid.Parse(channelID)
	if err != nil {
		return nil, err
	}

	dbGiveaways, err := c.giveawaysRepository.GetManyByChannelID(ctx, parsedChannelID)
	if err != nil {
		return nil, err
	}

	mappedGiveaways := make([]channels_giveaways.Giveaway, 0, len(dbGiveaways))
	for _, giveaway := range dbGiveaways {
		mappedGiveaways = append(mappedGiveaways, giveaway)
	}

	return mappedGiveaways, nil
}

func (c *Service) GiveawayRemove(
	ctx context.Context,
	giveawayID uuid.UUID,
	channelID string,
) (channels_giveaways.Giveaway, error) {
	dbGiveaway, err := c.GiveawayGet(ctx, giveawayID, channelID)
	if err != nil {
		return channels_giveaways.GiveawayNil, err
	}
	if dbGiveaway.IsNil() {
		return channels_giveaways.GiveawayNil, errors.NewNotFoundError("Giveaway with this ID was not found")
	}

	err = c.giveawaysRepository.Delete(ctx, giveawayID)
	if err != nil {
		return channels_giveaways.GiveawayNil, err
	}

	err = c.updateGiveawaysCacheForChannel(ctx, channelID)

	return dbGiveaway, err
}

type UpdateInput struct {
	StartedAt            *time.Time
	Keyword              *string
	StoppedAt            *time.Time
	MinWatchedTime       *int64
	MinMessages          *int32
	MinUsedChannelPoints *int64
	MinFollowDuration    *int64
	RequireSubscription  *bool
}

func (c *Service) GiveawayUpdate(
	ctx context.Context,
	giveawayID uuid.UUID,
	channelID string,
	input UpdateInput,
) (channels_giveaways.Giveaway, error) {
	currentGiveaway, err := c.GiveawayGet(ctx, giveawayID, channelID)
	if err != nil {
		return channels_giveaways.GiveawayNil, err
	}
	if currentGiveaway.IsNil() {
		return channels_giveaways.GiveawayNil, errors.NewNotFoundError("Giveaway with this ID was not found")
	}

	dbGiveaway, err := c.giveawaysRepository.Update(
		ctx, giveawayID, giveaways.UpdateInput{
			StartedAt:            input.StartedAt,
			Keyword:              input.Keyword,
			StoppedAt:            input.StoppedAt,
			MinWatchedTime:       input.MinWatchedTime,
			MinMessages:          input.MinMessages,
			MinUsedChannelPoints: input.MinUsedChannelPoints,
			MinFollowDuration:    input.MinFollowDuration,
			RequireSubscription:  input.RequireSubscription,
		},
	)
	if err != nil {
		return channels_giveaways.GiveawayNil, err
	}

	err = c.updateGiveawaysCacheForChannel(ctx, channelID)
	if err != nil {
		return channels_giveaways.GiveawayNil, err
	}

	return dbGiveaway, nil
}

type UpdateStatusInput struct {
	StartedAt null.Time
	StoppedAt null.Time
}

func (c *Service) GiveawayUpdateStatus(
	ctx context.Context,
	giveawayID uuid.UUID,
	channelID string,
	input UpdateStatusInput,
) (channels_giveaways.Giveaway, error) {
	currentGiveaway, err := c.GiveawayGet(ctx, giveawayID, channelID)
	if err != nil {
		return channels_giveaways.GiveawayNil, err
	}
	if currentGiveaway.IsNil() {
		return channels_giveaways.GiveawayNil, errors.NewNotFoundError("Giveaway with this ID was not found")
	}

	dbGiveaway, err := c.giveawaysRepository.UpdateStatuses(
		ctx,
		giveawayID,
		giveaways.UpdateStatusInput{
			StartedAt: input.StartedAt,
			StoppedAt: input.StoppedAt,
		},
	)
	if err != nil {
		return channels_giveaways.GiveawayNil, err
	}

	err = c.updateGiveawaysCacheForChannel(ctx, channelID)
	if err != nil {
		return channels_giveaways.GiveawayNil, err
	}

	return dbGiveaway, nil
}

func (c *Service) giveawayParticipantModelToEntity(
	m giveawaysparticipantsmodel.ChannelGiveawayParticipant,
) channels_giveaways.GiveawayParticipant {
	return channels_giveaways.GiveawayParticipant{
		UserLogin:   m.UserLogin,
		DisplayName: m.DisplayName,
		UserID:      m.UserID,
		IsWinner:    m.IsWinner,
		ID:          m.ID,
		GiveawayID:  m.GiveawayID,
	}
}

func (c *Service) giveawayWinnerBusModelToEntity(
	m giveawaysbus.Winner,
) channels_giveaways.GiveawayWinner {
	return channels_giveaways.GiveawayWinner{
		UserID:      uuid.MustParse(m.UserID),
		UserLogin:   m.UserLogin,
		DisplayName: m.UserDisplayName,
	}
}

/*
We need to update value of array of channels giveaways in Kv,
cus we search for keyword for every message and don't wanna use database calls in message handlers,
so we are do probably some unnecessarily work here but provide better consistency.
Also, limits for max giveaways per channel is low, so it will be fast, I suppose.
*/
func (c *Service) updateGiveawaysCacheForChannel(ctx context.Context, channelID string) error {
	parsedChannelID, err := uuid.Parse(channelID)
	if err != nil {
		return err
	}

	dbGiveaways, err := c.giveawaysRepository.GetManyActiveByChannelID(ctx, parsedChannelID)
	if err != nil {
		return err
	}

	err = c.giveawaysCacher.SetValue(ctx, channelID, dbGiveaways)
	if err != nil {
		return err
	}

	return nil
}

func (c *Service) sendWinnerMessage(
	ctx context.Context,
	channelID string,
	winners []channels_giveaways.GiveawayWinner,
) error {
	parsedChannelID, err := uuid.Parse(channelID)
	if err != nil {
		return err
	}

	settings, err := c.giveawaysSettingsRepository.GetByChannelID(ctx, channelID)
	if err != nil {
		c.logger.ErrorContext(ctx, "Cannot get giveaways settings", logger.Error(err))
		return err
	}

	if settings.IsNil() {
		c.logger.WarnContext(ctx, "Giveaways settings not found for channel", "channelId", channelID)
		return nil
	}

	channel, err := c.channelsRepository.GetByID(ctx, parsedChannelID)
	if err != nil {
		return err
	}
	if channel.TwitchPlatformID == nil {
		return fmt.Errorf("channel has no twitch platform id")
	}

	for _, winner := range winners {
		message := strings.ReplaceAll(settings.WinnerMessage, "{winner}", winner.UserLogin)

		err := c.twirBus.Bots.SendMessage.Publish(
			ctx,
			botsbus.SendMessageRequest{
				ChannelId:      *channel.TwitchPlatformID,
				Message:        message,
				SkipRateLimits: true,
			},
		)
		if err != nil {
			c.logger.ErrorContext(
				ctx,
				"Cannot send winner message",
				logger.Error(err),
				"winner",
				winner.UserLogin,
			)
			return err
		}
	}

	return nil
}
