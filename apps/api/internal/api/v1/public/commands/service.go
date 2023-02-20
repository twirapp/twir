package commands

import (
	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
	"github.com/satont/tsuwari/apps/api/internal/types"
	model "github.com/satont/tsuwari/libs/gomodels"
	"net/http"
	"sort"
)

type Command struct {
	Name         string   `json:"name"`
	Responses    []string `json:"responses"`
	Permission   string   `json:"permission"`
	Cooldown     int64    `json:"cooldown"`
	CooldownType string   `json:"cooldownType"`
	Aliases      []string `json:"aliases"`
	Description  *string  `json:"description"`
}

var CommandPerms = []string{"BROADCASTER", "MODERATOR", "VIP", "SUBSCRIBER", "FOLLOWER", "VIEWER"}

func handleGet(channelId string, services types.Services) ([]Command, error) {
	commands := []model.ChannelsCommands{}

	err := services.DB.
		Where(`"channelId" = ? AND "enabled" = ? AND "visible" = ?`, channelId, true, true).
		Preload("Responses").
		Find(&commands).Error

	if err != nil {
		return nil, fiber.NewError(http.StatusNotFound, "cannot find commands")
	}

	commandsResponse := []Command{}

	for _, cmd := range commands {
		responses := lo.Map(cmd.Responses, func(item model.ChannelsCommandsResponses, _ int) string {
			return item.Text.String
		})

		commandsResponse = append(commandsResponse, Command{
			Name:         cmd.Name,
			Responses:    responses,
			Permission:   cmd.Permission,
			Cooldown:     cmd.Cooldown.Int64,
			CooldownType: cmd.CooldownType,
			Aliases:      cmd.Aliases,
			Description:  cmd.Description.Ptr(),
		})
	}

	sort.Slice(commandsResponse, func(i, j int) bool {
		iPermIndex := lo.IndexOf(CommandPerms, commandsResponse[i].Permission)
		jPermIndex := lo.IndexOf(CommandPerms, commandsResponse[j].Permission)

		if iPermIndex == jPermIndex {
			return commandsResponse[i].Name < commandsResponse[j].Name
		}

		return iPermIndex < jPermIndex
	})

	return commandsResponse, nil
}
