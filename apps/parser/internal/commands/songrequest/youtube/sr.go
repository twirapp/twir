package sr_youtube

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/guregu/null"
	"github.com/satont/tsuwari/apps/parser/internal/config/twitch"
	"github.com/valyala/fasttemplate"
	"go.uber.org/zap"
	"log"
	"regexp"
	"strconv"
	"time"

	"github.com/samber/do"
	"github.com/satont/tsuwari/apps/parser/internal/di"
	"github.com/satont/tsuwari/libs/grpc/generated/websockets"

	"github.com/satont/tsuwari/apps/parser/internal/types"
	variables_cache "github.com/satont/tsuwari/apps/parser/internal/variablescache"
	model "github.com/satont/tsuwari/libs/gomodels"

	"github.com/satont/go-helix/v2"

	ytsr "github.com/SherlockYigit/youtube-go"
	ytdl "github.com/kkdai/youtube/v2"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"

	"github.com/samber/lo"
	youtube "github.com/satont/tsuwari/libs/types/types/api/modules"
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
		db := do.MustInvoke[gorm.DB](di.Provider)
		websocketGrpc := do.MustInvoke[websockets.WebsocketClient](di.Provider)
		logger := do.MustInvoke[zap.Logger](di.Provider)

		result := &types.CommandsHandlerResult{}

		if ctx.Text == nil {
			result.Result = append(result.Result, "You should provide text for song request")
			return result
		}
		var songId string

		moduleSettings := &model.ChannelModulesSettings{}
		parsedSettings := &youtube.YouTubeSettings{}
		err := db.
			Where(`"channelId" = ? AND "type" = ?`, ctx.ChannelId, "youtube_song_requests").
			First(moduleSettings).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			logger.Sugar().Error(err)
			result.Result = append(result.Result, "internal error")
			return result
		}
		if moduleSettings.ID != "" {
			err = json.Unmarshal(moduleSettings.Settings, parsedSettings)
			if err != nil {
				fmt.Println(err)
				result.Result = append(result.Result, "internal error")
				return result
			}

			if !*parsedSettings.Enabled {
				result.Result = append(result.Result, parsedSettings.Translations.NotEnabled)
				return result
			}
		} else {
			result.Result = append(result.Result, "Song requests not enabled")
			return result
		}

		findByRegexp := linkRegexp.FindStringSubmatch(*ctx.Text)
		if len(findByRegexp) > 0 {
			songId = findByRegexp[6]
		} else {
			res, err := ytsr.Search(*ctx.Text, ytsr.SearchOptions{
				Type:  "video",
				Limit: 1,
			})

			if err != nil {
				logger.Sugar().Error(err)
				result.Result = append(result.Result, parsedSettings.Translations.Song.NotFound)
				return result
			}

			if len(res) == 0 {
				result.Result = append(result.Result, parsedSettings.Translations.Song.NotFound)
				return result
			}

			songId = res[0].Video.Id
		}

		if songId == "" {
			result.Result = append(result.Result, parsedSettings.Translations.Song.NotFound)
			return result
		}

		alreadyRequestedSong := &model.RequestedSong{}
		db.Where(`"videoId" = ? AND "deletedAt" IS NULL AND "channelId" = ?`, songId, ctx.ChannelId).
			First(&alreadyRequestedSong)

		if alreadyRequestedSong.ID != "" {
			result.Result = append(result.Result, parsedSettings.Translations.Song.AlreadyInQueue)
			return result
		}

		ytdlSongInfo, err := YtDlClient.GetVideo(
			fmt.Sprintf("https://www.youtube.com/watch?v=%s", songId),
		)
		if err != nil {
			if err.Error() == "can't bypass age restriction: embedding of this video has been disabled" {
				result.Result = append(result.Result, parsedSettings.Translations.Song.AgeRestrictions)
			} else {
				result.Result = append(result.Result, parsedSettings.Translations.Song.CannotGetInformation)
			}
			return result
		}

		if ytdlSongInfo.Duration.Seconds() == 0 {
			result.Result = append(
				result.Result,
				parsedSettings.Translations.Song.Live,
			)
			return result
		}

		err = validate(
			ctx.ChannelId,
			ctx.SenderId,
			parsedSettings,
			ytdlSongInfo,
		)
		if err != nil {
			result.Result = append(result.Result, err.Error())
			return result
		}

		entity := model.RequestedSong{
			ID:                   uuid.NewV4().String(),
			ChannelID:            ctx.ChannelId,
			OrderedById:          ctx.SenderId,
			OrderedByName:        ctx.SenderName,
			OrderedByDisplayName: null.StringFrom(ctx.SenderDisplayName),
			VideoID:              ytdlSongInfo.ID,
			Title:                ytdlSongInfo.Title,
			Duration:             int32(ytdlSongInfo.Duration / time.Millisecond),
			CreatedAt:            time.Now().UTC(),
		}

		songsInQueue := []model.RequestedSong{}
		db.
			Where(
				`"channelId" = ? AND "id" != ? AND "deletedAt" IS NULL`,
				ctx.ChannelId,
				entity.ID,
			).
			Order(`"createdAt" asc`).
			Find(&songsInQueue)

		for i, s := range songsInQueue {
			s.QueuePosition = i + 1
			// db.Model(&model.RequestedSong{}).Where("id = ?", s.ID).Update("queuePosition", i+1)
			db.Save(&s)
		}

		entity.QueuePosition = len(songsInQueue) + 1

		err = db.Create(&entity).Error

		if err != nil {
			log.Fatal(err)
			result.Result = append(result.Result, "internal error")
			return result
		}

		timeForWait := 0 * time.Minute
		for _, s := range songsInQueue {
			timeForWait = time.Duration(s.Duration)*time.Millisecond + timeForWait
		}

		message := fasttemplate.ExecuteString(
			parsedSettings.Translations.Song.RequestedMessage,
			"{{", "}}",
			map[string]interface{}{
				"songId":    entity.VideoID,
				"songTitle": ytdlSongInfo.Title,
				"position":  strconv.Itoa(len(songsInQueue) + 1),
				"waitTime":  timeForWait.String(),
			},
		)
		result.Result = append(result.Result, message)

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
	settings *youtube.YouTubeSettings,
	song *ytdl.Video,
) error {
	db := do.MustInvoke[gorm.DB](di.Provider)
	twitchClient := do.MustInvoke[twitch.Twitch](di.Provider)

	if userId != channelId {
		return nil
	}

	if len(settings.DenyList.Users) > 0 {
		_, isUserDenied := lo.Find(
			settings.DenyList.Users,
			func(u youtube.YouTubeDenySettingsUsers) bool {
				return u.UserID == userId
			},
		)

		if isUserDenied {
			return errors.New(settings.Translations.User.Denied)
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
			return errors.New(settings.Translations.Channel.Denied)
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
			return errors.New(settings.Translations.Song.Denied)
		}
	}

	if *settings.AcceptOnlyWhenOnline {
		stream := &model.ChannelsStreams{}
		db.Where(`"userId" = ?`, channelId).First(stream)
		if stream.ID == "" {
			return errors.New(settings.Translations.AcceptOnlineWhenOnline)
		}
	}

	if settings.MaxRequests != 0 {
		var count int64
		db.Model(&model.RequestedSong{}).
			Where(`"channelId" = ? AND "deletedAt" IS NULL`, channelId).
			Count(&count)
		if count >= int64(settings.MaxRequests) {
			message := fasttemplate.ExecuteString(
				settings.Translations.Song.MaximumOrdered,
				"{{", "}}",
				map[string]interface{}{
					"maximum": strconv.Itoa(settings.MaxRequests),
				},
			)
			return errors.New(message)
		}
	}

	if settings.Song.MinViews != 0 && song.Views < settings.Song.MinViews {
		message := fasttemplate.ExecuteString(
			settings.Translations.Song.MinLength,
			"{{", "}}",
			map[string]interface{}{
				"songTitle":   song.Title,
				"songId":      song.ID,
				"songViews":   strconv.Itoa(song.Views),
				"neededViews": strconv.Itoa(settings.Song.MinViews),
			},
		)
		return errors.New(message)
	}

	songDuration := int(song.Duration.Minutes())
	if settings.Song.MaxLength != 0 && songDuration > settings.Song.MaxLength {
		message := fasttemplate.ExecuteString(
			settings.Translations.Song.MaxLength,
			"{{", "}}",
			map[string]interface{}{
				"songTitle": song.Title,
				"songId":    song.ID,
				"songViews": strconv.Itoa(song.Views),
				"maxLength": strconv.Itoa(settings.Song.MaxLength),
			},
		)
		return errors.New(message)
	}

	if settings.Song.MinLength != 0 && songDuration < settings.Song.MinLength {
		message := fasttemplate.ExecuteString(
			settings.Translations.Song.MinLength,
			"{{", "}}",
			map[string]interface{}{
				"songTitle": song.Title,
				"songId":    song.ID,
				"songViews": strconv.Itoa(song.Views),
				"minLength": strconv.Itoa(settings.Song.MinLength),
			},
		)
		return errors.New(message)
	}

	// TODO: check categories

	if settings.User.MaxRequests != 0 {
		var count int64
		db.
			Model(&model.RequestedSong{}).
			Where(`"orderedById" = ? AND "channelId" = ? AND "deletedAt" IS NULL`, userId, channelId).
			Count(&count)
		if count >= int64(settings.User.MaxRequests) {
			message := fasttemplate.ExecuteString(
				settings.Translations.Song.MaximumOrdered,
				"{{", "}}",
				map[string]interface{}{
					"count": strconv.Itoa(settings.User.MaxRequests),
				},
			)
			return errors.New(message)
		}
	}

	if settings.User.MinMessages != 0 || settings.User.MinWatchTime != 0 {
		user := &model.Users{}
		db.Where("id = ?", userId).Preload("Stats").First(&user)
		if user.ID == "" {
			return errors.New(
				"there is restrictions on user, but i cannot find you in db, sorry. :(",
			)
		}

		if settings.User.MinMessages != 0 &&
			user.Stats.Messages < int32(settings.User.MinMessages) {
			message := fasttemplate.ExecuteString(
				settings.Translations.User.MinMessages,
				"{{", "}}",
				map[string]interface{}{
					"minMessages":  strconv.Itoa(settings.User.MinMessages),
					"userMessages": strconv.Itoa(int(user.Stats.Messages)),
				},
			)
			return errors.New(message)
		}

		watchedInMinutes := time.Duration(user.Stats.Watched) * time.Millisecond
		if settings.User.MinWatchTime != 0 &&
			int64(watchedInMinutes.Minutes()) < settings.User.MinWatchTime {
			message := fasttemplate.ExecuteString(
				settings.Translations.User.MinWatched,
				"{{", "}}",
				map[string]interface{}{
					"minWatched":  strconv.Itoa(int(settings.User.MinWatchTime)),
					"userWatched": strconv.Itoa(int(watchedInMinutes)),
				},
			)
			return errors.New(message)
		}
	}

	if settings.User.MinFollowTime != 0 {
		neededDuration := time.Minute * time.Duration(settings.User.MinFollowTime)
		followReq, err := twitchClient.Client.GetUsersFollows(&helix.UsersFollowsParams{
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
			message := fasttemplate.ExecuteString(
				settings.Translations.User.MinFollow,
				"{{", "}}",
				map[string]interface{}{
					"minFollow":  strconv.Itoa(settings.User.MinFollowTime),
					"userFollow": strconv.Itoa(int(followDuration.Minutes())),
				},
			)
			return errors.New(message)
		}
	}

	return nil
}
