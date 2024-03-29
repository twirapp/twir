package workflows

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/goccy/go-json"
	"github.com/satont/twir/apps/events/internal/shared"
	model "github.com/satont/twir/libs/gomodels"
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

	var channelEvents []model.Event
	err := c.db.
		Where(`"channelId" = ? AND "type" = ? AND "enabled" = ?`, data.ChannelID, eventType, true).
		Preload("Channel").
		Preload("Operations").
		Preload("Operations.Filters").
		Find(&channelEvents).
		Error
	if err != nil {
		return err
	}

	var stream model.ChannelsStreams
	err = c.db.
		Where(`"userId" = ?`, data.ChannelID).
		Find(&stream).Error
	if err != nil {
		return err
	}

	var operations []model.EventOperation

	for _, entity := range channelEvents {
		if entity.ID == "" {
			continue
		}

		if entity.OnlineOnly && stream.ID == "" {
			continue
		}

		if entity.Channel != nil && (!entity.Channel.IsEnabled || entity.Channel.IsBanned) {
			continue
		}

		if entity.Type == model.EventTypeCommandUsed &&
			data.CommandID != "" &&
			entity.CommandID.Valid &&
			data.CommandID != entity.CommandID.String {
			continue
		}

		if entity.Type == model.EventTypeRedemptionCreated &&
			data.RewardID != "" &&
			entity.RewardID.Valid &&
			data.RewardID != entity.RewardID.String {
			continue
		}

		if entity.Type == model.EventTypeKeywordMatched &&
			data.KeywordID != "" &&
			entity.KeywordID.Valid &&
			data.KeywordID != entity.KeywordID.String {
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
			case model.OperationSendMessage:
				operationErr = workflow.ExecuteActivity(
					ctx,
					c.eventsActivity.SendMessage,
					operation,
					data,
				).Get(ctx, nil)
			case model.OperationBan, model.OperationTimeout:
				operationErr = workflow.ExecuteActivity(
					ctx,
					c.eventsActivity.Ban,
					operation,
					data,
				).Get(ctx, nil)
			case model.OperationBanRandom, model.OperationTimeoutRandom:
				operationErr = workflow.ExecuteActivity(
					ctx,
					c.eventsActivity.BanRandom,
					operation,
					data,
				).Get(ctx, nil)
			case model.OperationUnban:
				operationErr = workflow.ExecuteActivity(
					ctx,
					c.eventsActivity.Unban,
					operation,
					data,
				).Get(ctx, nil)
			case model.OperationChangeTitle:
				operationErr = workflow.ExecuteActivity(
					ctx,
					c.eventsActivity.ChangeTitle,
					operation,
					data,
				).Get(ctx, nil)
			case model.OperationChangeCategory:
				operationErr = workflow.ExecuteActivity(
					ctx,
					c.eventsActivity.ChangeCategory,
					operation,
					data,
				).Get(ctx, nil)
			case model.OperationAllowCommandToUser, model.OperationRemoveAllowCommandToUser:
				operationErr = workflow.ExecuteActivity(
					ctx,
					c.eventsActivity.CommandAllowOrRemoveUserPermission,
					operation,
					data,
				).Get(ctx, nil)
			case model.OperationDenyCommandToUser, model.OperationRemoveDenyCommandToUser:
				operationErr = workflow.ExecuteActivity(
					ctx,
					c.eventsActivity.CommandDenyOrRemoveUserPermission,
					operation,
					data,
				).Get(ctx, nil)
			case model.OperationCreateGreeting:
				operationErr = workflow.ExecuteActivity(
					ctx,
					c.eventsActivity.CreateGreeting,
					operation,
					data,
				).Get(ctx, nil)
			case model.OperationEnableEmoteOnly, model.OperationDisableEmoteOnly:
				operationErr = workflow.ExecuteActivity(
					ctx,
					c.eventsActivity.SwitchEmoteOnly,
					operation,
					data,
				).Get(ctx, nil)
			case model.OperationMod, model.OperationUnmod:
				operationErr = workflow.ExecuteActivity(
					ctx,
					c.eventsActivity.ModOrUnmod,
					operation,
					data,
				).Get(ctx, nil)
			case model.OperationUnmodRandom:
				operationErr = workflow.ExecuteActivity(
					ctx,
					c.eventsActivity.UnmodRandom,
					operation,
					data,
				).Get(ctx, nil)
			case model.OperationObsSetScene:
				operationErr = workflow.ExecuteActivity(
					ctx,
					c.eventsActivity.ObsSetScene,
					operation,
					data,
				).Get(ctx, nil)
			case model.OperationObsToggleSource:
				operationErr = workflow.ExecuteActivity(
					ctx,
					c.eventsActivity.ObsToggleSource,
					operation,
					data,
				).Get(ctx, nil)
			case model.OperationObsToggleAudio:
				operationErr = workflow.ExecuteActivity(
					ctx,
					c.eventsActivity.ObsToggleAudio,
					operation,
					data,
				).Get(ctx, nil)
			case model.OperationObsIncreaseVolume, model.OperationObsDecreaseVolume:
				operationErr = workflow.ExecuteActivity(
					ctx,
					c.eventsActivity.ObsAudioChangeVolume,
					operation,
					data,
				).Get(ctx, nil)
			case model.OperationObsSetVolume:
				operationErr = workflow.ExecuteActivity(
					ctx,
					c.eventsActivity.ObsAudioSetVolume,
					operation,
					data,
				).Get(ctx, nil)
			case model.OperationObsEnableAudio, model.OperationObsDisableAudio:
				operationErr = workflow.ExecuteActivity(
					ctx,
					c.eventsActivity.ObsEnableOrDisableAudio,
					operation,
					data,
				).Get(ctx, nil)
			case model.OperationObsStartStream, model.OperationObsStopStream:
				operationErr = workflow.ExecuteActivity(
					ctx,
					c.eventsActivity.ObsStartOrStopStream,
					operation,
					data,
				).Get(ctx, nil)
			case model.OperationEnableSubMode, model.OperationDisableSubMode:
				operationErr = workflow.ExecuteActivity(
					ctx,
					c.eventsActivity.SwitchSubMode,
					operation,
					data,
				).Get(ctx, nil)
			case model.OperationTTSSay:
				operationErr = workflow.ExecuteActivity(
					ctx,
					c.eventsActivity.TtsSay,
					operation,
					data,
				).Get(ctx, nil)
			case model.OperationTTSEnable, model.OperationTTSDisable:
				operationErr = workflow.ExecuteActivity(
					ctx,
					c.eventsActivity.TtsChangeState,
					operation,
					data,
				).Get(ctx, nil)
			case model.OperationTTSSwitchAutoRead, model.OperationTTSEnableAutoRead, model.OperationTTSDisableAutoRead:
				operationErr = workflow.ExecuteActivity(
					ctx,
					c.eventsActivity.TtsChangeAutoReadState,
					operation,
					data,
				).Get(ctx, nil)
			case model.OperationTTSSkip:
				operationErr = workflow.ExecuteActivity(
					ctx,
					c.eventsActivity.TtsSkip,
					operation,
					data,
				).Get(ctx, nil)
			case model.OperationChangeVariable:
				operationErr = workflow.ExecuteActivity(
					ctx,
					c.eventsActivity.ChangeVariableValue,
					operation,
					data,
				).Get(ctx, nil)
			case model.OperationIncrementVariable, model.OperationDecrementVariable:
				operationErr = workflow.ExecuteActivity(
					ctx,
					c.eventsActivity.IncrementORDecrementVariable,
					operation,
					data,
				).Get(ctx, nil)
			case model.OperationVip, model.OperationUnvip:
				operationErr = workflow.ExecuteActivity(
					ctx,
					c.eventsActivity.VipOrUnvip,
					operation,
					data,
				).Get(ctx, nil)
			case model.OperationUnvipRandom, model.OperationUnvipRandomIfNoSlots:
				operationErr = workflow.ExecuteActivity(
					ctx,
					c.eventsActivity.UnvipRandom,
					operation,
					data,
				).Get(ctx, nil)
			case model.OperationSevenTvAddEmote, model.OperationSevenTvRemoveEmote:
				operationErr = workflow.ExecuteActivity(
					ctx,
					c.eventsActivity.SevenTvEmoteManage,
					operation,
					data,
				).Get(ctx, nil)
			case model.OperationRaidChannel:
				operationErr = workflow.ExecuteActivity(
					ctx,
					c.eventsActivity.RaidChannel,
					operation,
					data,
				).Get(ctx, nil)
			case model.OperationTriggerAlert:
				operationErr = workflow.ExecuteActivity(
					ctx,
					c.eventsActivity.TriggerAlert,
					operation,
					data,
				).Get(ctx, nil)
			case model.OperationShoutoutChannel:
				operationErr = workflow.ExecuteActivity(
					ctx,
					c.eventsActivity.ShoutoutChannel,
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
	filters []*model.EventOperationFilter,
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
