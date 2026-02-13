package giveaways

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand"
	"time"

	"github.com/avito-tech/go-transaction-manager/trm/v2"
	"github.com/google/uuid"
	"github.com/guregu/null"
	"github.com/redis/go-redis/v9"
	"github.com/samber/lo"
	buscore "github.com/twirapp/twir/libs/bus-core"
	giveawaysbusmodel "github.com/twirapp/twir/libs/bus-core/giveaways"
	generic_cacher "github.com/twirapp/twir/libs/cache/generic-cacher"
	twitchcache "github.com/twirapp/twir/libs/cache/twitch"
	channels_giveaways "github.com/twirapp/twir/libs/entities/channels_giveaways"
	"github.com/twirapp/twir/libs/logger"
	"github.com/twirapp/twir/libs/repositories/giveaways"
	"github.com/twirapp/twir/libs/repositories/giveaways_participants"
	"github.com/twirapp/twir/libs/repositories/giveaways_participants/model"
	"github.com/twirapp/twir/libs/repositories/users"
	usersstats "github.com/twirapp/twir/libs/repositories/users_stats"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In
	LC fx.Lifecycle

	TxManager                       trm.Manager
	GiveawaysRepository             giveaways.Repository
	GiveawaysParticipantsRepository giveaways_participants.Repository
	GiveawaysCacher                 *generic_cacher.GenericCacher[[]channels_giveaways.Giveaway]
	UsersRepository                 users.Repository
	UsersStatsRepository            usersstats.Repository
	TwitchCache                     *twitchcache.CachedTwitchClient
	Logger                          *slog.Logger
	Redis                           *redis.Client
	TwirBus                         *buscore.Bus
}

func New(opts Opts) *Service {
	s := &Service{
		txManager:                       opts.TxManager,
		giveawaysRepository:             opts.GiveawaysRepository,
		giveawaysParticipantsRepository: opts.GiveawaysParticipantsRepository,
		giveawaysCacher:                 opts.GiveawaysCacher,
		usersRepository:                 opts.UsersRepository,
		usersStatsRepository:            opts.UsersStatsRepository,
		twitchCache:                     opts.TwitchCache,
		logger:                          opts.Logger,
		redis:                           opts.Redis,
		twirBus:                         opts.TwirBus,
	}

	opts.LC.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				return s.twirBus.Giveaways.ChooseWinner.SubscribeGroup(
					"giveaways",
					s.chooseWinner,
				)
			},
			OnStop: func(ctx context.Context) error {
				s.twirBus.Giveaways.ChooseWinner.Unsubscribe()

				return nil
			},
		},
	)

	return s
}

type Service struct {
	txManager                       trm.Manager
	giveawaysRepository             giveaways.Repository
	giveawaysParticipantsRepository giveaways_participants.Repository
	giveawaysCacher                 *generic_cacher.GenericCacher[[]channels_giveaways.Giveaway]
	usersRepository                 users.Repository
	usersStatsRepository            usersstats.Repository
	twitchCache                     *twitchcache.CachedTwitchClient
	logger                          *slog.Logger
	redis                           *redis.Client
	twirBus                         *buscore.Bus
}

const redisParticipantKey = "giveaways:%s:participants:%s"

func (c *Service) TryAddParticipant(
	ctx context.Context,
	userID string,
	userLogin string,
	userDisplayName string,
	giveawayID string,
) error {
	cacheKey := fmt.Sprintf(redisParticipantKey, giveawayID, userID)
	exists, _ := c.redis.Exists(ctx, cacheKey).Result()
	if exists >= 1 {
		return nil
	}

	// Get giveaway to check filters
	parsedGiveawayID, err := uuid.Parse(giveawayID)
	if err != nil {
		return err
	}

	giveaway, err := c.giveawaysRepository.GetByID(ctx, parsedGiveawayID)
	if err != nil {
		return err
	}
	if giveaway.IsNil() {
		return fmt.Errorf("giveaway not found")
	}

	// For KEYWORD giveaways, check filters before adding participant
	if giveaway.Type == channels_giveaways.GiveawayTypeKeyword {
		// Check user stats
		userStats, err := c.usersStatsRepository.GetByUserAndChannelID(ctx, userID, giveaway.ChannelID)
		if err != nil && userStats == nil {
			// If user stats don't exist, user doesn't meet any filter requirements
			return nil
		}

		// Apply filters
		if userStats != nil {
			if giveaway.MinWatchedTime != nil && userStats.Watched < *giveaway.MinWatchedTime {
				return nil
			}
			if giveaway.MinMessages != nil && userStats.Messages < *giveaway.MinMessages {
				return nil
			}
			if giveaway.MinUsedChannelPoints != nil && userStats.UsedChannelPoints < *giveaway.MinUsedChannelPoints {
				return nil
			}
			if giveaway.RequireSubscription && !userStats.IsSubscriber {
				return nil
			}
		}

		// Check follow duration if required
		if giveaway.MinFollowDuration != nil {
			followDuration, err := c.twitchCache.GetUserFollowDuration(ctx, userID, giveaway.ChannelID)
			if err != nil {
				c.logger.Error("cannot get user follow duration", logger.Error(err))
				return nil
			}
			if followDuration == nil || followDuration.Seconds() < float64(*giveaway.MinFollowDuration) {
				return nil
			}
		}
	}

	_, err = c.giveawaysParticipantsRepository.Create(
		ctx,
		giveaways_participants.CreateInput{
			GiveawayID:      giveawayID,
			UserID:          userID,
			UserLogin:       userLogin,
			UserDisplayName: userDisplayName,
		},
	)
	if err != nil {
		return err
	}

	if err := c.twirBus.Giveaways.NewParticipants.Publish(
		ctx,
		giveawaysbusmodel.NewParticipant{
			GiveawayID:      giveawayID,
			UserID:          userID,
			UserLogin:       userLogin,
			UserDisplayName: userDisplayName,
		},
	); err != nil {
		c.logger.Error("cannot publish new participant", logger.Error(err))
	}

	if err := c.redis.Set(ctx, cacheKey, 1, 1*time.Hour).Err(); err != nil {
		c.logger.Error("cannot set giveaway participant to redis cache", logger.Error(err))
	}

	return nil
}

func (c *Service) chooseWinner(
	ctx context.Context,
	req giveawaysbusmodel.ChooseWinnerRequest,
) (giveawaysbusmodel.ChooseWinnerResponse, error) {
	parsedGiveawayId, err := uuid.Parse(req.GiveawayID)
	if err != nil {
		return giveawaysbusmodel.ChooseWinnerResponse{}, err
	}

	giveaway, err := c.giveawaysRepository.GetByID(ctx, parsedGiveawayId)
	if err != nil {
		return giveawaysbusmodel.ChooseWinnerResponse{}, err
	}
	if giveaway.IsNil() {
		return giveawaysbusmodel.ChooseWinnerResponse{}, fmt.Errorf("giveaway not found")
	}

	var eligibleParticipants []model.ChannelGiveawayParticipant

	switch giveaway.Type {
	case channels_giveaways.GiveawayTypeKeyword:
		// For KEYWORD giveaways, get existing participants (filters already applied at join time)
		participants, err := c.giveawaysParticipantsRepository.GetManyByGiveawayID(
			ctx,
			req.GiveawayID,
			giveaways_participants.GetManyInput{
				IgnoreWinners: true,
			},
		)
		if err != nil {
			c.logger.Error("cannot get participants", logger.Error(err))
			return giveawaysbusmodel.ChooseWinnerResponse{}, err
		}
		eligibleParticipants = participants

	case channels_giveaways.GiveawayTypeOnlineChatters:
		// For ONLINE_CHATTERS, fetch online users with filters applied via SQL
		onlineUsers, err := c.usersRepository.GetOnlineUsersWithFilters(
			ctx,
			users.GetOnlineUsersWithFiltersInput{
				ChannelID:            giveaway.ChannelID,
				MinWatchedTime:       giveaway.MinWatchedTime,
				MinMessages:          giveaway.MinMessages,
				MinUsedChannelPoints: giveaway.MinUsedChannelPoints,
				RequireSubscription:  giveaway.RequireSubscription,
			},
		)
		if err != nil {
			c.logger.Error("cannot get online users with filters", logger.Error(err))
			return giveawaysbusmodel.ChooseWinnerResponse{}, err
		}

		// Filter by follow duration if required
		for _, user := range onlineUsers {
			if giveaway.MinFollowDuration != nil {
				followDuration, err := c.twitchCache.GetUserFollowDuration(ctx, user.UserID, giveaway.ChannelID)
				if err != nil {
					c.logger.Error("cannot get user follow duration", logger.Error(err))
					continue
				}
				if followDuration == nil || followDuration.Seconds() < float64(*giveaway.MinFollowDuration) {
					continue
				}
			}

			// Convert to participant model
			eligibleParticipants = append(eligibleParticipants, model.ChannelGiveawayParticipant{
				UserID:      user.UserID,
				UserLogin:   user.UserName, // UserName is actually the login
				DisplayName: user.UserName,
				GiveawayID:  parsedGiveawayId,
				IsWinner:    false,
			})
		}

	default:
		return giveawaysbusmodel.ChooseWinnerResponse{}, fmt.Errorf("unknown giveaway type: %s", giveaway.Type)
	}

	if len(eligibleParticipants) == 0 {
		return giveawaysbusmodel.ChooseWinnerResponse{}, fmt.Errorf("Cannot do roll with empty participants")
	}

	winners := make([]model.ChannelGiveawayParticipant, 0, 1)
	err = c.txManager.Do(
		ctx,
		func(txCtx context.Context) error {
			winnerInd := rand.Intn(len(eligibleParticipants))
			winnerData := eligibleParticipants[winnerInd]

			var winner model.ChannelGiveawayParticipant

			// For ONLINE_CHATTERS, we need to create the participant record
			if giveaway.Type == channels_giveaways.GiveawayTypeOnlineChatters {
				winner, err = c.giveawaysParticipantsRepository.Create(
					txCtx,
					giveaways_participants.CreateInput{
						GiveawayID:      req.GiveawayID,
						UserID:          winnerData.UserID,
						UserLogin:       winnerData.UserLogin,
						UserDisplayName: winnerData.DisplayName,
					},
				)
				if err != nil {
					c.logger.Error("create winner error", logger.Error(err))
					return err
				}

				// Mark as winner
				winner, err = c.giveawaysParticipantsRepository.Update(
					txCtx,
					winner.ID.String(),
					giveaways_participants.UpdateInput{
						IsWinner: lo.ToPtr(true),
					},
				)
				if err != nil {
					c.logger.Error("update winner error", logger.Error(err))
					return err
				}
			} else {
				// For KEYWORD giveaways, update existing participant
				winner, err = c.giveawaysParticipantsRepository.Update(
					txCtx,
					winnerData.ID.String(),
					giveaways_participants.UpdateInput{
						IsWinner: lo.ToPtr(true),
					},
				)
				if err != nil {
					c.logger.Error("update winner error", logger.Error(err))
					return err
				}
			}

			winners = append(winners, winner)

			_, err = c.giveawaysRepository.UpdateStatuses(
				txCtx,
				winner.GiveawayID,
				giveaways.UpdateStatusInput{
					StoppedAt: null.NewTime(time.Now(), true),
				},
			)
			if err != nil {
				c.logger.Error("update error", logger.Error(err))
				return err
			}

			return nil
		},
	)
	if err != nil {
		c.logger.Error("tx error", logger.Error(err))
		return giveawaysbusmodel.ChooseWinnerResponse{}, err
	}

	if err := c.giveawaysCacher.Invalidate(ctx, giveaway.ChannelID); err != nil {
		c.logger.Error("cannot invalidate giveaways cache", logger.Error(err))
	}

	mappedWinners := make([]giveawaysbusmodel.Winner, 0, len(winners))
	for _, winner := range winners {
		mappedWinners = append(
			mappedWinners,
			giveawaysbusmodel.Winner{
				UserID:          winner.UserID,
				UserLogin:       winner.UserLogin,
				UserDisplayName: winner.DisplayName,
			},
		)
	}

	return giveawaysbusmodel.ChooseWinnerResponse{
		Winners: mappedWinners,
	}, nil
}
