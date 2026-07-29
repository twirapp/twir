package twitchactions

import (
	"context"
	"fmt"
	"log/slog"
	"slices"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/aidenwallis/go-ratelimiting/redis"
	"github.com/kvizyx/twitchy/helix"
	"github.com/twirapp/twir/libs/bus-core/bots"
	channelentity "github.com/twirapp/twir/libs/entities/channel"
	channelplatformentity "github.com/twirapp/twir/libs/entities/channel_platform"
	"github.com/twirapp/twir/libs/logger"
	"github.com/twirapp/twir/libs/repositories/sentmessages"
	"github.com/twirapp/twir/libs/repositories/toxic_messages"
	"github.com/twirapp/twir/libs/twitch"
)

type SendMessageOpts struct {
	BroadcasterID        string
	SenderID             string
	Message              string
	ReplyParentMessageID string
	IsAnnounce           bool
	SkipToxicityCheck    bool
	SkipRateLimits       bool
	AnnounceColor        bots.AnnounceColor
}

const shoutOutPrefix = "/shoutout"

var allowedSlashCommands = []string{
	"/me",
	"/announce",
	"/announceblue",
	"/announcegreen",
	"/announceorange",
	"/announcepurple",
	"/shoutout",
	"/timeout",
	"/ban",
}

func validateResponseSlashes(response string) string {
	if slices.ContainsFunc(
		allowedSlashCommands, func(s string) bool {
			return strings.HasPrefix(response, s)
		},
	) {
		return response
	} else if strings.HasPrefix(response, "/") {
		return fmt.Sprintf(
			"Slash commands except %s is disallowed. This response wont be ever sended.",
			strings.Join(allowedSlashCommands, ", "),
		)
	} else if strings.HasPrefix(response, ".") {
		return `Message cannot start with "." symbol.`
	} else {
		return response
	}
}

func (c *TwitchActions) SendMessage(
	ctx context.Context,
	binding channelplatformentity.ChannelPlatform,
	opts SendMessageOpts,
) error {
	opts.BroadcasterID = binding.PlatformChannelID

	resp, err := c.rateLimiter.Use(
		ctx,
		&redis.SlidingWindowOptions{
			Key:             fmt.Sprintf("bots:rate_limit:send_message:%s", opts.BroadcasterID),
			MaximumCapacity: 200,
			Window:          30 * time.Second,
		},
	)
	if err != nil {
		return err
	}
	if !resp.Success {
		return nil
	}

	botConfig, active, err := activeTwitchBinding(binding)
	if err != nil {
		return fmt.Errorf("parse Twitch bot config: %w", err)
	}
	if !active {
		return nil
	}

	if opts.SenderID == "" {
		opts.SenderID = botConfig.BotID
	}

	if strings.HasPrefix(opts.Message, "/timeout") || strings.HasPrefix(opts.Message, "/ban") {
		return c.timeoutFromMessage(
			ctx,
			channelentity.Channel{Bindings: []channelplatformentity.ChannelPlatform{binding}},
			opts,
		)
	}

	if strings.HasPrefix(opts.Message, "/announce") && !opts.IsAnnounce {
		opts.IsAnnounce = true

		switch {
		case strings.HasPrefix(opts.Message, "/announceblue"):
			opts.AnnounceColor = bots.AnnounceColorBlue
		case strings.HasPrefix(opts.Message, "/announcegreen"):
			opts.AnnounceColor = bots.AnnounceColorGreen
		case strings.HasPrefix(opts.Message, "/announceorange"):
			opts.AnnounceColor = bots.AnnounceColorOrange
		case strings.HasPrefix(opts.Message, "/announcepurple"):
			opts.AnnounceColor = bots.AnnounceColorPurple
		}
	}

	var twitchClient *helix.Client
	if opts.IsAnnounce {
		twitchClient, err = c.createChannelBotClient(ctx, opts.SenderID, opts.BroadcasterID)
	} else {
		twitchClient, err = twitch.NewAppClientWithContext(ctx, c.config, c.twirBus)
	}
	if err != nil {
		return fmt.Errorf("create Twitch client: %w", err)
	}

	text := strings.ReplaceAll(opts.Message, "\n", " ")
	textParts := splitTextByLength(text)

	toxicity := make([]bool, len(textParts))
	// if !opts.SkipToxicityCheck {
	// 	t, err := c.toxicityCheck.CheckTextsToxicity(ctx, textParts)
	// 	if err != nil {
	// 		c.logger.Error("cannot check toxicity", logger.Error(err))
	// 		// return fmt.Errorf("cannot send message: %w", err)
	// 	} else {
	// 		toxicity = t
	// 	}
	// }

	for i, part := range textParts {
		// Do not send message if it was splitted more than 3 parts
		if i > 2 {
			return nil
		}

		message := part
		isToxic := !opts.SkipToxicityCheck && toxicity[i]
		if isToxic {
			if err := c.toxicMessagesRepository.Create(
				ctx,
				toxic_messages.CreateInput{
					ChannelID:     opts.BroadcasterID,
					ReplyToUserID: nil,
					Text:          part,
				},
			); err != nil {
				c.logger.Warn("Cannot save toxic message to db", logger.Error(err))
			}

			message = "[TwirApp] Redacted due toxicity validation. Contact support if you sure there is no toxicity."
		}

		if !opts.IsAnnounce {
			resp, err := twitchClient.Chat.SendChatMessage(
				ctx,
				helix.SendChatMessageRequest{
					BroadcasterID:        opts.BroadcasterID,
					SenderID:             opts.SenderID,
					Message:              validateResponseSlashes(message),
					ReplyParentMessageID: opts.ReplyParentMessageID,
				},
			)
			if err != nil {
				return fmt.Errorf("send chat message: %w", err)
			}

			var rateLimitGroup slog.Attr
			rateLimit := resp.Meta.RateLimit()
			if rateLimit.Valid() {
				rateLimitGroup = slog.Group(
					"rate_limit",
					slog.Int("limit", rateLimit.Limit()),
					slog.Int("remaining", rateLimit.Remaining()),
					slog.Int64("reset", rateLimit.Reset().Unix()),
				)
			}

			c.logger.Info(
				"✅ Message sent",
				slog.String("channel_id", opts.BroadcasterID),
				slog.String("text", message),
				slog.String("sender_id", opts.SenderID),
				slog.Bool("is_announce", opts.IsAnnounce),
				rateLimitGroup,
			)

			for _, m := range resp.Data {
				if m.DropReason != nil && m.DropReason.Message != "" {
					c.logger.Warn(
						"Message drop",
						slog.String("drop_reason", m.DropReason.Message),
						slog.String("code", m.DropReason.Code),
					)
					continue
				}

				go func() {
					createContext, cancel := context.WithTimeout(context.Background(), 10*time.Second)
					defer cancel()

					repoCreateError := c.sentMessagesRepository.Create(
						createContext,
						sentmessages.CreateInput{
							MessageTwitchID: m.MessageID,
							Content:         message,
							ChannelID:       opts.BroadcasterID,
							SenderID:        opts.SenderID,
						},
					)
					if repoCreateError != nil {
						c.logger.Warn("Cannot save message to db", slog.Any("err", repoCreateError))
					}
				}()
			}
		} else {
			var color bots.AnnounceColor
			if opts.AnnounceColor == bots.AnnounceColorRandom {
				color = bots.RandomAnnounceColor()
			} else {
				color = opts.AnnounceColor
			}

			_, err := twitchClient.Chat.SendChatAnnouncement(
				ctx,
				helix.SendChatAnnouncementRequest{
					BroadcasterID: opts.BroadcasterID,
					ModeratorID:   opts.SenderID,
					Message:       validateResponseSlashes(message),
					Color:         color.String(),
				},
			)
			if err != nil {
				return fmt.Errorf("send chat announcement: %w", err)
			}
		}
	}

	return nil
}

func activeTwitchBinding(
	binding channelplatformentity.ChannelPlatform,
) (channelplatformentity.TwitchBotConfig, bool, error) {
	twitchBinding, botConfig, found, err := (channelentity.Channel{
		Bindings: []channelplatformentity.ChannelPlatform{binding},
	}).TwitchBinding()
	if err != nil {
		return channelplatformentity.TwitchBotConfig{}, false, err
	}
	if !found || !twitchBinding.Enabled || twitchBinding.PlatformChannelID == "" ||
		!botConfig.IsBotMod || botConfig.IsTwitchBanned || botConfig.BotID == "" {
		return botConfig, false, nil
	}

	return botConfig, true, nil
}

const MAX_TWITCH_MESSAGE_LENGTH = 465

func splitTextByLength(text string) []string {
	var parts []string

	i := MAX_TWITCH_MESSAGE_LENGTH
	for utf8.RuneCountInString(text) > 0 {
		// the max twitch length of nickname is 32, and when we are replying to a message in chat,
		// and here is a reply, we need to reserve some space for that
		if utf8.RuneCountInString(text) < MAX_TWITCH_MESSAGE_LENGTH {
			parts = append(parts, text)
			break
		}
		runned := []rune(text)
		parts = append(parts, string(runned[:i]))
		text = string(runned[i:])
	}

	return parts
}
