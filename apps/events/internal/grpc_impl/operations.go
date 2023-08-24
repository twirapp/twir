package grpc_impl

import (
	"context"
	"github.com/satont/twir/libs/grpc/generated/websockets"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/samber/lo"
	"github.com/satont/twir/apps/events/internal"
	processor_module "github.com/satont/twir/apps/events/internal/grpc_impl/processor"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/twitch"
)

// check is filter returns false and if so - break operations loop
func (c *EventsGrpcImplementation) processFilters(
	processor *processor_module.Processor,
	filters []*model.EventOperationFilter,
	data internal.Data,
) bool {
	for _, filter := range filters {
		hydratedRight, _ := processor.HydrateStringWithData(filter.Right, data)
		hydratedLeft, _ := processor.HydrateStringWithData(filter.Left, data)

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

func (c *EventsGrpcImplementation) processOperations(
	channelId string,
	event model.Event,
	data internal.Data,
) {
	streamerApiClient, err := twitch.NewUserClient(channelId, *c.services.Cfg, c.services.TokensGrpc)
	if err != nil {
		c.services.Logger.Sugar().Error(err)
		return
	}

	// won't process stream if event setted to online only streams and stream is offline
	if event.OnlineOnly {
		stream := &model.ChannelsStreams{}
		err := c.services.DB.Where(`"userId" = ?`, channelId).Find(stream).Error
		if err != nil {
			c.services.Logger.Sugar().Error(err)
			return
		}

		if stream.ID == "" {
			return
		}
	}

	processor := processor_module.NewProcessor(
		processor_module.Opts{
			Services:          c.services,
			StreamerApiClient: streamerApiClient,
			Data:              &data,
			ChannelID:         channelId,
		},
	)

	sort.Slice(
		event.Operations, func(i, j int) bool {
			return event.Operations[i].Order < event.Operations[j].Order
		},
	)

	data.PrevOperation = &internal.DataFromPrevOperation{}

	var operationError error

operationsLoop:
	for _, operation := range event.Operations {
		if !operation.Enabled {
			continue
		}

		operation.Repeat++

		allFiltersOk := c.processFilters(processor, operation.Filters, data)
		if !allFiltersOk {
			continue
		}

		for i := 0; i < operation.Repeat; i++ {
			if operation.Delay != 0 {
				duration := time.Duration(operation.Delay) * time.Second
				time.Sleep(duration)
			}

			switch operation.Type {
			case model.OperationSendMessage:
				if operation.Input.Valid {
					operationError = processor.SendMessage(
						channelId,
						operation.Input.String,
						operation.UseAnnounce,
					)
				}
			case model.OperationBan, model.OperationUnban:
				if !operation.Input.Valid {
					continue
				}

				operationError = processor.BanOrUnban(
					operation.Input.String,
					operation.Type,
					operation.TimeoutMessage,
				)
			case model.OperationTimeout:
				if operation.Input.Valid {
					operationError = processor.Timeout(
						operation.Input.String,
						operation.TimeoutTime,
						operation.TimeoutMessage,
					)
				}
			case model.OperationTimeoutRandom:
				operationError = processor.BanRandom(
					operation.TimeoutTime,
					operation.TimeoutMessage,
				)
			case model.OperationBanRandom:
				operationError = processor.BanRandom(
					0,
					operation.TimeoutMessage,
				)
			case model.OperationVip, model.OperationUnvip:
				if !operation.Input.Valid {
					continue
				}

				operationError = processor.VipOrUnvip(operation.Input.String, operation.Type)
			case model.OperationUnvipRandom, model.OperationUnvipRandomIfNoSlots:
				operationError = processor.UnvipRandom(operation.Type, operation.Input.String)
			case model.OperationEnableSubMode, model.OperationDisableSubMode:
				operationError = processor.SwitchSubMode(operation.Type)
			case model.OperationEnableEmoteOnly, model.OperationDisableEmoteOnly:
				operationError = processor.SwitchEmoteOnly(operation.Type)
			case model.OperationChangeTitle:
				if !operation.Input.Valid {
					continue
				}
				operationError = processor.ChangeTitle(operation.Input.String)
			case model.OperationChangeCategory:
				if !operation.Input.Valid {
					continue
				}
				operationError = processor.ChangeCategory(operation.Input.String)
			case model.OperationMod, model.OperationUnmod:
				if !operation.Input.Valid {
					continue
				}

				operationError = processor.ModOrUnmod(operation.Input.String, operation.Type)
			case model.OperationUnmodRandom:
				operationError = processor.UnmodRandom()
			case model.OperationObsSetScene:
				if !operation.Target.Valid {
					continue
				}

				operationError = processor.ObsSetScene(operation.Target.String)
			case model.OperationObsToggleSource:
				if !operation.Target.Valid {
					continue
				}

				operationError = processor.ObsToggleSource(operation.Target.String)
			case model.OperationObsToggleAudio:
				if !operation.Target.Valid {
					continue
				}

				operationError = processor.ObsToggleAudio(operation.Target.String)
			case model.OperationObsIncreaseVolume, model.OperationObsDecreaseVolume:
				if !operation.Input.Valid || !operation.Target.Valid {
					continue
				}

				operationError = processor.ObsAudioChangeVolume(
					operation.Type,
					operation.Target.String,
					operation.Input.String,
				)
			case model.OperationObsEnableAudio, model.OperationObsDisableAudio:
				if !operation.Target.Valid {
					return
				}

				operationError = processor.ObsEnableOrDisableAudio(operation.Type, operation.Target.String)
			case model.OperationObsSetVolume:
				if !operation.Input.Valid || !operation.Target.Valid {
					continue
				}

				operationError = processor.ObsAudioSetVolume(
					operation.Target.String,
					operation.Input.String,
				)
			case model.OperationObsStartStream, model.OperationObsStopStream:
				operationError = processor.ObsStartOrStopStream(operation.Type)
			case model.OperationChangeVariable:
				if !operation.Input.Valid || !operation.Target.Valid {
					continue
				}

				operationError = processor.ChangeVariableValue(
					operation.Target.String,
					operation.Input.String,
				)
			case model.OperationIncrementVariable, model.OperationDecrementVariable:
				if !operation.Target.Valid {
					continue
				}

				operationError = processor.IncrementORDecrementVariable(
					operation.Type,
					operation.Target.String,
					operation.Input.String,
				)
			case model.OperationTTSSay:
				if !operation.Input.Valid {
					continue
				}

				operationError = processor.TtsSay(channelId, data.UserID, operation.Input.String)
			case model.OperationTTSSkip:
				operationError = processor.TtsSkip(channelId)
			case model.OperationTTSEnable, model.OperationTTSDisable:
				action := lo.If(operation.Type == model.OperationTTSEnable, true).Else(false)

				operationError = processor.TtsChangeState(channelId, action)
			case model.OperationTTSSwitchAutoRead, model.OperationTTSEnableAutoRead, model.OperationTTSDisableAutoRead:
				value := lo.
					If[*bool](operation.Type == model.OperationTTSSwitchAutoRead, nil).
					ElseIf(operation.Type == model.OperationTTSEnableAutoRead, lo.ToPtr(true)).
					Else(lo.ToPtr(false))

				operationError = processor.TtsChangeAutoReadState(channelId, value)
			case model.OperationAllowCommandToUser, model.OperationRemoveAllowCommandToUser:
				if !operation.Input.Valid || !operation.Target.Valid {
					continue
				}

				operationError = processor.AllowOrRemoveAllowCommandToUser(
					operation.Type,
					operation.Target.String,
					operation.Input.String,
				)
			case model.OperationDenyCommandToUser, model.OperationRemoveDenyCommandToUser:
				if !operation.Input.Valid || !operation.Target.Valid {
					continue
				}

				operationError = processor.DenyOrRemoveDenyCommandToUser(
					operation.Type,
					operation.Target.String,
					operation.Input.String,
				)
			case model.OperationTriggerAlert:
				if !operation.Target.Valid {
					continue
				}
				c.services.WebsocketsGrpc.TriggerAlert(
					context.TODO(),
					&websockets.TriggerAlertRequest{
						ChannelId: channelId,
						AlertId:   operation.Target.String,
					},
				)
			}

			if operationError != nil {
				if operationError != processor_module.InternalError {
					c.services.Logger.Sugar().Error(operationError)
				}
				break operationsLoop
			}
		}
	}
}
