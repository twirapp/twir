package sr_youtube

import (
	"context"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	command_arguments "github.com/twirapp/twir/apps/parser/internal/command-arguments"
	"github.com/twirapp/twir/apps/parser/internal/types"
	"github.com/twirapp/twir/apps/parser/internal/types/services"
	"github.com/twirapp/twir/libs/bus-core/ytsr"

	"github.com/guregu/null"
	"github.com/twirapp/twir/libs/twitch"
	"github.com/twirapp/twir/libs/grpc/websockets"
	"github.com/valyala/fasttemplate"

	model "github.com/twirapp/twir/libs/gomodels"

	"github.com/nicklaw5/helix/v2"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"

	"github.com/samber/lo"
)

type ReqError struct {
	Title string
	Error string
}

const (
	songRequestArgName = "name or link"
)

var SrCommand = &types.DefaultCommand{
	ChannelsCommands: &model.ChannelsCommands{
		Name:        "sr",
		Description: null.StringFrom("Song requests from youtube"),
		Module:      "SONGS",
		IsReply:     true,
		Visible:     true,
	},
	Args: []command_arguments.Arg{
		command_arguments.VariadicString{
			Name: songRequestArgName,
		},
	},
	Handler: func(ctx context.Context, parseCtx *types.ParseContext) (
		*types.CommandsHandlerResult,
		error,
	) {
		result := &types.CommandsHandlerResult{}

		moduleSettings := &model.ChannelSongRequestsSettings{}
		err := parseCtx.Services.Gorm.WithContext(ctx).
			Where(`"channel_id" = ?`, parseCtx.Channel.ID).
			First(moduleSettings).Error

		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return result, nil
			} else {
				return nil, &types.CommandHandlerError{
					Message: "cannot get song requests settings",
					Err:     err,
				}
			}
		}

		if !moduleSettings.Enabled {
			return result, nil
		}

		if moduleSettings.AcceptOnlyWhenOnline {
			stream := &model.ChannelsStreams{}
			parseCtx.Services.Gorm.WithContext(ctx).Where(
				`"userId" = ?`,
				parseCtx.Channel.ID,
			).First(stream)
			if stream.ID == "" {
				result.Result = append(result.Result, moduleSettings.TranslationsAcceptOnlineWhenOnline)
				return result, nil
			}
		}

		req, err := parseCtx.Services.Bus.YTSRSearch.Request(
			ctx,
			ytsr.SearchRequest{
				Search: parseCtx.ArgsParser.Get(songRequestArgName).String(),
			},
		)
		if err != nil {
			return nil, &types.CommandHandlerError{
				Message: "cannot search song",
				Err:     err,
			}
		}
		if len(req.Data.Songs) == 0 {
			result.Result = append(result.Result, moduleSettings.TranslationsSongNotFound)
			return result, nil
		}

		latestSong := &model.RequestedSong{}

		err = parseCtx.Services.Gorm.WithContext(ctx).
			Where(`"channelId" = ? AND "deletedAt" IS NULL`, parseCtx.Channel.ID).
			Order(`"createdAt" desc`).
			Find(&latestSong).Error
		if err != nil {
			return nil, &types.CommandHandlerError{
				Message: "cannot get latest song",
				Err:     err,
			}
		}

		requested := make([]*model.RequestedSong, 0, len(req.Data.Songs))
		errors := make([]*ReqError, 0, len(req.Data.Songs))

		var currentQueueCount int64
		err = parseCtx.Services.Gorm.WithContext(ctx).
			Where(`"channelId" = ? AND "deletedAt" IS NULL`, parseCtx.Channel.ID).
			Model(&model.RequestedSong{}).
			Count(&currentQueueCount).
			Error

		if err != nil {
			return nil, &types.CommandHandlerError{
				Message: "cannot get current queue count",
				Err:     err,
			}
		}

		for i, song := range req.Data.Songs {
			err = validate(
				ctx,
				parseCtx.Services,
				parseCtx.Channel.ID,
				parseCtx.Sender.ID,
				*moduleSettings,
				song,
			)

			if err != nil {
				errors = append(
					errors, &ReqError{
						Title: song.Title,
						Error: err.Error(),
					},
				)
			} else {
				model := &model.RequestedSong{
					ID:                   uuid.NewV4().String(),
					ChannelID:            parseCtx.Channel.ID,
					OrderedById:          parseCtx.Sender.ID,
					OrderedByName:        parseCtx.Sender.Name,
					OrderedByDisplayName: null.StringFrom(parseCtx.Sender.DisplayName),
					VideoID:              song.Id,
					Title:                song.Title,
					Duration:             int32(song.Duration),
					CreatedAt:            time.Now().UTC(),
					QueuePosition:        int(currentQueueCount) + (i + 1),
					SongLink:             null.StringFromPtr(song.Link),
				}

				err = parseCtx.Services.Gorm.WithContext(ctx).Create(model).Error
				if err == nil {
					requested = append(requested, model)
				}
			}
		}

		if len(requested) > 0 {
			requestedMapped := lo.Map(
				requested, func(item *model.RequestedSong, _ int) string {
					return fmt.Sprintf("%s (#%v)", item.Title, item.QueuePosition)
				},
			)

			result.Result = append(result.Result, "✅ "+strings.Join(requestedMapped, " · "))
		}

		if len(errors) > 0 {
			errorsMapped := lo.Map(
				errors, func(item *ReqError, _ int) string {
					return item.Title + " - " + item.Error
				},
			)
			result.Result = append(result.Result, "❌"+strings.Join(errorsMapped, " · "))
		}

		for _, song := range requested {
			parseCtx.Services.GrpcClients.WebSockets.YoutubeAddSongToQueue(
				context.Background(),
				&websockets.YoutubeAddSongToQueueRequest{
					ChannelId: parseCtx.Channel.ID,
					EntityId:  song.ID,
				},
			)
		}

		return result, nil
	},
}

func validate(
	ctx context.Context,
	services *services.Services,
	channelId, userId string,
	settings model.ChannelSongRequestsSettings,
	song ytsr.Song,
) error {
	alreadyRequestedSong := &model.RequestedSong{}
	services.Gorm.WithContext(ctx).Where(
		`"videoId" = ? AND "deletedAt" IS NULL AND "channelId" = ?`,
		song.Id,
		channelId,
	).
		Find(&alreadyRequestedSong)

	if alreadyRequestedSong.ID != "" {
		return errors.New(settings.TranslationsSongAlreadyInQueue)
	}

	// if channelId == userId {
	//	return nil
	// }

	twitchClient, err := twitch.NewAppClientWithContext(
		ctx,
		*services.Config,
		services.Bus,
	)
	if err != nil {
		return err
	}

	if song.IsLive {
		return errors.New(settings.TranslationsSongLive)
	}

	if len(settings.DenyListWords) > 0 {
		for _, word := range settings.DenyListWords {
			if word == "" {
				continue
			}
			if strings.Contains(strings.ToLower(song.Title), strings.ToLower(word)) {
				return errors.New(settings.TranslationsSongDenied)
			}
		}

		for _, word := range settings.DenyListWords {
			if word == "" {
				continue
			}

			if strings.Contains(strings.ToLower(song.Author.Name), strings.ToLower(word)) {
				return errors.New(settings.TranslationsSongDenied)
			}
		}
	}

	if len(settings.DenyListUsers) > 0 {
		_, isUserDenied := lo.Find(
			settings.DenyListUsers,
			func(u string) bool {
				return u == userId
			},
		)

		if isUserDenied {
			return errors.New(settings.TranslationsUserDenied)
		}
	}

	if len(settings.DenyListChannels) > 0 {
		_, isChannelBlacklisted := lo.Find(
			settings.DenyListChannels,
			func(u string) bool {
				return u == song.Author.ChannelId
			},
		)

		if isChannelBlacklisted {
			return errors.New(settings.TranslationsChannelDenied)
		}
	}

	if len(settings.DenyListSongs) > 0 {
		_, isSongBlackListed := lo.Find(
			settings.DenyListSongs,
			func(u string) bool {
				return u == song.Id
			},
		)

		if isSongBlackListed {
			return errors.New(settings.TranslationsSongDenied)
		}
	}

	if settings.MaxRequests != 0 {
		var count int64
		services.Gorm.WithContext(ctx).Model(&model.RequestedSong{}).
			Where(`"channelId" = ? AND "deletedAt" IS NULL`, channelId).
			Count(&count)
		if count >= int64(settings.MaxRequests) {
			message := fasttemplate.ExecuteString(
				settings.TranslationsSongMaximumOrdered,
				"{{", "}}",
				map[string]interface{}{
					"maximum": strconv.Itoa(settings.MaxRequests),
				},
			)
			return errors.New(message)
		}
	}

	if settings.SongMinViews != 0 && int(song.Views) < settings.SongMinViews {
		message := fasttemplate.ExecuteString(
			settings.TranslationsSongMinViews,
			"{{", "}}",
			map[string]interface{}{
				"songTitle":   song.Title,
				"songId":      song.Id,
				"songViews":   strconv.Itoa(int(song.Views)),
				"neededViews": strconv.Itoa(settings.SongMinViews),
			},
		)
		return errors.New(message)
	}

	songDuration := time.Duration(song.Duration) * time.Second
	if settings.SongMaxLength != 0 && int(math.Round(songDuration.Minutes())) > settings.SongMaxLength {
		message := fasttemplate.ExecuteString(
			settings.TranslationsSongMaxLength,
			"{{", "}}",
			map[string]interface{}{
				"songTitle": song.Title,
				"songId":    song.Id,
				"songViews": strconv.Itoa(int(song.Views)),
				"maxLength": strconv.Itoa(settings.SongMaxLength),
			},
		)
		return errors.New(message)
	}

	if settings.SongMinLength != 0 && int(math.Round(songDuration.Minutes())) < settings.SongMinLength {
		message := fasttemplate.ExecuteString(
			settings.TranslationsSongMinLength,
			"{{", "}}",
			map[string]interface{}{
				"songTitle": song.Title,
				"songId":    song.Id,
				"songViews": strconv.Itoa(int(song.Views)),
				"minLength": strconv.Itoa(settings.SongMinLength),
			},
		)
		return errors.New(message)
	}

	// TODO: check categories

	if settings.UserMaxRequests != 0 {
		var count int64
		services.Gorm.WithContext(ctx).
			Model(&model.RequestedSong{}).
			Where(`"orderedById" = ? AND "channelId" = ? AND "deletedAt" IS NULL`, userId, channelId).
			Count(&count)
		if count >= int64(settings.UserMaxRequests) {
			message := fasttemplate.ExecuteString(
				settings.TranslationsUserMaxRequests,
				"{{", "}}",
				map[string]interface{}{
					"count":   strconv.Itoa(settings.UserMaxRequests),
					"maximum": strconv.Itoa(settings.UserMaxRequests),
				},
			)
			return errors.New(message)
		}
	}

	if settings.UserMinMessages != 0 || settings.UserMinWatchTime != 0 {
		user := &model.Users{}
		services.Gorm.WithContext(ctx).Where("id = ?", userId).Preload("Stats").First(&user)
		if user.ID == "" {
			return errors.New(
				"there is restrictions on user, but i cannot find you in db, sorry. :(",
			)
		}

		if settings.UserMinMessages != 0 &&
			user.Stats.Messages < int32(settings.UserMinMessages) {
			message := fasttemplate.ExecuteString(
				settings.TranslationsUserMinMessages,
				"{{", "}}",
				map[string]interface{}{
					"neededMessages": strconv.Itoa(settings.UserMinMessages),
					"userMessages":   strconv.Itoa(int(user.Stats.Messages)),
				},
			)
			return errors.New(message)
		}

		watchedInMinutes := time.Duration(user.Stats.Watched) * time.Millisecond
		if settings.UserMinWatchTime != 0 &&
			int64(watchedInMinutes.Minutes()) < settings.UserMinWatchTime {
			message := fasttemplate.ExecuteString(
				settings.TranslationsUserMinWatched,
				"{{", "}}",
				map[string]interface{}{
					"minWatched":  strconv.Itoa(int(settings.UserMinWatchTime)),
					"userWatched": strconv.Itoa(int(watchedInMinutes)),
				},
			)
			return errors.New(message)
		}
	}

	if settings.UserMinFollowTime != 0 {
		neededDuration := time.Minute * time.Duration(settings.UserMinFollowTime)
		followReq, err := twitchClient.GetUsersFollows(
			&helix.UsersFollowsParams{
				FromID: userId,
				ToID:   channelId,
			},
		)
		if err != nil {
			return errors.New("Internal error when checking follow")
		}
		if followReq.Data.Total == 0 {
			return errors.New("For request song you need to be a followed")
		}

		followDuration := time.Since(followReq.Data.Follows[0].FollowedAt)
		if followDuration.Minutes() < neededDuration.Minutes() {
			message := fasttemplate.ExecuteString(
				settings.TranslationsUserMinFollow,
				"{{", "}}",
				map[string]interface{}{
					"minFollow":  strconv.Itoa(settings.UserMinFollowTime),
					"userFollow": strconv.Itoa(int(followDuration.Minutes())),
				},
			)
			return errors.New(message)
		}
	}

	return nil
}
