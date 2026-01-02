package discord_go

import (
	"context"

	"github.com/diamondburned/arikawa/v3/discord"
)

type DiscordChannelType string

const (
	DiscordChannelTypeText         DiscordChannelType = "text"
	DiscordChannelTypeVoice        DiscordChannelType = "voice"
	DiscordChannelTypeCategory     DiscordChannelType = "category"
	DiscordChannelTypeAnnouncement DiscordChannelType = "announcement"
	DiscordChannelTypeUnknown      DiscordChannelType = "unknown"
)

type DiscordChannel struct {
	ID              string
	Name            string
	Type            DiscordChannelType
	CanSendMessages bool
}

func (c *Discord) GetGuildChannels(_ context.Context, guildID string) ([]DiscordChannel, error) {
	guildIdSnowflake, err := discord.ParseSnowflake(guildID)
	if err != nil {
		return nil, err
	}

	ch, err := c.api.Channels(discord.GuildID(guildIdSnowflake))
	if err != nil {
		return nil, err
	}

	result := make([]DiscordChannel, 0, len(ch))
	for _, channel := range ch {
		var channelType DiscordChannelType
		switch channel.Type {
		case discord.GuildText:
			channelType = DiscordChannelTypeText
		case discord.GuildAnnouncement:
			channelType = DiscordChannelTypeAnnouncement
		case discord.GuildVoice, discord.GuildStageVoice:
			channelType = DiscordChannelTypeVoice
		case discord.GuildCategory:
			channelType = DiscordChannelTypeCategory
		default:
			channelType = DiscordChannelTypeUnknown
		}

		result = append(
			result, DiscordChannel{
				ID:              channel.ID.String(),
				Name:            channel.Name,
				Type:            channelType,
				CanSendMessages: channel.SelfPermissions.Has(discord.PermissionSendMessages),
			},
		)
	}

	return result, nil
}
