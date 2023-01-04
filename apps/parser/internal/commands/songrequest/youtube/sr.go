package sr_youtube

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/samber/do"
	"github.com/satont/tsuwari/apps/parser/internal/di"
	"github.com/satont/tsuwari/libs/grpc/generated/websockets"

	"github.com/satont/tsuwari/apps/parser/internal/config/twitch"
	"github.com/satont/tsuwari/apps/parser/internal/types"
	variables_cache "github.com/satont/tsuwari/apps/parser/internal/variablescache"
	model "github.com/satont/tsuwari/libs/gomodels"

	"github.com/satont/go-helix/v2"

	ytsr "github.com/SherlockYigit/youtube-go"
	ytdl "github.com/kkdai/youtube/v2"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"

	youtube "github.com/satont/tsuwari/libs/types/types/api/modules"

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
		websocketGrpc := do.MustInvoke[websockets.WebsocketClient](di.Provider)

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
			result.Result = append(result.Result, "Song not found")
			return result
		}

		err := ctx.Services.Db.
			Where(`"videoId" = ? AND "deletedAt" IS NULL`, songId).
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
			if err.Error() == "can't bypass age restriction: embedding of this video has been disabled" {
				result.Result = append(result.Result, "Age restriction on that track.")
			} else {
				result.Result = append(result.Result, "Cannot get information about song.")
			}
			return result
		}

		if ytdlSongInfo.Duration.Seconds() == 0 {
			result.Result = append(
				result.Result,
				"Seems like that song is live, which is disallowed.",
			)
			return result
		}

		moduleSettings := &model.ChannelModulesSettings{}
		err = ctx.Services.Db.Where(`"channelId" = ?`, ctx.ChannelId).First(moduleSettings).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			fmt.Println(err)
			result.Result = append(result.Result, "internal error")
			return result
		}
		if moduleSettings.ID != "" {
			parsedSettings := &youtube.YouTubeSettings{}
			err = json.Unmarshal(moduleSettings.Settings, parsedSettings)
			if err != nil {
				fmt.Println(err)
				result.Result = append(result.Result, "internal error")
				return result
			}

			if !*parsedSettings.Enabled {
				result.Result = append(result.Result, "Song requests not enabled")
			}

			err = validate(
				ctx.ChannelId,
				ctx.SenderId,
				ctx.Services.Db,
				ctx.Services.Twitch,
				parsedSettings,
				ytdlSongInfo,
			)
			if err != nil {
				result.Result = append(result.Result, err.Error())
				return result
			}
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

		songsInQueue := []model.RequestedSong{}
		ctx.Services.Db.
			Where(
				`"channelId" = ? AND "id" != ? AND "deletedAt" IS NULL`,
				ctx.ChannelId,
				entity.ID,
			).
			Order(`"createdAt" asc`).
			Find(&songsInQueue)

		for i, s := range songsInQueue {
			s.QueuePosition = i + 1
			// ctx.Services.Db.Model(&model.RequestedSong{}).Where("id = ?", s.ID).Update("queuePosition", i+1)
			ctx.Services.Db.Save(&s)
		}

		entity.QueuePosition = len(songsInQueue) + 1

		err = ctx.Services.Db.Create(&entity).Error

		if err != nil {
			log.Fatal(err)
			result.Result = append(result.Result, "internal error")
			return result
		}

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

		websocketGrpc.YoutubeAddSongToQueue(
			context.Background(),
			&websockets.YoutubeAddSongToQueueRequest{
				ChannelId: ctx.ChannelId,
				EntityId:  entity.ID,
			},
		)

		return result
	},
}

func validate(
	channelId, userId string,
	db *gorm.DB,
	tw *twitch.Twitch,
	settings *youtube.YouTubeSettings,
	song *ytdl.Video,
) error {
	if userId != channelId {
		return nil
	}

	if len(settings.DenyList.Users) > 0 {
		_, isUserBlackListed := lo.Find(
			settings.DenyList.Users,
			func(u youtube.YouTubeDenySettingsUsers) bool {
				return u.UserID == userId
			},
		)

		if isUserBlackListed {
			return errors.New("You cannot request any songs.")
		}
	}

	if len(settings.DenyList.Channels) > 0 {
		_, isChannelBlacklisted := lo.Find(
			settings.DenyList.Channels,
			func(u youtube.YouTubeDenySettingsChannels) bool {
				return u.ID == song.ChannelID
			},
		)

		if isChannelBlacklisted {
			return errors.New("This channel is denied for requests.")
		}
	}

	if len(settings.DenyList.Songs) > 0 {
		_, isSongBlackListed := lo.Find(
			settings.DenyList.Songs,
			func(u youtube.YouTubeDenySettingsSongs) bool {
				return u.ID == song.ID
			},
		)

		if isSongBlackListed {
			return errors.New("This song is denied to request.")
		}
	}

	if *settings.AcceptOnlyWhenOnline {
		stream := &model.ChannelsStreams{}
		db.Where(`"userId" = ?`, channelId).First(stream)
		if stream.ID == "" {
			return errors.New("Requests accepted only on online streams")
		}
	}

	if settings.MaxRequests != 0 {
		var count int64
		db.Model(&model.RequestedSong{}).
			Where(`"channelId" = ? AND "deletedAt" IS NULL`, channelId).
			Count(&count)
		if count >= int64(settings.MaxRequests) {
			return errors.New("Maximum number of tracks ordered now, try later")
		}
	}

	if settings.Song.MinViews != 0 && song.Views < settings.Song.MinViews {
		return errors.New(
			fmt.Sprintf("Song haven't %v views for request", settings.Song.MinViews),
		)
	}

	songDuration := int(song.Duration.Minutes())
	if settings.Song.MaxLength != 0 && songDuration > settings.Song.MaxLength {
		return errors.New(fmt.Sprintf("Maximum length of song is %v", settings.Song.MaxLength))
	}

	if settings.Song.MinLength != 0 && songDuration < settings.Song.MinLength {
		return errors.New(fmt.Sprintf("Minimum length of song is %v", settings.Song.MinLength))
	}

	// TODO: check categories

	if settings.User.MaxRequests != 0 {
		var count int64
		db.
			Model(&model.RequestedSong{}).
			Where(`"orderedById" = ? AND "channelId" = ? AND "deletedAt" IS NULL`, userId, channelId).
			Count(&count)
		if count >= int64(settings.User.MaxRequests) {
			return errors.New("Maximum number of tracks ordered now, try later")
		}
	}

	if settings.User.MinMessages != 0 || settings.User.MinWatchTime != 0 {
		user := &model.Users{}
		db.Where("id = ?", userId).Preload("Stats").First(&user)
		if user.ID == "" {
			return errors.New(
				"There is restrictions on user, but i cannot find you in db, sorry. :(",
			)
		}

		if settings.User.MinMessages != 0 &&
			user.Stats.Messages < int32(settings.User.MinMessages) {
			return errors.New(
				fmt.Sprintf("You haven't %v messages for request song", user.Stats.Messages),
			)
		}

		watchedInMinutes := time.Duration(user.Stats.Watched) * time.Millisecond
		if settings.User.MinWatchTime != 0 &&
			int64(watchedInMinutes.Minutes()) < settings.User.MinWatchTime {

			return errors.New(
				fmt.Sprintf(
					"You haven't watched stream for %v minutes for request song",
					time.Minute*time.Duration(settings.User.MinWatchTime),
				),
			)
		}
	}

	if settings.User.MinFollowTime != 0 {
		neededDuration := time.Minute * time.Duration(settings.User.MinFollowTime)
		followReq, err := tw.Client.GetUsersFollows(&helix.UsersFollowsParams{
			FromID: userId,
			ToID:   channelId,
		})
		if err != nil {
			return errors.New("Internal error when checking follow")
		}
		if followReq.Data.Total == 0 {
			return errors.New("For request song you need to be a followed")
		}

		followDuration := time.Since(followReq.Data.Follows[0].FollowedAt)
		if followDuration.Minutes() < neededDuration.Minutes() {
			return errors.New(
				fmt.Sprintf(
					"You need to be follower at least %v minutes for request song",
					neededDuration.Minutes(),
				),
			)
		}
	}

	return nil
}
