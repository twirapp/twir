package buslistener

import (
	"context"
	"fmt"

	buscore "github.com/twirapp/twir/libs/bus-core"
	botssettings "github.com/twirapp/twir/libs/bus-core/bots-settings"
	"github.com/twirapp/twir/libs/cache"
	"github.com/twirapp/twir/libs/logger"
	channelscommandsprefixmodel "github.com/twirapp/twir/libs/repositories/channels_commands_prefix/model"
	"go.uber.org/fx"
)

type BusListener struct {
	bus         *buscore.Bus
	logger      logger.Logger
	prefixCache cache.Cache[channelscommandsprefixmodel.ChannelsCommandsPrefix]
}

type Opts struct {
	fx.In
	Lc fx.Lifecycle

	Bus         *buscore.Bus
	Logger      logger.Logger
	PrefixCache cache.Cache[channelscommandsprefixmodel.ChannelsCommandsPrefix]
}

func New(opts Opts) (*BusListener, error) {
	bl := &BusListener{
		bus:         opts.Bus,
		logger:      opts.Logger,
		prefixCache: opts.PrefixCache,
	}

	opts.Lc.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				if err := bl.bus.BotsSettings.UpdatePrefix.Subscribe(
					bl.updateChannelCommandPrefix,
				); err != nil {
					return err
				}

				return nil
			},
			OnStop: func(ctx context.Context) error {
				bl.bus.BotsSettings.UpdatePrefix.Unsubscribe()
				return nil
			},
		},
	)

	return bl, nil
}

func (bl *BusListener) updateChannelCommandPrefix(
	ctx context.Context,
	data botssettings.UpdatePrefixRequest,
) (struct{}, error) {
	prefix := channelscommandsprefixmodel.ChannelsCommandsPrefix{
		ID:        data.ID,
		ChannelID: data.ChannelID,
		Prefix:    data.Prefix,
		CreatedAt: data.CreatedAt,
		UpdatedAt: data.UpdatedAt,
	}

	if err := bl.prefixCache.Set(ctx, data.ChannelID, prefix); err != nil {
		return struct{}{}, fmt.Errorf("set channel command prefix cache: %s", err)
	}

	return struct{}{}, nil
}
