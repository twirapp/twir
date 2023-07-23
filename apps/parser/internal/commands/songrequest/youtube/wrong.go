package sr_youtube

import (
	"context"
	"fmt"
	"log"
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
	Handler: func(ctx context.Context, parseCtx *types.ParseContext) *types.CommandsHandlerResult {
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
			log.Fatal(err)
			result.Result = append(result.Result, "internal error")
			return result
		}

		if len(songs) == 0 {
			result.Result = append(result.Result, `You haven't requested any song`)
			return result
		}

		number := 1

		if parseCtx.Text != nil {
			newNumber, err := strconv.Atoi(*parseCtx.Text)
			if err != nil {
				result.Result = append(result.Result, "Seems like you provided not a number.")
				return result
			}
			number = newNumber
		}

		if number > len(songs)+1 || number <= 0 {
			result.Result = append(
				result.Result,
				fmt.Sprintf("there is only %v songs", len(songs)),
			)
			return result
		}

		choosedSong := songs[number-1]
		choosedSong.DeletedAt = lo.ToPtr(time.Now().UTC())
		err = parseCtx.Services.Gorm.WithContext(ctx).Updates(&choosedSong).Error
		if err != nil {
			result.Result = append(result.Result, "Cannot delete song")
			return result
		}

		_, err = parseCtx.Services.GrpcClients.WebSockets.YoutubeRemoveSongToQueue(
			ctx,
			&websockets.YoutubeRemoveSongFromQueueRequest{
				ChannelId: parseCtx.Channel.ID,
				EntityId:  choosedSong.ID,
			},
		)
		if err != nil {
			log.Fatal(err)
			result.Result = append(result.Result, "Internal error happened when we removing song from queue")
			return result
		}

		result.Result = append(
			result.Result,
			fmt.Sprintf("Song %s deleted from queue", choosedSong.Title),
		)

		return result
	},
}
