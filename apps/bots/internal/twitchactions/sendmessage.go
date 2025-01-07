package twitchactions

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"unicode/utf8"

	"github.com/aidenwallis/go-ratelimiting/redis"
	"github.com/google/uuid"
	"github.com/nicklaw5/helix/v2"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/twitch"
	"github.com/twirapp/twir/libs/repositories/sentmessages"
)

type SendMessageOpts struct {
	BroadcasterID        string
	SenderID             string
	Message              string
	ReplyParentMessageID string
	IsAnnounce           bool
}

func validateResponseSlashes(response string) string {
	if strings.HasPrefix(response, "/me") || strings.HasPrefix(
		response,
		"/announce",
	) || strings.HasPrefix(response, "/shoutout") {
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
	resp, err := c.rateLimiter.Use(
		ctx,
		&redis.SlidingWindowOptions{
			Key:             fmt.Sprintf("bots:rate_limit:send_message:%s", opts.BroadcasterID),
			MaximumCapacity: 20,
			Window:          30,
		},
	)
	if err != nil {
		return err
	}
	if !resp.Success {
		return nil
	}

	channel := &model.Channels{}
	if err = c.gorm.
		WithContext(ctx).
		Where("channels.id = ?", opts.BroadcasterID).
		Joins("User").
		First(channel).Error; err != nil {
		return err
	}
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
	if twitchClientErr != nil {
		return twitchClientErr
	}

	text := strings.ReplaceAll(opts.Message, "\n", " ")
	textParts := splitTextByLength(text)

	for i, part := range textParts {
		if i > 2 {
			return nil
		}

		var msgErr error
		var errorMessage string

		if !opts.IsAnnounce {
			resp, err := twitchClient.SendChatMessage(
				&helix.SendChatMessageParams{
					BroadcasterID:        opts.BroadcasterID,
					SenderID:             opts.SenderID,
					Message:              validateResponseSlashes(part),
					ReplyParentMessageID: opts.ReplyParentMessageID,
				},
			)
			msgErr = err

			for _, m := range resp.Data.Messages {
				err := c.sentMessagesRepository.Create(
					ctx, sentmessages.CreateInput{
						MessageTwitchID: m.MessageID,
						Content:         part,
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
					Message:       validateResponseSlashes(part),
				},
			)
			msgErr = err

			if resp != nil {
				errorMessage = resp.ErrorMessage
			}

			if resp.ErrorMessage != "" && err == nil {
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
			return fmt.Errorf("cannot send message: %w", errorMessage)
		}
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
