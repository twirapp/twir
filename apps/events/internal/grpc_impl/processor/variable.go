package processor

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/samber/lo"
	model "github.com/satont/twir/libs/gomodels"
)

func (c *Processor) ChangeVariableValue(variableId, input string) error {
	hydratedInput, err := c.HydrateStringWithData(input, c.data)

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

	if variable.Type == model.CustomVarNumber {
		newVarValue, err := strconv.Atoi(hydratedInput)
		if err != nil {
			return InternalError
		}
		variable.Response = fmt.Sprintf("%v", newVarValue)
	} else {
		variable.Response = hydratedInput
	}

	if err = c.services.DB.Save(variable).Error; err != nil {
		return err
	}

	return nil
}

func (c *Processor) IncrementORDecrementVariable(
	operationType model.EventOperationType, variableId, input string,
) error {
	hydratedInput, err := c.HydrateStringWithData(input, c.data)

	if err != nil || len(hydratedInput) == 0 {
		return fmt.Errorf("cannot hydrate string %w", err)
	}

	hydratedInput = strings.TrimSpace(hydratedInput)

	variable := &model.ChannelsCustomvars{}
	err = c.services.DB.
		Where(`"channelId" = ? AND "id" = ? AND "type" = ?`, c.channelId, variableId, model.CustomVarNumber).
		Find(variable).
		Error
	if err != nil {
		return err
	}

	if variable.ID == "" {
		return InternalError
	}

	currentVariableNumber, err := strconv.Atoi(variable.Response)
	if err != nil {
		return InternalError
	}

	newVariableNumber, err := strconv.Atoi(hydratedInput)
	if err != nil {
		return InternalError
	}

	variable.Response = fmt.Sprintf(
		"%v",
		lo.
			If(operationType == model.OperationIncrementVariable, currentVariableNumber+newVariableNumber).
			Else(currentVariableNumber-newVariableNumber),
	)

	if err = c.services.DB.Save(variable).Error; err != nil {
		return err
	}

	return nil
}
