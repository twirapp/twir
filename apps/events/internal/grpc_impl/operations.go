package grpc_impl

import (
	"github.com/samber/lo"
	"github.com/satont/tsuwari/apps/events/internal"
	processor_module "github.com/satont/tsuwari/apps/events/internal/grpc_impl/processor"
	model "github.com/satont/tsuwari/libs/gomodels"
	"github.com/satont/tsuwari/libs/twitch"
	"sort"
	"time"
)

func (c *EventsGrpcImplementation) processOperations(channelId string, event model.Event, data internal.Data) {
	streamerApiClient, err := twitch.NewUserClient(channelId, *c.services.Cfg, c.services.TokensGrpc)
	if err != nil {
		c.services.Logger.Sugar().Error(err)
		return
	}

	// wont process stream if event setted to online only streams and stream is offline
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

	processor := processor_module.NewProcessor(processor_module.Opts{
		Services:          c.services,
		StreamerApiClient: streamerApiClient,
		Data:              &data,
		ChannelID:         channelId,
	})

	sort.Slice(event.Operations, func(i, j int) bool {
		return event.Operations[i].Order < event.Operations[j].Order
	})

	data.PrevOperation = &internal.DataFromPrevOperation{}

	var operationError error

operationsLoop:
	for _, operation := range event.Operations {
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

				operationError = processor.ObsAudioSetVolume(operation.Target.String, operation.Input.String)
			case model.OperationObsStartStream, model.OperationObsStopStream:
				operationError = processor.ObsStartOrStopStream(operation.Type)
			case model.OperationChangeVariable:
				if !operation.Input.Valid || !operation.Target.Valid {
					continue
				}

				operationError = processor.ChangeVariableValue(operation.Target.String, operation.Input.String)
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
