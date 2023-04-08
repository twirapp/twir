package sr_youtube

import (
	"context"
	"fmt"
	"github.com/guregu/null"
	"gorm.io/gorm"
	"log"
	"strconv"
	"time"

	"github.com/samber/do"
	"github.com/satont/tsuwari/apps/parser/internal/di"
	"github.com/satont/tsuwari/apps/parser/internal/types"
	variables_cache "github.com/satont/tsuwari/apps/parser/internal/variablescache"
	model "github.com/satont/tsuwari/libs/gomodels"
	"github.com/satont/tsuwari/libs/grpc/generated/websockets"

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
	Handler: func(ctx *variables_cache.ExecutionContext) *types.CommandsHandlerResult {
		db := do.MustInvoke[gorm.DB](di.Provider)
		websocketGrpc := do.MustInvoke[websockets.WebsocketClient](di.Provider)

		result := &types.CommandsHandlerResult{}

		songs := []model.RequestedSong{}
		err := db.
			Where(
				`"channelId" = ? AND "orderedById" = ? AND "deletedAt" IS NULL`,
				ctx.ChannelId,
				ctx.SenderId,
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

		if ctx.Text != nil {
			newNumber, err := strconv.Atoi(*ctx.Text)
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
		err = db.Updates(&choosedSong).Error
		if err != nil {
			result.Result = append(result.Result, "Cannot delete song")
			return result
		}

		_, err = websocketGrpc.YoutubeRemoveSongToQueue(
			context.Background(),
			&websockets.YoutubeRemoveSongFromQueueRequest{
				ChannelId: ctx.ChannelId,
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
