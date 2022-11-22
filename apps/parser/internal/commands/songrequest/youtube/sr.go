package sr_youtube

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"regexp"
	"time"
	model "tsuwari/models"
	"tsuwari/parser/internal/config/twitch"
	"tsuwari/parser/internal/types"
	variables_cache "tsuwari/parser/internal/variablescache"

	"github.com/satont/go-helix/v2"

	youtubenats "github.com/satont/tsuwari/libs/nats/youtube"
	"google.golang.org/protobuf/proto"

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
			result.Result = append(result.Result, "cannot get information about song.")
			return result
		}

		if ytdlSongInfo.Duration.Seconds() == 0 {
			result.Result = append(result.Result, "seems like that song is live, which is disallowed.")
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
			parsedSettings := &youtube.YoutubeSettings{}
			err = json.Unmarshal(moduleSettings.Settings, parsedSettings)
			if err != nil {
				fmt.Println(err)
				result.Result = append(result.Result, "internal error")
				return result
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

		err = ctx.Services.Db.Create(&entity).Error

		if err != nil {
			log.Fatal(err)
			result.Result = append(result.Result, "internal error")
			return result
		}

		songsInQueue := []model.RequestedSong{}
		ctx.Services.Db.
			Where(
				`"channelId" = ? AND "id" != ? AND "deletedAt" IS NULL`,
				ctx.ChannelId,
				entity.ID,
			).
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

func validate(
	channelId, userId string,
	db *gorm.DB,
	tw *twitch.Twitch,
	settings *youtube.YoutubeSettings,
	song *ytdl.Video,
) error {
	if userId != channelId {
		return nil
	}

	if settings.BlackList != nil {
		if settings.BlackList.Users != nil {
			_, isUserBlackListed := lo.Find(
				settings.BlackList.Users,
				func(u youtube.YoutubeBlacklistSettingsUsers) bool {
					return u.UserID == userId
				},
			)

			if isUserBlackListed {
				return errors.New("you cannot request song because you are blacklisted")
			}
		}

		if settings.BlackList.Channels != nil {
			_, isChannelBlacklisted := lo.Find(
				settings.BlackList.Channels,
				func(u youtube.YoutubeBlacklistSettingsChannels) bool {
					return u.ID == song.ChannelID
				},
			)

			if isChannelBlacklisted {
				return errors.New("you cannot request that song because channel is blacklisted")
			}
		}

	}

	if settings.AcceptOnlyWhenOnline != nil && *settings.AcceptOnlyWhenOnline {
		stream := &model.ChannelsStreams{}
		db.Where(`"userId" = ?`, channelId).First(stream)
		if stream.ID == "" {
			return errors.New("requests accepted only on online streams")
		}
	}

	if settings.MaxRequests != nil {
		var count int64
		db.Model(&model.RequestedSong{}).
			Where(`"channelId" = ? AND "deletedAt" IS NULL`, channelId).
			Count(&count)
		if count >= int64(*settings.MaxRequests) {
			return errors.New("maximum number of tracks ordered now, try later")
		}
	}

	if settings.Song != nil {
		if settings.Song.MinViews != nil && song.Views < *settings.Song.MinViews {
			return errors.New(
				fmt.Sprintf("song haven't %v views for request", *settings.Song.MinViews),
			)
		}

		if settings.Song.MaxLength != nil &&
			song.Duration > time.Minute*time.Duration(*settings.Song.MaxLength) {
			return errors.New("that song is to long for request")
		}

		// TODO: check categories
	}

	if settings.User != nil {
		if settings.User.MaxRequests != nil {
			var count int64
			db.
				Model(&model.RequestedSong{}).
				Where(`"orderedById" = ? AND "channelId" = ? AND "deletedAt" IS NULL`, userId, channelId).
				Count(&count)
			if count >= int64(*settings.User.MaxRequests) {
				return errors.New("maximum number of tracks ordered now, try later")
			}
		}

		if settings.User.MinMessages != nil || settings.User.MinWatchTime != nil {
			user := &model.Users{}
			db.Where("id = ?", userId).Preload("Stats").First(&user)
			if user.ID == "" {
				return errors.New(
					"there is restrictions on user, but i cannot find you in db, sorry. :(",
				)
			}

			if settings.User.MinMessages != nil &&
				user.Stats.Messages < *settings.User.MinMessages {
				return errors.New(
					fmt.Sprintf("you haven't %v messages for request song", user.Stats.Messages),
				)
			}

			watchedInMinutes := time.Duration(user.Stats.Watched) * time.Millisecond
			if settings.User.MinWatchTime != nil &&
				int64(watchedInMinutes.Minutes()) < *settings.User.MinWatchTime {

				return errors.New(
					fmt.Sprintf(
						"you haven't watched stream for %v minutes for request song",
						time.Minute*time.Duration(*settings.User.MinWatchTime),
					),
				)
			}
		}

		if settings.User.MinFollowTime != nil {
			neededDuration := time.Minute * time.Duration(*settings.User.MinFollowTime)
			followReq, err := tw.Client.GetUsersFollows(&helix.UsersFollowsParams{
				FromID: userId,
				ToID:   channelId,
			})
			if err != nil {
				return errors.New("internal error when checking follow")
			}
			if followReq.Data.Total == 0 {
				return errors.New("for request song you need to be a followed")
			}

			followDuration := time.Since(followReq.Data.Follows[0].FollowedAt)
			if followDuration.Minutes() < neededDuration.Minutes() {
				return errors.New(
					fmt.Sprintf(
						"you need to be follower at least %v minutes for request song",
						neededDuration.Minutes(),
					),
				)
			}
		}
	}

	return nil
}
