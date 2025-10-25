package handler

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/goccy/go-json"
	"github.com/google/uuid"
	"github.com/kvizyx/twitchy/eventsub"
	"github.com/redis/go-redis/v9"
	emotes_cacher "github.com/twirapp/twir/libs/bus-core/emotes-cacher"
	"github.com/twirapp/twir/libs/bus-core/events"
	"github.com/twirapp/twir/libs/bus-core/twitch"
	"github.com/twirapp/twir/libs/redis_keys"
	channelscommandsprefixrepository "github.com/twirapp/twir/libs/repositories/channels_commands_prefix"
	channelscommandsprefixmodel "github.com/twirapp/twir/libs/repositories/channels_commands_prefix/model"
	streamsmodel "github.com/twirapp/twir/libs/repositories/streams/model"
	"go.opentelemetry.io/otel/attribute"
	"golang.org/x/sync/errgroup"
)

func convertFragmentTypeToEnumValue(t string) twitch.FragmentType {
	switch t {
	case "text":
		return twitch.FragmentType_TEXT
	case "cheermote":
		return twitch.FragmentType_CHEERMOTE
	case "emote":
		return twitch.FragmentType_EMOTE
	case "mention":
		return twitch.FragmentType_MENTION
	default:
		return twitch.FragmentType_TEXT
	}
}

func (c *Handler) HandleChannelChatMessage(
	ctx context.Context,
	event eventsub.ChannelChatMessageEvent,
	meta eventsub.WebsocketNotificationMetadata,
) {
	_, span := c.tracer.Start(ctx, "HandleChannelChatMessage")
	span.SetAttributes(
		attribute.String("message_id", event.MessageId),
		attribute.String("channel_id", event.BroadcasterUserId),
	)
	defer span.End()

	fragments := make([]twitch.ChatMessageMessageFragment, 0, len(event.Message.Fragments))

	startFragmentPosition := 0
	for _, fragment := range event.Message.Fragments {
		var cheerMote *twitch.ChatMessageMessageFragmentCheermote
		var emote *twitch.ChatMessageMessageFragmentEmote
		var mention *twitch.ChatMessageMessageFragmentMention

		if fragment.Cheermote != nil {
			cheerMote = &twitch.ChatMessageMessageFragmentCheermote{
				Prefix: fragment.Cheermote.Prefix,
				Bits:   int64(fragment.Cheermote.Bits),
				Tier:   int64(fragment.Cheermote.Tier),
			}
		}

		if fragment.Emote != nil {
			formats := make([]string, 0, len(fragment.Emote.Format))
			for _, f := range fragment.Emote.Format {
				formats = append(formats, string(f))
			}

			emote = &twitch.ChatMessageMessageFragmentEmote{
				Id:         fragment.Emote.Id,
				EmoteSetId: fragment.Emote.EmoteSetId,
				OwnerId:    fragment.Emote.OwnerId,
				Format:     formats,
			}
		}

		if fragment.Mention != nil {
			mention = &twitch.ChatMessageMessageFragmentMention{
				UserId:    fragment.Mention.UserId,
				UserName:  fragment.Mention.UserName,
				UserLogin: fragment.Mention.UserLogin,
			}
		}

		position := twitch.ChatMessageMessageFragmentPosition{
			Start: startFragmentPosition,
			End:   startFragmentPosition + utf8.RuneCountInString(fragment.Text),
		}

		fragments = append(
			fragments,
			twitch.ChatMessageMessageFragment{
				Type:      convertFragmentTypeToEnumValue(string(fragment.Type)),
				Text:      fragment.Text,
				Cheermote: cheerMote,
				Emote:     emote,
				Mention:   mention,
				Position:  position,
			},
		)

		startFragmentPosition += utf8.RuneCountInString(fragment.Text)
	}

	badges := make([]twitch.ChatMessageBadge, 0, len(event.Badges))
	for _, badge := range event.Badges {
		badges = append(
			badges,
			twitch.ChatMessageBadge{
				Id:    badge.Id,
				SetId: badge.SetId,
				Info:  badge.Info,
			},
		)
	}

	var cheer *twitch.ChatMessageCheer
	if event.Cheer != nil {
		cheer = &twitch.ChatMessageCheer{Bits: int64(event.Cheer.Bits)}
	}

	var reply *twitch.ChatMessageReply
	if event.Reply != nil {
		reply = &twitch.ChatMessageReply{
			ParentMessageId:   event.Reply.ParentMessageId,
			ParentMessageBody: event.Reply.ParentMessageBody,
			ParentUserId:      event.Reply.ParentUserId,
			ParentUserName:    event.Reply.ParentUserName,
			ParentUserLogin:   event.Reply.ParentUserLogin,
			ThreadMessageId:   event.Reply.ThreadMessageId,
			ThreadUserId:      event.Reply.ThreadUserId,
			ThreadUserName:    event.Reply.ThreadUserName,
			ThreadUserLogin:   event.Reply.ThreadUserLogin,
		}
	}

	data := twitch.TwitchChatMessage{
		ID:                   event.MessageId,
		BroadcasterUserId:    event.BroadcasterUserId,
		BroadcasterUserName:  event.BroadcasterUserName,
		BroadcasterUserLogin: event.BroadcasterUserLogin,
		ChatterUserId:        event.ChatterUserId,
		ChatterUserName:      event.ChatterUserName,
		ChatterUserLogin:     event.ChatterUserLogin,
		MessageId:            event.MessageId,
		Message: &twitch.ChatMessageMessage{
			Text:      event.Message.Text,
			Fragments: fragments,
		},
		Color:                       event.Color,
		Badges:                      badges,
		MessageType:                 string(event.Type),
		Cheer:                       cheer,
		Reply:                       reply,
		ChannelPointsCustomRewardId: event.ChannelPointsCustomRewardId,
	}

	var errwg errgroup.Group

	errwg.Go(
		func() error {
			emotes, emotesErr := c.chatMessageCountEmotes(ctx, data)
			if emotesErr != nil {
				c.logger.Error("cannot count emotes", slog.Any("err", emotesErr))
			}
			data.EnrichedData.UsedEmotesWithThirdParty = emotes

			return nil
		},
	)

	errwg.Go(
		func() error {
			commandsPrefix, err := c.chatMessageGetChannelCommandPrefix(ctx, data.BroadcasterUserId)
			if err != nil {
				return err
			}
			data.EnrichedData.ChannelCommandPrefix = commandsPrefix

			return nil
		},
	)

	errwg.Go(
		func() error {
			channel, err := c.channelsCache.Get(ctx, data.BroadcasterUserId)
			if err != nil {
				return err
			}
			data.EnrichedData.DbChannel = channel

			return nil
		},
	)

	errwg.Go(
		func() error {
			stream, err := c.chatMessageGetChannelStream(ctx, data.BroadcasterUserId)
			if err != nil {
				return err
			}

			data.EnrichedData.ChannelStream = stream

			return nil
		},
	)

	if err := errwg.Wait(); err != nil {
		c.logger.Error("cannot handle message", slog.Any("err", err))
		return
	}

	if err := c.twirBus.ChatMessages.Publish(ctx, data); err != nil {
		c.logger.Error("cannot publish message for handle", slog.Any("err", err))
	}

	isCommand := strings.HasPrefix(data.Message.Text, data.EnrichedData.ChannelCommandPrefix)
	// ignore bot himself from chat commands
	if isCommand && data.ChatterUserId == data.EnrichedData.DbChannel.BotID && c.config.AppEnv == "production" {
		return
	} else if isCommand && data.EnrichedData.DbChannel.IsEnabled {
		if err := c.twirBus.Parser.ProcessMessageAsCommand.Publish(ctx, data); err != nil {
			c.logger.Error("cannot publish process command", slog.Any("err", err))
		}
	}
}

func (c *Handler) HandleChannelChatMessageDelete(
	ctx context.Context,
	event eventsub.ChannelChatMessageDeleteEvent,
	meta eventsub.WebsocketNotificationMetadata,
) {
	c.logger.Info(
		"message delete",
		slog.String("channelId", event.BroadcasterUserId),
		slog.String("channelName", event.BroadcasterUserName),
		slog.String("userId", event.TargetUserId),
		slog.String("userName", event.TargetUserLogin),
	)

	if err := c.twirBus.Events.ChannelMessageDelete.Publish(
		ctx,
		events.ChannelMessageDeleteMessage{
			BaseInfo: events.BaseInfo{
				ChannelID:   event.BroadcasterUserId,
				ChannelName: event.BroadcasterUserLogin,
			},
			UserId:               event.TargetUserId,
			UserName:             event.TargetUserLogin,
			UserLogin:            event.TargetUserName,
			BroadcasterUserName:  event.BroadcasterUserName,
			BroadcasterUserLogin: event.BroadcasterUserLogin,
			MessageId:            event.MessageId,
		},
	); err != nil {
		c.logger.Error(err.Error(), slog.Any("err", err))
	}
}

func (c *Handler) chatMessageCountEmotes(
	ctx context.Context,
	msg twitch.TwitchChatMessage,
) (map[string]int, error) {
	if msg.Message == nil {
		return nil, nil
	}

	splittedMsg := strings.Fields(msg.Message.Text)
	emotes := make(map[string]int, len(splittedMsg))

	var (
		wg            errgroup.Group
		globalEmotes  []emotes_cacher.Emote
		channelEmotes []emotes_cacher.Emote
	)

	wg.Go(
		func() error {
			e, err := c.twirBus.EmotesCacher.GetGlobalEmotes.Request(
				ctx,
				emotes_cacher.GetGlobalEmotesRequest{},
			)
			if err != nil {
				return err
			}

			globalEmotes = e.Data.Emotes

			return nil
		},
	)

	wg.Go(
		func() error {
			e, err := c.twirBus.EmotesCacher.GetChannelEmotes.Request(
				ctx,
				emotes_cacher.GetChannelEmotesRequest{
					ChannelID: msg.BroadcasterUserId,
				},
			)
			if err != nil {
				return err
			}

			channelEmotes = e.Data.Emotes

			return nil
		},
	)

	if err := wg.Wait(); err != nil {
		c.logger.Error("failed to fetch global emotes", slog.Any("err", err))
		return nil, err
	}

	for _, f := range msg.Message.Fragments {
		if f.Type != twitch.FragmentType_EMOTE || f.Text == "" || f.Text == " " {
			continue
		}
		emotes[f.Text] += 1
	}

	for _, part := range splittedMsg {
		if part == "" || part == " " {
			continue
		}

		// do not make redis requests if part already present in map
		var isTwitchEmote bool
		for _, fragment := range msg.Message.Fragments {
			if fragment.Emote != nil && strings.TrimSpace(fragment.Text) == part {
				isTwitchEmote = true
				break
			}

			if fragment.Cheermote != nil && strings.TrimSpace(fragment.Text) == part {
				isTwitchEmote = true
				break
			}
		}

		if emote, ok := emotes[part]; !isTwitchEmote && ok {
			emotes[part] = emote + 1
			continue
		}

		for _, globalEmote := range globalEmotes {
			if globalEmote.Name == part {
				emotes[part] += 1
				continue
			}
		}

		for _, channelEmote := range channelEmotes {
			if channelEmote.Name == part {
				emotes[part] += 1
				continue
			}
		}
	}

	return emotes, nil
}

func (c *Handler) chatMessageGetChannelCommandPrefix(ctx context.Context, channelId string) (
	string,
	error,
) {
	commandsPrefix := "!"
	fetchedCommandsPrefix, err := c.prefixCache.Get(ctx, channelId)
	if err != nil && !errors.Is(err, channelscommandsprefixrepository.ErrNotFound) {
		return "", err
	}

	if fetchedCommandsPrefix != channelscommandsprefixmodel.Nil {
		commandsPrefix = fetchedCommandsPrefix.Prefix
	} else {
		prefixCtx := context.WithoutCancel(ctx)

		go func() {
			if err := c.prefixCache.SetValue(
				prefixCtx,
				channelId,
				channelscommandsprefixmodel.ChannelsCommandsPrefix{
					ID:        uuid.New(),
					ChannelID: channelId,
					Prefix:    commandsPrefix,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
			); err != nil {
				c.logger.Error("cannot set default command prefix", slog.Any("err", err))
			}
		}()
	}

	return commandsPrefix, nil
}

func (c *Handler) chatMessageGetChannelStream(
	ctx context.Context,
	channelId string,
) (*streamsmodel.Stream, error) {
	cacheKey := redis_keys.StreamByChannelID(channelId)
	cachedBytes, err := c.redisClient.Get(ctx, cacheKey).Bytes()
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, fmt.Errorf("failed to get stream cache: %w", err)
	}

	if len(cachedBytes) > 0 {
		var stream streamsmodel.Stream
		if err := json.Unmarshal(cachedBytes, &stream); err != nil {
			return nil, err
		}

		return &stream, nil
	}

	stream, err := c.streamsrepository.GetByChannelID(ctx, channelId)
	if err != nil {
		return nil, fmt.Errorf("failed to get stream by channel id: %w", err)
	}

	if stream.ID == "" {
		return nil, nil
	}

	streamBytes, err := json.Marshal(stream)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal stream: %w", err)
	}

	if err := c.redisClient.Set(ctx, cacheKey, streamBytes, 30*time.Second).Err(); err != nil {
		return nil, fmt.Errorf("failed to set stream cache: %w", err)
	}

	return &stream, nil
}
