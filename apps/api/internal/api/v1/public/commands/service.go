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
	Cooldown     int64    `json:"cooldown"`
	CooldownType string   `json:"cooldownType"`
	Aliases      []string `json:"aliases"`
	Description  *string  `json:"description"`
	Permissions  []string `json:"permissions"`
}

func handleGet(channelId string, services types.Services) ([]Command, error) {
	commands := []model.ChannelsCommands{}
	channelRoles := []model.ChannelRole{}

	err := services.DB.
		Where(`"channelId" = ? AND "enabled" = ? AND "visible" = ?`, channelId, true, true).
		Preload("Responses").
		Find(&commands).Error

	err = services.DB.Where(`"channelId" = ?`, channelId).Find(&channelRoles).Error

	if err != nil {
		return nil, fiber.NewError(http.StatusNotFound, "cannot find commands")
	}

	commandsResponse := []Command{}

	for _, cmd := range commands {
		responses := lo.Map(cmd.Responses, func(item model.ChannelsCommandsResponses, _ int) string {
			return item.Text.String
		})

		roles := []string{}

		for _, role := range cmd.RolesIDS {
			r, ok := lo.Find(channelRoles, func(item model.ChannelRole) bool {
				return item.ID == role
			})
			if !ok {
				continue
			}
			roles = append(roles, r.Name)
		}

		commandsResponse = append(commandsResponse, Command{
			Name:         cmd.Name,
			Responses:    responses,
			Cooldown:     cmd.Cooldown.Int64,
			CooldownType: cmd.CooldownType,
			Aliases:      cmd.Aliases,
			Description:  cmd.Description.Ptr(),
			Permissions:  roles,
		})
	}

	sort.Slice(commandsResponse, func(i, j int) bool {
		return commandsResponse[i].Name < commandsResponse[j].Name
	})

	return commandsResponse, nil
}
