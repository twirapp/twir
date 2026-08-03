package discord_go

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	discordapi "github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/session"
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
	api     *discordapi.Client
	session *session.Session

	logger      *slog.Logger
	discordRepo channelsintegrationsdiscord.Repository
}

func New(opts Opts) (*Discord, error) {
	if opts.Config.DiscordBotToken == "" {
		return &Discord{}, nil
	}

	log := logger.WithComponent(opts.Logger, "discord")
	discordSession := session.NewWithIntents(
		"Bot "+opts.Config.DiscordBotToken,
		gateway.IntentGuilds,
		gateway.IntentGuildMessages,
		gateway.IntentMessageContent,
	)
	d := &Discord{
		logger:      log,
		discordRepo: opts.DiscordRepo,
		api:         discordSession.Client,
		session:     discordSession,
	}
	discordSession.AddHandler(d.handleShardReady)
	discordSession.AddHandler(d.handleGuildDelete)

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

				if err := d.session.Open(ctx); err != nil {
					return fmt.Errorf("open Discord gateway: %w", err)
				}

				return nil
			},
			OnStop: func(ctx context.Context) error {
				closeResult := make(chan error, 1)
				go func() {
					closeResult <- d.session.Close()
				}()

				select {
				case err := <-closeResult:
					if errors.Is(err, session.ErrClosed) {
						return nil
					}
					return err
				case <-ctx.Done():
					return ctx.Err()
				}
			},
		},
	)

	return d, nil
}

func (c *Discord) AddHandler(handler any) func() {
	if c.session == nil {
		return func() {}
	}

	return c.session.AddHandler(handler)
}

func (c *Discord) Messages(
	ctx context.Context,
	channelID discord.ChannelID,
	limit uint,
) ([]discord.Message, error) {
	return c.api.WithContext(ctx).Messages(channelID, limit)
}

func (c *Discord) Message(
	ctx context.Context,
	channelID discord.ChannelID,
	messageID discord.MessageID,
) (*discord.Message, error) {
	return c.api.WithContext(ctx).Message(channelID, messageID)
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
