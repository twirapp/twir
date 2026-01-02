package discord_go

import (
	"context"
	"log/slog"

	discordapi "github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	cfg "github.com/twirapp/twir/libs/config"
	"github.com/twirapp/twir/libs/logger"
	channelsintegrationsdiscord "github.com/twirapp/twir/libs/repositories/channels_integrations_discord"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In

	LC          fx.Lifecycle
	Config      cfg.Config
	Logger      *slog.Logger
	DiscordRepo channelsintegrationsdiscord.Repository
}

type Discord struct {
	api *discordapi.Client

	logger      *slog.Logger
	discordRepo channelsintegrationsdiscord.Repository
}

func New(opts Opts) (*Discord, error) {
	if opts.Config.DiscordBotToken == "" {
		return &Discord{}, nil
	}

	discordapi.NewClient(opts.Config.DiscordBotToken)

	log := logger.WithComponent(opts.Logger, "discord")
	d := &Discord{
		logger:      log,
		discordRepo: opts.DiscordRepo,
		api:         discordapi.NewClient("Bot " + opts.Config.DiscordBotToken),
	}

	opts.LC.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				botInfo, err := d.api.Me()
				if err != nil {
					return err
				}
				log.Info(
					"Starting Discord bot",
					slog.String("bot_name", botInfo.Username),
					slog.String("bot_id", botInfo.ID.String()),
				)

				return nil
			},
			OnStop: func(ctx context.Context) error {
				return nil
			},
		},
	)

	return d, nil
}

func (c *Discord) DeleteMessage(_ context.Context, channelID, messageID, reason string) error {
	channelIdSnowlake, err := discord.ParseSnowflake(channelID)
	if err != nil {
		return err
	}

	messageIdSnowlake, err := discord.ParseSnowflake(messageID)
	if err != nil {
		return err
	}

	return c.api.DeleteMessage(
		discord.ChannelID(channelIdSnowlake),
		discord.MessageID(messageIdSnowlake),
		discordapi.AuditLogReason(reason),
	)
}

type SendMessageResponse struct {
	MessageID string
}

func (c *Discord) SendMessage(_ context.Context, channelID, message string, embeds ...discord.Embed) (SendMessageResponse, error) {
	channelIdSnowlake, err := discord.ParseSnowflake(channelID)
	if err != nil {
		return SendMessageResponse{}, err
	}

	resp, err := c.api.SendMessage(discord.ChannelID(channelIdSnowlake), message, embeds...)
	if err != nil {
		return SendMessageResponse{}, err
	}

	return SendMessageResponse{
		MessageID: resp.ID.String(),
	}, nil
}

func (c *Discord) EditMessage(_ context.Context, channelID, messageID, newMessage string, embeds ...discord.Embed) error {
	channelIdSnowlake, err := discord.ParseSnowflake(channelID)
	if err != nil {
		return err
	}

	messageIdSnowlake, err := discord.ParseSnowflake(messageID)
	if err != nil {
		return err
	}

	_, err = c.api.EditMessage(
		discord.ChannelID(channelIdSnowlake),
		discord.MessageID(messageIdSnowlake),
		newMessage,
		embeds...,
	)
	return err
}
