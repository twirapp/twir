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
	"github.com/twirapp/twir/libs/twitch"
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
	SkipRateLimits       bool
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

	channel, err := c.channelsCache.Get(ctx, opts.BroadcasterID)
	if err != nil {
		return err
	}
	if !channel.IsEnabled || !channel.IsBotMod || channel.IsTwitchBanned {
		return nil
	}

	if opts.SenderID == "" {
		opts.SenderID = channel.BotID
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
		twitchClient, twitchClientErr = twitch.NewAppClientWithContext(ctx, c.config, c.twirBus)
	} else {
		twitchClient, twitchClientErr = twitch.NewBotClientWithContext(
			ctx,
			opts.SenderID,
			c.config,
			c.twirBus,
		)
	}

	if twitchClientErr != nil {
		return twitchClientErr
	}

	text := strings.ReplaceAll(opts.Message, "\n", " ")
	textParts := splitTextByLength(text)

	toxicity := make([]bool, len(textParts))
	// if !opts.SkipToxicityCheck {
	// 	t, err := c.toxicityCheck.CheckTextsToxicity(ctx, textParts)
	// 	if err != nil {
	// 		c.logger.Error("cannot check toxicity", slog.Any("err", err))
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

		var msgErr error
		var errorMessage string

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
				c.logger.Warn("Cannot save toxic message to db", slog.Any("err", err))
			}

			message = "[TwirApp] Redacted due toxicity validation. Contact support if you sure there is no toxicity."
		}

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
			if resp == nil {
				return fmt.Errorf("cannot send message with unknown reason: %w", err)
			}

			var rateLimitGroup slog.Attr
			if resp != nil {
				rateLimitGroup = slog.Group(
					"rate_limit",
					slog.Int("limit", resp.GetRateLimit()),
					slog.Int("remaining", resp.GetRateLimitRemaining()),
					slog.Int("reset", resp.GetRateLimitReset()),
				)
			}

			c.logger.Info(
				"Message sent",
				slog.String("channel_id", opts.BroadcasterID),
				slog.String("sender_id", opts.SenderID),
				rateLimitGroup,
			)

			for _, m := range resp.Data.Messages {
				if m.DropReasons.Data.Message != "" {
					c.logger.Warn(
						"Message drop",
						slog.String("drop_reason", m.DropReasons.Data.Message),
						slog.String("code", m.DropReasons.Data.Code),
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
