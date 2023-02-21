package processor

import (
	"fmt"
	model "github.com/satont/tsuwari/libs/gomodels"
	"strings"
)

func (c *Processor) ChangeVariableValue(variableId, input string) error {
	hydratedInput, err := hydrateStringWithData(input, c.data)

	if err != nil || len(hydratedInput) == 0 {
		return fmt.Errorf("cannot hydrate string %w", err)
	}

	hydratedInput = strings.TrimSpace(hydratedInput)

	variable := &model.ChannelsCustomvars{}
	err = c.services.DB.Where(`"channelId" = ? AND "id" = ?`, c.channelId, variableId).Find(variable).Error
	if err != nil {
		return err
	}

	if variable.ID == "" {
		return InternalError
	}

	variable.Response = hydratedInput

	if err = c.services.DB.Save(variable).Error; err != nil {
		return err
	}

	return nil
}
