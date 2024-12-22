package http_public

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/twirapp/twir/apps/api-gql/internal/services/channels"
)

type commandDto struct {
	Name        string               `json:"name"`
	Description *string              `json:"description"`
	Module      string               `json:"module"`
	Group       *string              `json:"group"`
	Responses   []commandDtoResponse `json:"responses"`
}

type commandDtoResponse struct {
	Text string `json:"text"`
}

func (p *Public) HandleChannelCommandsGet(c *gin.Context) {
	// TODO: refactor to service
	channel, err := p.channelsService.GetByID(c.Request.Context(), c.Param("channelId"))
	if err != nil {
		if errors.Is(err, channels.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "channel not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !channel.IsEnabled {
		c.JSON(http.StatusNotFound, gin.H{"error": "channel is disabled"})
		return
	}

	// TODO: it must use generic cacher with repository instead of gorm cacher
	commands, err := p.cachedCommands.Get(c.Request.Context(), c.Param("channelId"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	commands = lo.Filter(
		commands, func(item model.ChannelsCommands, index int) bool {
			return item.Enabled && item.Visible
		},
	)

	mappedCommands := make([]commandDto, 0, len(commands))

	for _, command := range commands {
		mappedCmd := commandDto{
			Name:        command.Name,
			Description: command.Description.Ptr(),
			Module:      command.Module,
			Group:       nil,
			Responses:   make([]commandDtoResponse, 0, len(command.Responses)),
		}

		if command.Group != nil {
			mappedCmd.Group = &command.Group.Name
		}

		for _, response := range command.Responses {
			mappedCmd.Responses = append(
				mappedCmd.Responses,
				commandDtoResponse{
					Text: response.Text.String,
				},
			)
		}

		mappedCommands = append(mappedCommands, mappedCmd)
	}

	c.JSON(http.StatusOK, mappedCommands)
}
