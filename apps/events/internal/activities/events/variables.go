package events

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/events/internal/shared"
	"github.com/twirapp/twir/libs/repositories/events/model"
	"go.temporal.io/sdk/activity"

	variablesmodel "github.com/twirapp/twir/libs/repositories/variables/model"
)

func (c *Activity) ChangeVariableValue(
	ctx context.Context,
	operation model.EventOperation,
	data shared.EventData,
) error {
	activity.RecordHeartbeat(ctx, nil)

	if operation.Target == nil || *operation.Target == "" {
		return fmt.Errorf("target is required for change variable value operation")
	}

	hydratedInput, hydrateErr := c.hydrator.HydrateStringWithData(
		data.ChannelID,
		*operation.Input,
		data,
	)
	if hydrateErr != nil {
		return fmt.Errorf("cannot hydrate string %w", hydrateErr)
	}

	hydratedInput = strings.TrimSpace(hydratedInput)

	parsedUuid, err := uuid.Parse(*operation.Target)
	if err != nil {
		return fmt.Errorf("cannot parse target as uuid %w", err)
	}

	variable, err := c.variablesRepository.GetByID(ctx, parsedUuid)
	if err != nil {
		return fmt.Errorf("cannot get variable by ID %w", err)
	}

	if variable.Type == variablesmodel.CustomVarNumber {
		newVarValue, atoiErr := strconv.Atoi(hydratedInput)
		if atoiErr != nil {
			return fmt.Errorf("cannot convert string to int %w", atoiErr)
		}
		variable.Response = fmt.Sprintf("%v", newVarValue)
	} else {
		variable.Response = hydratedInput
	}

	if err = c.db.WithContext(ctx).Save(variable).Error; err != nil {
		return err
	}

	return nil
}

func (c *Activity) IncrementORDecrementVariable(
	ctx context.Context,
	operation model.EventOperation,
	data shared.EventData,
) error {
	activity.RecordHeartbeat(ctx, nil)

	if operation.Target == nil || operation.Input == nil {
		return fmt.Errorf("target and input are required for increment or decrement variable operation")
	}

	hydratedInput, err := c.hydrator.HydrateStringWithData(
		data.ChannelID,
		*operation.Input,
		data,
	)
	if err != nil {
		return fmt.Errorf("cannot hydrate string %w", err)
	}

	hydratedInput = strings.TrimSpace(hydratedInput)

	parsedUuid, err := uuid.Parse(*operation.Target)
	if err != nil {
		return fmt.Errorf("cannot parse target as uuid %w", err)
	}

	variable, err := c.variablesRepository.GetByID(ctx, parsedUuid)
	if err != nil {
		return fmt.Errorf("cannot get variable by ID %w", err)
	}

	currentVariableNumber, err := strconv.Atoi(variable.Response)
	if err != nil {
		return fmt.Errorf("cannot convert string to int %w", err)
	}

	newVariableNumber, err := strconv.Atoi(hydratedInput)
	if err != nil {
		return fmt.Errorf("cannot convert string to int %w", err)
	}

	variable.Response = fmt.Sprintf(
		"%v",
		lo.
			If(
				operation.Type == model.EventOperationTypeIncrementVariable,
				currentVariableNumber+newVariableNumber,
			).
			Else(currentVariableNumber-newVariableNumber),
	)

	if err = c.db.WithContext(ctx).Save(variable).Error; err != nil {
		return err
	}

	return nil
}
