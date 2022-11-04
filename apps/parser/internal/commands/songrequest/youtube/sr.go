package sr_youtube

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"time"
	model "tsuwari/models"
	"tsuwari/parser/internal/types"
	variables_cache "tsuwari/parser/internal/variablescache"

	youtubenats "github.com/satont/tsuwari/libs/nats/youtube"
	"google.golang.org/protobuf/proto"

	ytsr "github.com/SherlockYigit/youtube-go"
	ytdl "github.com/kkdai/youtube/v2"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"

	"github.com/samber/lo"
)

var (
	ytContext  = context.Background()
	YtDlClient = ytdl.Client{}
	linkRegexp = regexp.MustCompile(
		`(?m)^((?:https?:)?\/\/)?((?:www|m)\.)?((?:youtube(-nocookie)?\.com|youtu.be))(\/(?:[\w\-]+\?v=|embed\/|v\/)?)(?P<id>[\w\-]+)(\S+)?$`,
	)
)

var SrCommand = types.DefaultCommand{
	Command: types.Command{
		Name:               "ytsr",
		Description:        lo.ToPtr("Song requests from youtube"),
		Permission:         "VIEWER",
		Visible:            false,
		Module:             lo.ToPtr("SONGREQUEST"),
		IsReply:            true,
		KeepResponsesOrder: lo.ToPtr(false),
	},
	Handler: func(ctx variables_cache.ExecutionContext) *types.CommandsHandlerResult {
		result := &types.CommandsHandlerResult{}

		if ctx.Text == nil {
			result.Result = append(result.Result, "You should provide text for song request")
			return result
		}

		var songId string

		findByRegexp := linkRegexp.FindStringSubmatch(*ctx.Text)
		if len(findByRegexp) > 0 {
			songId = findByRegexp[6]
		} else {
			res := ytsr.Search(*ctx.Text, ytsr.SearchOptions{
				Type:  "video",
				Limit: 1,
			})

			if len(res) == 0 {
				result.Result = append(result.Result, "Song not found.")
				return result
			}

			songId = res[0].Video.Id
		}

		if songId == "" {
			result.Result = append(result.Result, "song not found")
			return result
		}

		err := ctx.Services.Db.
			Where(`"videoId" = ? AND "deletedAt" = null`, songId).
			First(&model.RequestedSong{}).
			Error

		if err != nil && err != gorm.ErrRecordNotFound {
			log.Fatal(err)
			result.Result = append(result.Result, "internal error")
			return result
		}

		if err == nil {
			result.Result = append(result.Result, "That song already in queue")
			return result
		}

		ytdlSongInfo, err := YtDlClient.GetVideo(
			fmt.Sprintf("https://www.youtube.com/watch?v=%s", songId),
		)
		if err != nil {
			result.Result = append(result.Result, "cannot get information about song.")
			return result
		}

		entity := model.RequestedSong{
			ID:            uuid.NewV4().String(),
			ChannelID:     ctx.ChannelId,
			OrderedById:   ctx.SenderId,
			OrderedByName: ctx.SenderName,
			VideoID:       ytdlSongInfo.ID,
			Title:         ytdlSongInfo.Title,
			Duration:      int32(ytdlSongInfo.Duration / time.Millisecond),
			CreatedAt:     time.Now(),
		}

		err = ctx.Services.Db.Create(&entity).Error

		if err != nil {
			log.Fatal(err)
			result.Result = append(result.Result, "internal error")
			return result
		}

		songsInQueue := []model.RequestedSong{}
		ctx.Services.Db.
			Where(`"channelId" = ? AND "id" != ?`, ctx.ChannelId, entity.ID).
			Select("duration").
			Order(`"createdAt" desc`).
			Find(&songsInQueue)

		timeForWait := 0 * time.Minute
		for _, s := range songsInQueue {
			timeForWait = time.Duration(s.Duration)*time.Millisecond + timeForWait
		}

		result.Result = append(
			result.Result,
			fmt.Sprintf(
				`Song "%s" requested, your position #%d. Estimated wait time before your track will be played is %s.`,
				ytdlSongInfo.Title,
				len(songsInQueue)+1,
				timeForWait.String(),
			),
		)

		natsData, err := proto.Marshal(
			&youtubenats.AddSongToQueue{ChannelId: ctx.ChannelId, EntityId: entity.ID},
		)
		if err == nil {
			ctx.Services.Nats.Publish(youtubenats.SUBJECTS_ADD_SONG_TO_QUEUE, natsData)
		}

		return result
	},
}
