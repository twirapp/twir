package twitchactions

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/aidenwallis/go-ratelimiting/redis"
	"github.com/google/uuid"
	"github.com/nicklaw5/helix/v2"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/twitch"
	"github.com/twirapp/twir/libs/repositories/sentmessages"
	"github.com/twirapp/twir/libs/repositories/toxic_messages"
)

type SendMessageOpts struct {
	BroadcasterID        string
	SenderID             string
	Message              string
	ReplyParentMessageID string
	IsAnnounce           bool
	SkipToxicityCheck    bool
}

const shoutOutPrefix = "/shoutout"

func validateResponseSlashes(response string) string {
	if strings.HasPrefix(response, "/me") || strings.HasPrefix(
		response,
		"/announce",
	) || strings.HasPrefix(response, shoutOutPrefix) {
		return response
	} else if strings.HasPrefix(response, "/") {
		return "Slash commands except /me, /announce and /shoutout is disallowed. This response wont be ever sended."
	} else if strings.HasPrefix(response, ".") {
		return `Message cannot start with "." symbol.`
	} else {
		return response
	}
}

func (c *TwitchActions) SendMessage(ctx context.Context, opts SendMessageOpts) error {
	rateLimiterStart := time.Now()
	resp, err := c.rateLimiter.Use(
		ctx,
		&redis.SlidingWindowOptions{
			Key:             fmt.Sprintf("bots:rate_limit:send_message:%s", opts.BroadcasterID),
			MaximumCapacity: 20,
			Window:          30,
		},
	)
	c.logger.Info("Rate limiter", slog.Duration("duration", time.Since(rateLimiterStart)))
	if err != nil {
		return err
	}
	if !resp.Success {
		return nil
	}

	channel := &model.Channels{}
	channelsStart := time.Now()
	if err = c.gorm.
		WithContext(ctx).
		Where("channels.id = ?", opts.BroadcasterID).
		Joins("User").
		First(channel).Error; err != nil {
		return err
	}
	c.logger.Info("Get channel", slog.Duration("duration", time.Since(channelsStart)))
	if !channel.IsEnabled || !channel.IsBotMod || channel.IsTwitchBanned || channel.User.IsBanned {
		return nil
	}

	c.logger.Info(
		"Sending message",
		slog.String("channel_id", opts.BroadcasterID),
		slog.String("sender_id", opts.SenderID),
		slog.Bool("is_announce", opts.IsAnnounce),
	)

	var twitchClient *helix.Client
	var twitchClientErr error
	twitchClientStart := time.Now()
	if !opts.IsAnnounce {
		twitchClient, twitchClientErr = twitch.NewAppClientWithContext(ctx, c.config, c.tokensGrpc)
	} else {
		twitchClient, twitchClientErr = twitch.NewBotClientWithContext(
			ctx,
			opts.SenderID,
			c.config,
			c.tokensGrpc,
		)
	}
	c.logger.Info("Twitch client", slog.Duration("duration", time.Since(twitchClientStart)))
	if twitchClientErr != nil {
		return twitchClientErr
	}

	text := strings.ReplaceAll(opts.Message, "\n", " ")
	textParts := splitTextByLength(text)

	toxicity := make([]bool, len(textParts))
	if !opts.SkipToxicityCheck {
		start := time.Now()
		t, err := c.toxicityCheck.CheckTextsToxicity(ctx, textParts)
		c.logger.Info("Toxicity check", slog.Duration("duration", time.Since(start)))
		if err != nil {
			return fmt.Errorf("cannot send message: %w", err)
		}
		toxicity = t
	}

	for i, part := range textParts {
		// Do not send message if it was splitted more than 3 parts
		if i > 2 {
			return nil
		}

		var msgErr error
		var errorMessage string

		message := part
		isToxic := !opts.SkipToxicityCheck && toxicity[i]
		if isToxic {
			toxicMessageStart := time.Now()
			if err := c.toxicMessagesRepository.Create(
				ctx,
				toxic_messages.CreateInput{
					ChannelID:     opts.BroadcasterID,
					ReplyToUserID: nil,
					Text:          part,
				},
			); err != nil {
				c.logger.Warn("Cannot save toxic message to db", slog.Any("err", err))
			} else {
				c.logger.Info(
					"Save toxic message",
					slog.Duration("duration", time.Since(toxicMessageStart)),
				)
			}

			message = "[TwirApp] Redacted due toxicity validation. Contact support if you sure there is no toxicity."
		}

		sendStart := time.Now()
		if !opts.IsAnnounce {
			resp, err := twitchClient.SendChatMessage(
				&helix.SendChatMessageParams{
					BroadcasterID:        opts.BroadcasterID,
					SenderID:             opts.SenderID,
					Message:              validateResponseSlashes(message),
					ReplyParentMessageID: opts.ReplyParentMessageID,
				},
			)
			msgErr = err

			for _, m := range resp.Data.Messages {
				err := c.sentMessagesRepository.Create(
					ctx, sentmessages.CreateInput{
						MessageTwitchID: m.MessageID,
						Content:         message,
						ChannelID:       opts.BroadcasterID,
						SenderID:        opts.SenderID,
					},
				)
				if err != nil {
					c.logger.Warn("Cannot save message to db", slog.Any("err", err))
				}
			}

			if resp != nil {
				errorMessage = resp.ErrorMessage
			}
		} else {
			resp, err := twitchClient.SendChatAnnouncement(
				&helix.SendChatAnnouncementParams{
					BroadcasterID: opts.BroadcasterID,
					ModeratorID:   opts.SenderID,
					Message:       validateResponseSlashes(message),
				},
			)
			msgErr = err

			if resp != nil {
				errorMessage = resp.ErrorMessage
			}

			if resp != nil && resp.ErrorMessage != "" && err == nil {
				err := c.sentMessagesRepository.Create(
					ctx, sentmessages.CreateInput{
						MessageTwitchID: uuid.NewString(),
						Content:         part,
						ChannelID:       opts.BroadcasterID,
						SenderID:        opts.SenderID,
					},
				)
				if err != nil {
					c.logger.Warn("Cannot save message to db", slog.Any("err", err))
				}
			}
		}

		if msgErr != nil {
			return err
		}

		if errorMessage != "" {
			return fmt.Errorf("cannot send message: %s", errorMessage)
		}

		c.logger.Info("Message sent", slog.Duration("duration", time.Since(sendStart)))
	}

	return nil
}

func splitTextByLength(text string) []string {
	var parts []string

	i := 500
	for utf8.RuneCountInString(text) > 0 {
		if utf8.RuneCountInString(text) < 500 {
			parts = append(parts, text)
			break
		}
		runned := []rune(text)
		parts = append(parts, string(runned[:i]))
		text = string(runned[i:])
	}

	return parts
}
