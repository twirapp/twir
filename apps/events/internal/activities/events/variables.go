package events

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/samber/lo"
	"github.com/satont/twir/apps/events/internal/shared"
	model "github.com/satont/twir/libs/gomodels"
	"go.temporal.io/sdk/activity"
)

func (c *Activity) ChangeVariableValue(
	ctx context.Context,
	operation model.EventOperation,
	data shared.EventData,
) error {
	activity.RecordHeartbeat(ctx, nil)

	if operation.Target.String == "" {
		return fmt.Errorf("target is empty")
	}
	hydratedInput, hydrateErr := c.hydrator.HydrateStringWithData(
		data.ChannelID,
		operation.Input.String,
		data,
	)
	if hydrateErr != nil {
		return fmt.Errorf("cannot hydrate string %w", hydrateErr)
	}

	hydratedInput = strings.TrimSpace(hydratedInput)

	variable := &model.ChannelsCustomvars{}
	err := c.db.WithContext(ctx).Where(
		`"channelId" = ? AND "id" = ?`,
		data.ChannelID,
		operation.Target.String,
	).First(variable).Error
	if err != nil {
		return err
	}

	if variable.Type == model.CustomVarNumber {
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

	if operation.Target.String == "" {
		return fmt.Errorf("target is empty")
	}
	hydratedInput, err := c.hydrator.HydrateStringWithData(
		data.ChannelID,
		operation.Input.String,
		data,
	)
	if err != nil {
		return fmt.Errorf("cannot hydrate string %w", err)
	}

	hydratedInput = strings.TrimSpace(hydratedInput)

	variable := &model.ChannelsCustomvars{}
	err = c.db.
		WithContext(ctx).
		Where(
			`"channelId" = ? AND "id" = ? AND "type" = ?`,
			data.ChannelID,
			operation.Target.String,
			model.CustomVarNumber,
		).
		First(variable).
		Error
	if err != nil {
		return err
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
				operation.Type == model.OperationIncrementVariable,
				currentVariableNumber+newVariableNumber,
			).
			Else(currentVariableNumber-newVariableNumber),
	)

	if err = c.db.WithContext(ctx).Save(variable).Error; err != nil {
		return err
	}

	return nil
}
