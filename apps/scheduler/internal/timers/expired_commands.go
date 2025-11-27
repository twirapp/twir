package timers

import (
	"context"
	"log/slog"
	"time"

	"github.com/guregu/null"
	buscore "github.com/twirapp/twir/libs/bus-core"
	commandscache "github.com/twirapp/twir/libs/cache/commands"
	generic_cacher "github.com/twirapp/twir/libs/cache/generic-cacher"
	config "github.com/twirapp/twir/libs/config"
	model "github.com/twirapp/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/logger"
	commandswithgroupsandresponsesrepository "github.com/twirapp/twir/libs/repositories/commands_with_groups_and_responses"
	commandswithgroupsandresponsesmodel "github.com/twirapp/twir/libs/repositories/commands_with_groups_and_responses/model"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type ExpiredCommandsOpts struct {
	fx.In
	Lc fx.Lifecycle

	Logger *slog.Logger
	Config config.Config

	Gorm         *gorm.DB
	TwirBus      *buscore.Bus
	CommandsRepo commandswithgroupsandresponsesrepository.Repository
}

type expiredCommands struct {
	config        config.Config
	logger        *slog.Logger
	db            *gorm.DB
	commandsCache *generic_cacher.GenericCacher[[]commandswithgroupsandresponsesmodel.CommandWithGroupAndResponses]
}

func NewExpiredCommands(opts ExpiredCommandsOpts) *expiredCommands {
	timeTick := 15 * time.Second
	if opts.Config.AppEnv == "production" {
		timeTick = 5 * time.Minute
	}
	ticker := time.NewTicker(timeTick)

	ctx, cancel := context.WithCancel(context.Background())

	s := &expiredCommands{
		config:        opts.Config,
		logger:        opts.Logger,
		db:            opts.Gorm,
		commandsCache: commandscache.New(opts.CommandsRepo, opts.TwirBus),
	}

	opts.Lc.Append(
		fx.Hook{
			OnStart: func(_ context.Context) error {
				go func() {
					for {
						select {
						case <-ctx.Done():
							ticker.Stop()
							return
						case <-ticker.C:
							s.checkForExpiredCommands(ctx)
						}
					}
				}()

				return nil
			},
			OnStop: func(_ context.Context) error {
				cancel()
				return nil
			},
		},
	)

	return s
}

func (s *expiredCommands) checkForExpiredCommands(ctx context.Context) {
	var commands []model.ChannelsCommands
	if err := s.db.WithContext(ctx).Where(
		`"expires_at" < ?`,
		time.Now().UTC(),
	).Find(&commands).Error; err != nil {
		s.logger.Error("failed to get commands", logger.Error(err))
		return
	}

	for _, c := range commands {
		s.logger.Info("Command expired", slog.Any("command", c))

		if !c.ExpiresAt.Valid || c.ExpiresType == nil {
			continue
		}

		if *c.ExpiresType == model.ChannelCommandExpiresTypeDisable && c.Enabled {
			c.Enabled = false
			c.ExpiresType = nil
			c.ExpiresAt = null.Time{}
			s.db.WithContext(ctx).Save(&c)
		} else if *c.ExpiresType == model.ChannelCommandExpiresTypeDelete && !c.Default {
			err := s.db.WithContext(ctx).Delete(
				&c,
			).Error
			if err != nil {
				s.logger.Error("failed to delete command", logger.Error(err))
			}
		}

		err := s.commandsCache.Invalidate(
			ctx,
			c.ChannelID,
		)
		if err != nil {
			s.logger.Error("failed to invalidate commands cache", logger.Error(err))
			return
		}
	}
}
