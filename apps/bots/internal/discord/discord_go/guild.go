package discord_go

import (
	"context"
	"fmt"

	"github.com/diamondburned/arikawa/v3/discord"
)

type Guild struct {
	ID   string
	Name string
	Icon string
}

func (c *Discord) GetGuild(_ context.Context, guildID string) (Guild, error) {
	guildIdSnowflake, err := discord.ParseSnowflake(guildID)
	if err != nil {
		return Guild{}, err
	}

	guild, err := c.api.Guild(discord.GuildID(guildIdSnowflake))
	if err != nil {
		return Guild{}, err
	}
	if guild == nil {
		return Guild{}, fmt.Errorf("guild not found")
	}

	return Guild{
		ID:   guild.ID.String(),
		Name: guild.Name,
		Icon: guild.Icon,
	}, nil
}

func (c *Discord) LeaveGuild(_ context.Context, guildID string) error {
	guildIdSnowflake, err := discord.ParseSnowflake(guildID)
	if err != nil {
		return err
	}

	return c.api.LeaveGuild(discord.GuildID(guildIdSnowflake))
}

type GuildRole struct {
	ID    string
	Name  string
	Color string
}

func (c *Discord) GetGuildRoles(_ context.Context, guildID string) ([]GuildRole, error) {
	guildIdSnowflake, err := discord.ParseSnowflake(guildID)
	if err != nil {
		return nil, err
	}

	roles, err := c.api.Roles(discord.GuildID(guildIdSnowflake))
	if err != nil {
		return nil, err
	}

	result := make([]GuildRole, 0, len(roles))
	for _, role := range roles {
		result = append(
			result, GuildRole{
				ID:    role.ID.String(),
				Name:  role.Name,
				Color: fmt.Sprint(role.Color),
			},
		)
	}

	return result, nil
}
