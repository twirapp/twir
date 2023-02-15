package grpc_impl

import (
	"github.com/satont/tsuwari/apps/events/internal"
	"github.com/satont/tsuwari/apps/events/internal/grpc_impl/processor"
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

	processor := processor.NewProcessor(processor.Opts{
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
				if data.UserName == "" {
					continue
				}

				processor.BanOrUnban(operation.Type)
			case model.OperationTimeout:
				if operation.Input.Valid {
					operationError = processor.Timeout(operation.Input.String, operation.TimeoutTime)
				}
			case model.OperationTimeoutRandom:
				operationError = processor.BanRandom(operation.TimeoutTime)
			case model.OperationBanRandom:
				operationError = processor.BanRandom(0)
			case model.OperationVip, model.OperationUnvip:
				if data.UserName == "" {
					continue
				}

				operationError = processor.VipOrUnvip(operation.Type)
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
				processor.ChangeTitle(operation.Input.String)
			case model.OperationChangeCategory:
				if !operation.Input.Valid {
					continue
				}
				operationError = processor.ChangeCategory(operation.Input.String)
			case model.OperationMod, model.OperationUnmod:
				if data.UserName == "" {
					continue
				}

				operationError = processor.ModOrUnmod(operation.Type)
			case model.OperationUnmodRandom:
				operationError = processor.UnmodRandom()
			}

			if operationError != nil {
				c.services.Logger.Sugar().Error(err)
				break operationsLoop
			}
		}
	}
}
