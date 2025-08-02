package workflows

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/goccy/go-json"
	"github.com/twirapp/twir/apps/events/internal/shared"
	deprecatedmodel "github.com/twirapp/twir/libs/gomodels"
	channelmodel "github.com/twirapp/twir/libs/repositories/channels/model"
	"github.com/twirapp/twir/libs/repositories/events/model"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

func (c *EventWorkflow) Flow(
	ctx workflow.Context,
	eventType model.EventType,
	data shared.EventData,
) error {
	options := workflow.ActivityOptions{
		StartToCloseTimeout: time.Second * 15,
		HeartbeatTimeout:    time.Second * 5,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:        time.Second,
			BackoffCoefficient:     2.0,
			MaximumInterval:        time.Second * 100,
			MaximumAttempts:        3,
			NonRetryableErrorTypes: []string{},
		},
		TaskQueue: shared.EventsWorkerTaskQueueName,
	}
	ctx = workflow.WithActivityOptions(ctx, options)

	if data.ChannelID == "" {
		return errors.New("channel id is empty")
	}

	eventsCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	channelEvents, err := c.channelsEventsWithOperationsCache.Get(eventsCtx, data.ChannelID)
	if err != nil {
		return err
	}

	var stream deprecatedmodel.ChannelsStreams
	err = c.db.
		Where(`"userId" = ?`, data.ChannelID).
		Find(&stream).Error
	if err != nil {
		return err
	}

	channel, err := c.channelsCache.Get(eventsCtx, data.ChannelID)
	if err != nil {
		return err
	}

	if channel == channelmodel.Nil {
		return errors.New("channel not found")
	}

	var operations []model.EventOperation

	for _, entity := range channelEvents {
		if !entity.Enabled {
			continue
		}

		if entity.ID == "" {
			continue
		}

		if entity.OnlineOnly && stream.ID == "" {
			continue
		}

		if entity.Type != eventType {
			continue
		}

		if !channel.IsEnabled || channel.IsTwitchBanned {
			continue
		}

		if entity.Type == model.EventTypeCommandUsed &&
			data.CommandID != "" &&
			entity.CommandID != nil &&
			data.CommandID != *entity.CommandID {
			continue
		}

		if entity.Type == model.EventTypeRedemptionCreated &&
			data.RewardID != "" &&
			entity.RewardID != nil &&
			data.RewardID != *entity.RewardID {
			continue
		}

		if entity.Type == model.EventTypeKeywordMatched &&
			data.KeywordID != "" &&
			entity.KeywordID != nil &&
			data.KeywordID != *entity.KeywordID {
			continue
		}

		operations = append(operations, entity.Operations...)
	}

	workflow.GetLogger(ctx).Info("Scheduled workflow")

	// set workflow execution state
	info := workflow.GetInfo(ctx)
	if info.WorkflowExecution.ID == "" {
		return errors.New("workflow id is empty")
	}
	bytes, bytesErr := json.Marshal(&shared.EventsWorkflowExecutionState{})
	if bytesErr != nil {
		return bytesErr
	}

	redisErr := c.redis.Set(
		context.Background(),
		"events:workflows:"+info.WorkflowExecution.ID,
		bytes,
		7*24*time.Hour,
	).Err()
	if redisErr != nil {
		return redisErr
	}
	// end set workflow execution state

	workflow.GetLogger(ctx).Info("Got operations", "size", len(operations))

	// execute event operations
	for _, operation := range operations {
		if !operation.Enabled {
			continue
		}
		operation.Repeat++

		filtersOk := c.filtersOk(data.ChannelID, operation.Filters, data)
		if !filtersOk {
			continue
		}

		for i := 0; i < operation.Repeat; i++ {
			if operation.Delay != 0 {
				workflow.Sleep(ctx, time.Duration(operation.Delay)*time.Second)
			}
			var operationErr error
			switch operation.Type {
			case model.EventOperationTypeSendMessage:
				operationErr = workflow.ExecuteActivity(
					ctx,
					c.eventsActivity.SendMessage,
					operation,
					data,
				).Get(ctx, nil)
			case model.EventOperationTypeBan, model.EventOperationTypeTimeout:
				operationErr = workflow.ExecuteActivity(
					ctx,
					c.eventsActivity.Ban,
					operation,
					data,
				).Get(ctx, nil)
			case model.EventOperationTypeBanRandom, model.EventOperationTypeTimeoutRandom:
				operationErr = workflow.ExecuteActivity(
					ctx,
					c.eventsActivity.BanRandom,
					operation,
					data,
				).Get(ctx, nil)
			case model.EventOperationTypeUnban:
				operationErr = workflow.ExecuteActivity(
					ctx,
					c.eventsActivity.Unban,
					operation,
					data,
				).Get(ctx, nil)
			case model.EventOperationTypeChangeTitle:
				operationErr = workflow.ExecuteActivity(
					ctx,
					c.eventsActivity.ChangeTitle,
					operation,
					data,
				).Get(ctx, nil)
			case model.EventOperationTypeChangeCategory:
				operationErr = workflow.ExecuteActivity(
					ctx,
					c.eventsActivity.ChangeCategory,
					operation,
					data,
				).Get(ctx, nil)
			case model.EventOperationTypeAllowCommandToUser, model.EventOperationTypeRemoveAllowCommandToUser:
				operationErr = workflow.ExecuteActivity(
					ctx,
					c.eventsActivity.CommandAllowOrRemoveUserPermission,
					operation,
					data,
				).Get(ctx, nil)
			case model.EventOperationTypeDenyCommandToUser, model.EventOperationTypeRemoveDenyCommandToUser:
				operationErr = workflow.ExecuteActivity(
					ctx,
					c.eventsActivity.CommandDenyOrRemoveUserPermission,
					operation,
					data,
				).Get(ctx, nil)
			case model.EventOperationTypeCreateGreeting:
				operationErr = workflow.ExecuteActivity(
					ctx,
					c.eventsActivity.CreateGreeting,
					operation,
					data,
				).Get(ctx, nil)
			case model.EventOperationTypeDisableEmoteOnly, model.EventOperationTypeEnableEmoteOnly:
				operationErr = workflow.ExecuteActivity(
					ctx,
					c.eventsActivity.SwitchEmoteOnly,
					operation,
					data,
				).Get(ctx, nil)
			case model.EventOperationTypeMod, model.EventOperationTypeUnmod:
				operationErr = workflow.ExecuteActivity(
					ctx,
					c.eventsActivity.ModOrUnmod,
					operation,
					data,
				).Get(ctx, nil)
			case model.EventOperationTypeUnmodRandom:
				operationErr = workflow.ExecuteActivity(
					ctx,
					c.eventsActivity.UnmodRandom,
					operation,
					data,
				).Get(ctx, nil)
			case model.EventOperationTypeObsChangeScene:
				operationErr = workflow.ExecuteActivity(
					ctx,
					c.eventsActivity.ObsSetScene,
					operation,
					data,
				).Get(ctx, nil)
			case model.EventOperationTypeObsToggleSource:
				operationErr = workflow.ExecuteActivity(
					ctx,
					c.eventsActivity.ObsToggleSource,
					operation,
					data,
				).Get(ctx, nil)
			case model.EventOperationTypeObsToggleAudio:
				operationErr = workflow.ExecuteActivity(
					ctx,
					c.eventsActivity.ObsToggleAudio,
					operation,
					data,
				).Get(ctx, nil)
			case model.EventOperationTypeObsIncreaseAudioVolume, model.EventOperationTypeObsDecreaseAudioVolume:
				operationErr = workflow.ExecuteActivity(
					ctx,
					c.eventsActivity.ObsAudioChangeVolume,
					operation,
					data,
				).Get(ctx, nil)
			case model.EventOperationTypeObsSetAudioVolume:
				operationErr = workflow.ExecuteActivity(
					ctx,
					c.eventsActivity.ObsAudioSetVolume,
					operation,
					data,
				).Get(ctx, nil)
			case model.EventOperationTypeObsEnableAudio, model.EventOperationTypeObsDisableAudio:
				operationErr = workflow.ExecuteActivity(
					ctx,
					c.eventsActivity.ObsEnableOrDisableAudio,
					operation,
					data,
				).Get(ctx, nil)
			case model.EventOperationTypeObsStartStream, model.EventOperationTypeObsStopStream:
				operationErr = workflow.ExecuteActivity(
					ctx,
					c.eventsActivity.ObsStartOrStopStream,
					operation,
					data,
				).Get(ctx, nil)
			case model.EventOperationTypeDisableSubmode, model.EventOperationTypeEnableSubmode:
				operationErr = workflow.ExecuteActivity(
					ctx,
					c.eventsActivity.SwitchSubMode,
					operation,
					data,
				).Get(ctx, nil)
			case model.EventOperationTypeTtsSay:
				operationErr = workflow.ExecuteActivity(
					ctx,
					c.eventsActivity.TtsSay,
					operation,
					data,
				).Get(ctx, nil)
			case model.EventOperationTypeTtsEnable, model.EventOperationTypeTtsDisable:
				operationErr = workflow.ExecuteActivity(
					ctx,
					c.eventsActivity.TtsChangeState,
					operation,
					data,
				).Get(ctx, nil)
			case model.EventOperationTypeTtsDisableAutoread, model.EventOperationTypeTtsEnableAutoread, model.EventOperationTypeTtsSwitchAutoread:
				operationErr = workflow.ExecuteActivity(
					ctx,
					c.eventsActivity.TtsChangeAutoReadState,
					operation,
					data,
				).Get(ctx, nil)
			case model.EventOperationTypeTtsSkip:
				operationErr = workflow.ExecuteActivity(
					ctx,
					c.eventsActivity.TtsSkip,
					operation,
					data,
				).Get(ctx, nil)
			case model.EventOperationTypeChangeVariable:
				operationErr = workflow.ExecuteActivity(
					ctx,
					c.eventsActivity.ChangeVariableValue,
					operation,
					data,
				).Get(ctx, nil)
			case model.EventOperationTypeIncrementVariable, model.EventOperationTypeDecrementVariable:
				operationErr = workflow.ExecuteActivity(
					ctx,
					c.eventsActivity.IncrementORDecrementVariable,
					operation,
					data,
				).Get(ctx, nil)
			case model.EventOperationTypeVip, model.EventOperationTypeUnvip:
				operationErr = workflow.ExecuteActivity(
					ctx,
					c.eventsActivity.VipOrUnvip,
					operation,
					data,
				).Get(ctx, nil)
			case model.EventOperationTypeUnvipRandom, model.EventOperationTypeUnvipRandomIfNoSlots:
				operationErr = workflow.ExecuteActivity(
					ctx,
					c.eventsActivity.UnvipRandom,
					operation,
					data,
				).Get(ctx, nil)
			case model.EventOperationTypeSeventvAddEmote, model.EventOperationTypeSeventvRemoveEmote:
				operationErr = workflow.ExecuteActivity(
					ctx,
					c.eventsActivity.SevenTvEmoteManage,
					operation,
					data,
				).Get(ctx, nil)
			case model.EventOperationTypeRaidChannel:
				operationErr = workflow.ExecuteActivity(
					ctx,
					c.eventsActivity.RaidChannel,
					operation,
					data,
				).Get(ctx, nil)
			case model.EventOperationTypeTriggerAlert:
				operationErr = workflow.ExecuteActivity(
					ctx,
					c.eventsActivity.TriggerAlert,
					operation,
					data,
				).Get(ctx, nil)
			case model.EventOperationTypeShoutoutChannel:
				operationErr = workflow.ExecuteActivity(
					ctx,
					c.eventsActivity.ShoutoutChannel,
					operation,
					data,
				).Get(ctx, nil)
			case model.EventOperationTypeMessageDelete:
				operationErr = workflow.ExecuteActivity(
					ctx,
					c.eventsActivity.MessageDelete,
					operation,
					data,
				).Get(ctx, nil)
			}

			if operationErr != nil {
				return operationErr
			}
		}
	}

	return nil
}

func (c *EventWorkflow) filtersOk(
	channelId string,
	filters []model.EventOperationFilter,
	data shared.EventData,
) bool {
	for _, filter := range filters {
		hydratedRight, _ := c.hydrator.HydrateStringWithData(channelId, filter.Right, data)
		hydratedLeft, _ := c.hydrator.HydrateStringWithData(channelId, filter.Left, data)

		numericRight, _ := strconv.Atoi(hydratedRight)
		numericLeft, _ := strconv.Atoi(hydratedLeft)

		if filter.Type == model.EventOperationFilterTypeEquals {
			if hydratedLeft != hydratedRight {
				return false
			}
		}

		if filter.Type == model.EventOperationFilterTypeNotEquals {
			if hydratedLeft == hydratedRight {
				return false
			}
		}

		if filter.Type == model.EventOperationFilterTypeContains {
			if !strings.Contains(hydratedLeft, hydratedRight) {
				return false
			}
		}

		if filter.Type == model.EventOperationFilterTypeStartsWith {
			if !strings.HasPrefix(hydratedLeft, hydratedRight) {
				return false
			}
		}

		if filter.Type == model.EventOperationFilterTypeEndsWith {
			if !strings.HasSuffix(hydratedLeft, hydratedRight) {
				return false
			}
		}

		if filter.Type == model.EventOperationFilterTypeNotContains {
			if strings.Contains(hydratedLeft, hydratedRight) {
				return false
			}
		}

		if filter.Type == model.EventOperationFilterTypeGreaterThan {
			if numericLeft <= numericRight {
				return false
			}
		}

		if filter.Type == model.EventOperationFilterTypeLessThan {
			if numericLeft >= numericRight {
				return false
			}
		}

		if filter.Type == model.EventOperationFilterTypeGreaterThanOrEquals {
			if numericLeft < numericRight {
				return false
			}
		}

		if filter.Type == model.EventOperationFilterTypeLessThanOrEquals {
			if numericLeft > numericRight {
				return false
			}
		}

		if filter.Type == model.EventOperationFilterTypeIsEmpty {
			if hydratedLeft != "" {
				return false
			}
		}

		if filter.Type == model.EventOperationFilterTypeIsNotEmpty {
			if hydratedLeft == "" {
				return false
			}
		}
	}

	return true
}
