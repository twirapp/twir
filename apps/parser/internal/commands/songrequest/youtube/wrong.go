package sr_youtube

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
	model "tsuwari/models"
	"tsuwari/parser/internal/types"
	variables_cache "tsuwari/parser/internal/variablescache"

	"github.com/samber/lo"
	youtubenats "github.com/satont/tsuwari/libs/nats/youtube"
	"google.golang.org/protobuf/proto"
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
			result.Result = append(result.Result, `you haven't requested any song`)
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
			result.Result = append(result.Result, "seems like you provided not a number.")
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
			result.Result = append(result.Result, "cannot delete song")
			return result
		}

		natsData, err := proto.Marshal(
			&youtubenats.RemoveSongFromQueue{ChannelId: ctx.ChannelId, EntityId: choosedSong.ID},
		)
		if err == nil {
			ctx.Services.Nats.Publish(youtubenats.SUBJECTS_REMOVE_SONG_FROM_QUEUE, natsData)
		} else {
			log.Fatal(err)
			result.Result = append(result.Result, "internal error happend when we removing song from queue")
			return result
		}
		result.Result = append(
			result.Result,
			fmt.Sprintf("Song %s deleted from queue", choosedSong.Title),
		)

		return result
	},
}
