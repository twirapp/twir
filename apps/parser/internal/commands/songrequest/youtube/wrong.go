package sr_youtube

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/guregu/null"
	"github.com/satont/twir/apps/parser/internal/types"

	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/grpc/generated/websockets"

	"github.com/samber/lo"
)

var WrongCommand = &types.DefaultCommand{
	ChannelsCommands: &model.ChannelsCommands{
		Name:        "sr wrong",
		Description: null.StringFrom("Delete wrong song from queue"),
		Module:      "SONGS",
		IsReply:     true,
		Visible:     true,
	},
	Handler: func(ctx context.Context, parseCtx *types.ParseContext) (
		*types.CommandsHandlerResult,
		error,
	) {
		result := &types.CommandsHandlerResult{}

		var songs []*model.RequestedSong
		err := parseCtx.Services.Gorm.WithContext(ctx).
			Where(
				`"channelId" = ? AND "orderedById" = ? AND "deletedAt" IS NULL`,
				parseCtx.Channel.ID,
				parseCtx.Sender.ID,
			).
			Limit(5).
			Order(`"createdAt" desc`).
			Find(&songs).
			Error
		if err != nil {
			return nil, &types.CommandHandlerError{
				Message: "cannot get songs from queue",
				Err:     err,
			}
		}

		if len(songs) == 0 {
			result.Result = append(result.Result, `You haven't requested any song`)
			return result, nil
		}

		number := 1

		if parseCtx.Text != nil {
			newNumber, err := strconv.Atoi(*parseCtx.Text)
			if err != nil {
				result.Result = append(result.Result, "Seems like you provided not a number.")
				return result, nil
			}
			number = newNumber
		}

		if number > len(songs)+1 || number <= 0 {
			result.Result = append(
				result.Result,
				fmt.Sprintf("there is only %v songs", len(songs)),
			)
			return result, nil
		}

		choosedSong := songs[number-1]
		choosedSong.DeletedAt = lo.ToPtr(time.Now().UTC())
		err = parseCtx.Services.Gorm.WithContext(ctx).Updates(&choosedSong).Error
		if err != nil {
			return nil, &types.CommandHandlerError{
				Message: "cannot update song",
				Err:     err,
			}
		}

		_, err = parseCtx.Services.GrpcClients.WebSockets.YoutubeRemoveSongToQueue(
			ctx,
			&websockets.YoutubeRemoveSongFromQueueRequest{
				ChannelId: parseCtx.Channel.ID,
				EntityId:  choosedSong.ID,
			},
		)
		if err != nil {
			return nil, &types.CommandHandlerError{
				Message: "cannot remove song from queue",
				Err:     err,
			}
		}

		result.Result = append(
			result.Result,
			fmt.Sprintf("Song %s deleted from queue", choosedSong.Title),
		)

		return result, nil
	},
}
