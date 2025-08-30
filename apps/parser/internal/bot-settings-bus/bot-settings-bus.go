package botsettingsbus

import (
	"context"
	"fmt"

	"github.com/twirapp/twir/apps/parser/internal/types/services"
	buscore "github.com/twirapp/twir/libs/bus-core"
	botssettings "github.com/twirapp/twir/libs/bus-core/bots-settings"
	channelscommandsprefixmodel "github.com/twirapp/twir/libs/repositories/channels_commands_prefix/model"
)

type BotSettingsBus struct {
	bus      *buscore.Bus
	services *services.Services
}

func New(
	bus *buscore.Bus,
	services *services.Services,
) BotSettingsBus {
	return BotSettingsBus{
		bus:      bus,
		services: services,
	}
}

func (bss BotSettingsBus) Subscribe() error {
	if err := bss.bus.BotsSettings.UpdatePrefix.Subscribe(
		bss.updateChannelCommandPrefix,
	); err != nil {
		return err
	}

	return nil
}

func (bss BotSettingsBus) Unsubscribe() {
	bss.bus.BotsSettings.UpdatePrefix.Unsubscribe()
}

func (bss BotSettingsBus) updateChannelCommandPrefix(
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

	if err := bss.services.CommandsPrefixCache.Set(ctx, data.ChannelID, prefix); err != nil {
		return struct{}{}, fmt.Errorf("set channel command prefix cache: %s", err)
	}

	return struct{}{}, nil
}
