package sr_youtube

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/guregu/null"
	config "github.com/satont/tsuwari/libs/config"
	"github.com/satont/tsuwari/libs/grpc/generated/tokens"
	"github.com/satont/tsuwari/libs/grpc/generated/ytsr"
	"github.com/satont/tsuwari/libs/twitch"
	"github.com/valyala/fasttemplate"
	"go.uber.org/zap"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/samber/do"
	"github.com/satont/tsuwari/apps/parser/internal/di"
	"github.com/satont/tsuwari/libs/grpc/generated/websockets"

	"github.com/satont/tsuwari/apps/parser/internal/types"
	variables_cache "github.com/satont/tsuwari/apps/parser/internal/variablescache"
	model "github.com/satont/tsuwari/libs/gomodels"

	"github.com/satont/go-helix/v2"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"

	"github.com/samber/lo"
	youtube "github.com/satont/tsuwari/libs/types/types/api/modules"
)

type ReqError struct {
	Title string
	Error string
}

var SrCommand = types.DefaultCommand{
	Command: types.Command{
		Name:               "sr",
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
		search := do.MustInvoke[ytsr.YtsrClient](di.Provider)

		result := &types.CommandsHandlerResult{}

		if ctx.Text == nil {
			result.Result = append(result.Result, "You should provide text for song request")
			return result
		}

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

		req, err := search.Search(context.Background(), &ytsr.SearchRequest{Search: *ctx.Text})
		if err != nil {
			fmt.Println(err)
			return result
		}
		if len(req.Songs) == 0 {
			result.Result = append(result.Result, parsedSettings.Translations.Song.NotFound)
			return result
		}

		var songsCount int64

		err = db.
			Model(&model.RequestedSong{}).
			Where(`"channelId" = ? AND "deletedAt" IS NULL`, ctx.ChannelId).
			Order(`"createdAt" asc`).
			Count(&songsCount).Error

		requested := make([]*model.RequestedSong, 0, len(req.Songs))
		errors := make([]*ReqError, 0, len(req.Songs))

		for i, song := range req.Songs {
			err = validate(
				ctx.ChannelId,
				ctx.SenderId,
				parsedSettings,
				song,
			)

			if err != nil {
				errors = append(errors, &ReqError{
					Title: song.Title,
					Error: err.Error(),
				})
			} else {
				model := &model.RequestedSong{
					ID:                   uuid.NewV4().String(),
					ChannelID:            ctx.ChannelId,
					OrderedById:          ctx.SenderId,
					OrderedByName:        ctx.SenderName,
					OrderedByDisplayName: null.StringFrom(ctx.SenderDisplayName),
					VideoID:              song.Id,
					Title:                song.Title,
					Duration:             int32(song.Duration),
					CreatedAt:            time.Now().UTC(),
					QueuePosition:        int(songsCount) + i + 1,
				}

				err = db.Create(model).Error
				if err == nil {
					requested = append(requested, model)
				}
			}
		}

		if len(requested) > 0 {
			requestedMapped := lo.Map(requested, func(item *model.RequestedSong, _ int) string {
				return fmt.Sprintf("%s (#%v)", item.Title, item.QueuePosition)
			})

			result.Result = append(result.Result, "✅ "+strings.Join(requestedMapped, " · "))
		}

		if len(errors) > 0 {
			errorsMapped := lo.Map(errors, func(item *ReqError, _ int) string {
				return item.Title + " - " + item.Error
			})
			result.Result = append(result.Result, "❌"+strings.Join(errorsMapped, " · "))
		}

		for _, song := range requested {
			websocketGrpc.YoutubeAddSongToQueue(
				context.Background(),
				&websockets.YoutubeAddSongToQueueRequest{
					ChannelId: ctx.ChannelId,
					EntityId:  song.ID,
				},
			)
		}

		return result
	},
}

func validate(
	channelId, userId string,
	settings *youtube.YouTubeSettings,
	song *ytsr.Song,
) error {
	db := do.MustInvoke[gorm.DB](di.Provider)
	cfg := do.MustInvoke[config.Config](di.Provider)
	tokensGrpc := do.MustInvoke[tokens.TokensClient](di.Provider)

	twitchClient, err := twitch.NewAppClient(cfg, tokensGrpc)

	if err != nil {
		return err
	}

	alreadyRequestedSong := &model.RequestedSong{}
	db.Where(`"videoId" = ? AND "deletedAt" IS NULL AND "channelId" = ?`, song.Id, channelId).
		Find(&alreadyRequestedSong)

	if alreadyRequestedSong.ID != "" {
		return errors.New(settings.Translations.Song.AlreadyInQueue)
	}

	if song.IsLive {
		return errors.New(settings.Translations.Song.Live)
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

	if len(settings.DenyList.Channels) > 0 && song.Author != nil {
		_, isChannelBlacklisted := lo.Find(
			settings.DenyList.Channels,
			func(u youtube.YouTubeDenySettingsChannels) bool {
				return u.ID == song.Author.ChannelId
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
				return u.ID == song.Id
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

	if settings.Song.MinViews != 0 && int(song.Views) < settings.Song.MinViews {
		message := fasttemplate.ExecuteString(
			settings.Translations.Song.MinViews,
			"{{", "}}",
			map[string]interface{}{
				"songTitle":   song.Title,
				"songId":      song.Id,
				"songViews":   strconv.Itoa(int(song.Views)),
				"neededViews": strconv.Itoa(settings.Song.MinViews),
			},
		)
		return errors.New(message)
	}

	songDuration := time.Duration(song.Duration) * time.Millisecond
	if settings.Song.MaxLength != 0 && int(math.Round(songDuration.Minutes())) > settings.Song.MaxLength {
		message := fasttemplate.ExecuteString(
			settings.Translations.Song.MaxLength,
			"{{", "}}",
			map[string]interface{}{
				"songTitle": song.Title,
				"songId":    song.Id,
				"songViews": strconv.Itoa(int(song.Views)),
				"maxLength": strconv.Itoa(settings.Song.MaxLength),
			},
		)
		return errors.New(message)
	}

	if settings.Song.MinLength != 0 && int(math.Round(songDuration.Minutes())) < settings.Song.MinLength {
		message := fasttemplate.ExecuteString(
			settings.Translations.Song.MinLength,
			"{{", "}}",
			map[string]interface{}{
				"songTitle": song.Title,
				"songId":    song.Id,
				"songViews": strconv.Itoa(int(song.Views)),
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
				settings.Translations.User.MaxRequests,
				"{{", "}}",
				map[string]interface{}{
					"count":   strconv.Itoa(settings.User.MaxRequests),
					"maximum": strconv.Itoa(settings.User.MaxRequests),
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
		followReq, err := twitchClient.GetUsersFollows(&helix.UsersFollowsParams{
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
