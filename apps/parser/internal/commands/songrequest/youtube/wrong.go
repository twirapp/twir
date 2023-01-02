package sr_youtube

import (
	"context"
	"fmt"
	"github.com/samber/do"
	"github.com/satont/tsuwari/apps/parser/internal/di"
	"github.com/satont/tsuwari/apps/parser/internal/types"
	variables_cache "github.com/satont/tsuwari/apps/parser/internal/variablescache"
	model "github.com/satont/tsuwari/libs/gomodels"
	"github.com/satont/tsuwari/libs/grpc/generated/websocket"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/samber/lo"
)

var WrongCommand = types.DefaultCommand{
	Command: types.Command{
		Name:               "ytsr wrong",
		Description:        lo.ToPtr("Delete wrong song from queue"),
		Permission:         "VIEWER",
		Visible:            false,
		Module:             lo.ToPtr("SONGREQUEST"),
		IsReply:            true,
		KeepResponsesOrder: lo.ToPtr(false),
	},
	Handler: func(ctx variables_cache.ExecutionContext) *types.CommandsHandlerResult {
		websocketGrpc := do.MustInvoke[websocket.WebsocketClient](di.Provider)

		result := &types.CommandsHandlerResult{}

		songs := []model.RequestedSong{}
		err := ctx.Services.Db.
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

		if ctx.Text == nil {
			mappedSongs := lo.Map(songs, func(s model.RequestedSong, index int) string {
				return fmt.Sprintf("%v. %s", index+1, s.Title)
			})
			result.Result = append(
				result.Result,
				fmt.Sprintf("Choose song number: %s", strings.Join(mappedSongs, ", ")),
			)
			return result
		}

		number, err := strconv.Atoi(*ctx.Text)
		if err != nil {
			result.Result = append(result.Result, "Seems like you provided not a number.")
			return result
		}

		if number > len(songs)+1 || number <= 0 {
			result.Result = append(
				result.Result,
				fmt.Sprintf("there is only %v songs", len(songs)),
			)
			return result
		}

		choosedSong := songs[number-1]
		choosedSong.DeletedAt = lo.ToPtr(time.Now())
		err = ctx.Services.Db.Updates(&choosedSong).Error
		if err != nil {
			result.Result = append(result.Result, "Cannot delete song")
			return result
		}

		_, err = websocketGrpc.YoutubeRemoveSongToQueue(
			context.Background(),
			&websocket.YoutubeRemoveSongFromQueueRequest{
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
