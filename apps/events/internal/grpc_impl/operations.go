package grpc_impl

import (
	"github.com/satont/tsuwari/apps/events/internal"
	processor_module "github.com/satont/tsuwari/apps/events/internal/grpc_impl/processor"
	model "github.com/satont/tsuwari/libs/gomodels"
	"github.com/satont/tsuwari/libs/twitch"
	"sort"
	"time"
)

func (c *EventsGrpcImplementation) processOperations(channelId string, operations []model.EventOperation, data internal.Data) {
	streamerApiClient, err := twitch.NewUserClient(channelId, *c.services.Cfg, c.services.TokensGrpc)
	if err != nil {
		c.services.Logger.Sugar().Error(err)
		return
	}

	processor := processor_module.NewProcessor(processor_module.Opts{
		Services:          c.services,
		StreamerApiClient: streamerApiClient,
		Data:              &data,
		ChannelID:         channelId,
	})

	sort.Slice(operations, func(i, j int) bool {
		return operations[i].Order < operations[j].Order
	})

	data.PrevOperation = &internal.DataFromPrevOperation{}

	var operationError error

operationsLoop:
	for _, operation := range operations {
		for i := 0; i < operation.Repeat; i++ {
			if operation.Delay != 0 {
				duration := time.Duration(operation.Delay) * time.Second
				time.Sleep(duration)
			}

			switch operation.Type {
			case model.OperationSendMessage:
				if operation.Input.Valid {
					operationError = processor.SendMessage(channelId, operation.Input.String, operation.UseAnnounce)
				}
			case model.OperationBan, model.OperationUnban:
				if !operation.Input.Valid {
					continue
				}

				operationError = processor.BanOrUnban(operation.Input.String, operation.Type)
			case model.OperationTimeout:
				if operation.Input.Valid {
					operationError = processor.Timeout(operation.Input.String, operation.TimeoutTime)
				}
			case model.OperationTimeoutRandom:
				operationError = processor.BanRandom(operation.TimeoutTime)
			case model.OperationBanRandom:
				operationError = processor.BanRandom(0)
			case model.OperationVip, model.OperationUnvip:
				if !operation.Input.Valid {
					continue
				}

				operationError = processor.VipOrUnvip(operation.Input.String, operation.Type)
			case model.OperationUnvipRandom:
				operationError = processor.UnvipRandom()
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
				if !operation.ObsTargetName.Valid {
					continue
				}

				operationError = processor.ObsSetScene(operation.ObsTargetName.String)
			case model.OperationObsToggleSource:
				if !operation.ObsTargetName.Valid {
					continue
				}

				operationError = processor.ObsToggleSource(operation.ObsTargetName.String)
			case model.OperationObsToggleAudio:
				if !operation.ObsTargetName.Valid {
					continue
				}

				operationError = processor.ObsToggleAudio(operation.ObsTargetName.String)
			case model.OperationObsIncreaseVolume, model.OperationObsDecreaseVolume:
				if !operation.Input.Valid || !operation.ObsTargetName.Valid {
					continue
				}

				operationError = processor.ObsAudioChangeVolume(
					operation.Type,
					operation.ObsTargetName.String,
					operation.Input.String,
				)
			case model.OperationObsEnableAudio, model.OperationObsDisableAudio:
				if !operation.ObsTargetName.Valid {
					return
				}

				operationError = processor.ObsEnableOrDisableAudio(operation.Type, operation.ObsTargetName.String)
			case model.OperationObsSetVolume:
				if !operation.Input.Valid || !operation.ObsTargetName.Valid {
					continue
				}

				operationError = processor.ObsAudioSetVolume(operation.ObsTargetName.String, operation.Input.String)
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
