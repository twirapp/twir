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
		Data:              data,
		ChannelID:         channelId,
	})

	sort.Slice(operations, func(i, j int) bool {
		return operations[i].Order < operations[j].Order
	})

	for _, operation := range operations {
		for i := 0; i < operation.Repeat; i++ {
			if operation.Delay.Valid {
				duration := time.Duration(operation.Delay.Int64) * time.Second
				time.Sleep(duration)
			}

			switch operation.Type {
			case model.OperationSendMessage:
				if operation.Input.Valid {
					processor.SendMessage(channelId, operation.Input.String)
				}
			case model.OperationBan, model.OperationUnban:
				if data.UserName == "" {
					continue
				}

				processor.BanOrUnban(operation.Type)
			case model.OperationBanRandom:
				processor.BanRandom()
			case model.OperationVip, model.OperationUnvip:
				if data.UserName == "" {
					continue
				}

				processor.VipOrUnvip(operation.Type)
			case model.OperationEnableSubMode, model.OperationDisableSubMode:
				processor.SwitchSubMode(operation.Type)
			case model.OperationEnableEmoteOnly, model.OperationDisableEmoteOnly:
				processor.SwitchEmoteOnly(operation.Type)
			case model.OperationChangeTitle:
				if !operation.Input.Valid {
					continue
				}
				processor.ChangeTitle(operation.Input.String)
			case model.OperationChangeCategory:
				if !operation.Input.Valid {
					continue
				}
				processor.ChangeTitle(operation.Input.String)
			}
		}
	}
}
