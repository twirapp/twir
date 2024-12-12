package resolvers

import (
	"context"
	"fmt"
	"slices"

	model "github.com/satont/twir/libs/gomodels"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

var commandAlreadyExistsError = gqlerror.Error{
	Message:   "command with this name or aliase already exists",
	Path:      nil,
	Locations: nil,
	Extensions: map[string]interface{}{
		"code": "ALREADY_EXISTS",
	},
}

func (r *Resolver) checkIsCommandWithNameOrAliaseExists(
	ctx context.Context,
	commandId *string,
	name string,
	aliases []string,
) error {
	dashboardId, err := r.sessions.GetSelectedDashboard(ctx)
	if err != nil {
		return err
	}

	var existedCommands []model.ChannelsCommands

	query := r.gorm.WithContext(ctx).
		Where(`"channelId" = ?`, dashboardId).
		Select([]string{"id", "name", "aliases", `"channelId"`})

	if commandId != nil {
		query = query.Not(`"id" = ?`, *commandId)
	}

	if err := query.
		Find(&existedCommands).Error; err != nil {
		return fmt.Errorf("failed to get existed commands: %w", err)
	}

	for _, existedCommand := range existedCommands {
		if existedCommand.Name == name {
			return &commandAlreadyExistsError
		}
		for _, aliase := range existedCommand.Aliases {
			if aliase == name {
				return &commandAlreadyExistsError
			}
		}

		for _, aliase := range aliases {
			if existedCommand.Name == aliase {
				return &commandAlreadyExistsError
			}

			if slices.Contains(existedCommand.Aliases, aliase) {
				return &commandAlreadyExistsError
			}
		}
	}

	return nil
}
