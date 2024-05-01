package public

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
	model "github.com/satont/twir/libs/gomodels"
)

func (p *Public) HandleChannelCommandsGet(c *gin.Context) {
	channel := model.Channels{}
	if err := p.gorm.
		WithContext(c.Request.Context()).
		Where(`"id" = ?`, c.Param("channelId")).
		First(&channel).
		Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "channel not found"})
		return
	}

	if !channel.IsEnabled {
		c.JSON(http.StatusNotFound, gin.H{"error": "channel is disabled"})
		return
	}

	commands, err := p.cachedCommands.GetCommands(c.Request.Context(), c.Param("channelId"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	commands = lo.Filter(
		commands, func(item model.ChannelsCommands, index int) bool {
			return item.Enabled && item.Visible
		},
	)

	mappedCommands := make([]map[string]any, 0, len(commands))

	for _, command := range commands {
		cmd := map[string]any{
			"name":        command.Name,
			"description": command.Description,
			"module":      command.Module,
		}

		if command.Group != nil {
			cmd["group"] = command.Group.Name
		}

		responses := make([]map[string]any, 0, len(command.Responses))
		if len(command.Responses) > 0 {
			for _, response := range command.Responses {
				responses = append(
					responses, map[string]any{
						"text": response.Text.String,
					},
				)
			}

			cmd["responses"] = responses
		}

		mappedCommands = append(mappedCommands, cmd)
	}

	c.JSON(http.StatusOK, mappedCommands)
}
