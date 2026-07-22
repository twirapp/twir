package messagehandler

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/twirapp/twir/apps/bots/internal/channelbinding"
	emotes_cacher "github.com/twirapp/twir/libs/bus-core/emotes-cacher"
	"github.com/twirapp/twir/libs/bus-core/generic"
	"github.com/twirapp/twir/libs/entities/platform"
	"github.com/twirapp/twir/libs/logger"
	channelplatformsmodel "github.com/twirapp/twir/libs/repositories/channel_platforms/model"
	channelsmodel "github.com/twirapp/twir/libs/repositories/channels/model"
	channelscommandsprefixrepository "github.com/twirapp/twir/libs/repositories/channels_commands_prefix"
	streamsmodel "github.com/twirapp/twir/libs/repositories/streams/model"
	usersmodel "github.com/twirapp/twir/libs/repositories/users/model"
	usersstatsmodel "github.com/twirapp/twir/libs/repositories/users_stats/model"
	"golang.org/x/sync/errgroup"
)

type enrichedChatMessage struct {
	generic.ChatMessage
	EnrichedData chatMessageEnrichedData
}

type chatMessageEnrichedData struct {
	Binding                  channelplatformsmodel.ChannelPlatform
	DbChannel                channelsmodel.Channel
	BotPlatformID            string
	ChannelStream            *streamsmodel.Stream
	DbUser                   *usersmodel.User
	DbUserChannelStat        *usersstatsmodel.UserStat
	ChannelCommandPrefix     string
	UsedEmotesWithThirdParty map[string]int
}

func findChatMessageBinding(
	channel channelsmodel.Channel,
	message generic.ChatMessage,
	messagePlatform platform.Platform,
) (channelplatformsmodel.ChannelPlatform, bool, error) {
	if message.ChannelBindingID == "" {
		binding, found := channelbinding.Find(channel, messagePlatform)
		return binding, found, nil
	}

	bindingID, err := uuid.Parse(message.ChannelBindingID)
	if err != nil {
		return channelplatformsmodel.ChannelPlatform{}, false, fmt.Errorf("parse message channel binding id: %w", err)
	}

	binding, found := channelbinding.FindByID(channel, bindingID)
	return binding, found, nil
}

func (c *MessageHandler) enrichChatMessage(
	ctx context.Context,
	message generic.ChatMessage,
) (enrichedChatMessage, error) {
	result := enrichedChatMessage{
		ChatMessage: message,
		EnrichedData: chatMessageEnrichedData{
			ChannelCommandPrefix: "!",
		},
	}

	channelID, err := uuid.Parse(message.ChannelID)
	if err != nil {
		return result, fmt.Errorf("parse message channel id: %w", err)
	}

	channel, err := c.channelsRepository.GetByID(ctx, channelID)
	if err != nil {
		return result, fmt.Errorf("get message channel: %w", err)
	}
	result.EnrichedData.DbChannel = channel

	messagePlatform := chatMessagePlatform(message)
	binding, found, err := findChatMessageBinding(channel, message, messagePlatform)
	if err != nil {
		return result, err
	}
	if !found || !binding.Enabled {
		return result, nil
	}
	result.EnrichedData.Binding = binding
	if messagePlatform == platform.PlatformTwitch {
		botConfig, err := channelbinding.ParseTwitchBotConfig(binding)
		if err != nil {
			return result, fmt.Errorf("parse Twitch bot config: %w", err)
		}
		result.EnrichedData.BotPlatformID = botConfig.BotID
	}

	userID, err := uuid.Parse(message.UserID)
	if err != nil {
		return result, fmt.Errorf("parse message user id: %w", err)
	}

	user, err := c.usersRepository.GetByID(ctx, userID)
	if err != nil {
		return result, fmt.Errorf("get message user: %w", err)
	}
	result.EnrichedData.DbUser = &user

	userStats, err := c.usersstatsRepository.GetByUserAndChannelID(ctx, userID, channelID)
	if err != nil {
		return result, fmt.Errorf("get message user stats: %w", err)
	}
	result.EnrichedData.DbUserChannelStat = userStats

	if c.prefixCache != nil {
		prefix, err := c.prefixCache.Get(ctx, channelID.String())
		if err != nil && !errors.Is(err, channelscommandsprefixrepository.ErrNotFound) {
			return result, fmt.Errorf("get message command prefix: %w", err)
		}
		if err == nil && prefix.Prefix != "" {
			result.EnrichedData.ChannelCommandPrefix = prefix.Prefix
		}
	}

	stream, err := c.streamsRepository.GetByChannelID(ctx, channelID, messagePlatform)
	if err != nil {
		return result, fmt.Errorf("get message channel stream: %w", err)
	}
	if stream.ID != "" {
		result.EnrichedData.ChannelStream = &stream
	}

	emotes, err := c.chatMessageCountEmotes(ctx, message)
	if err != nil {
		c.logger.Error("cannot count emotes", logger.Error(err))
	} else {
		result.EnrichedData.UsedEmotesWithThirdParty = emotes
	}

	return result, nil
}

func chatMessagePlatform(message generic.ChatMessage) platform.Platform {
	if message.Platform == "" {
		return platform.PlatformTwitch
	}

	return platform.Platform(message.Platform)
}

func (c *MessageHandler) chatMessageCountEmotes(
	ctx context.Context,
	message generic.ChatMessage,
) (map[string]int, error) {
	if message.Message == nil || chatMessagePlatform(message) != platform.PlatformTwitch {
		return nil, nil
	}

	emotes := make(map[string]int)
	for _, fragment := range message.Message.Fragments {
		if fragment.Type == generic.FragmentType_EMOTE && fragment.Text != "" && fragment.Text != " " {
			emotes[fragment.Text]++
		}
	}

	var (
		globalEmotes  []emotes_cacher.Emote
		channelEmotes []emotes_cacher.Emote
		group         errgroup.Group
	)

	group.Go(
		func() error {
			response, err := c.twirBus.EmotesCacher.GetGlobalEmotes.Request(
				ctx,
				emotes_cacher.GetGlobalEmotesRequest{},
			)
			if err != nil {
				return err
			}

			globalEmotes = response.Data.Emotes
			return nil
		},
	)
	group.Go(
		func() error {
			response, err := c.twirBus.EmotesCacher.GetChannelEmotes.Request(
				ctx,
				emotes_cacher.GetChannelEmotesRequest{
					Platform:  platform.PlatformTwitch,
					ChannelID: message.PlatformChannelID,
				},
			)
			if err != nil {
				return err
			}

			channelEmotes = response.Data.Emotes
			return nil
		},
	)

	if err := group.Wait(); err != nil {
		return nil, err
	}

	for _, part := range strings.Fields(message.Message.Text) {
		if part == "" || part == " " {
			continue
		}

		var isNativeEmote bool
		for _, fragment := range message.Message.Fragments {
			if (fragment.Emote != nil || fragment.Cheermote != nil) && strings.TrimSpace(fragment.Text) == part {
				isNativeEmote = true
				break
			}
		}

		if emote, ok := emotes[part]; !isNativeEmote && ok {
			emotes[part] = emote + 1
			continue
		}

		for _, emote := range globalEmotes {
			if emote.Name == part {
				emotes[part]++
			}
		}
		for _, emote := range channelEmotes {
			if emote.Name == part {
				emotes[part]++
			}
		}
	}

	return emotes, nil
}
