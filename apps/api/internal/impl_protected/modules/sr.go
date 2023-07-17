package modules

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/raitonoberu/ytsearch"
	"github.com/samber/lo"
	loParallel "github.com/samber/lo/parallel"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/grpc/generated/api/modules_sr"
	"github.com/satont/twir/libs/types/types/api/modules"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
	"strings"
	"sync"
)

const SrType = "youtube_song_requests"

func getThumbNailUrl(url string) string {
	return strings.Replace(url, "http", "https", 1)
}

func (c *Modules) ModulesSRGet(
	ctx context.Context,
	_ *emptypb.Empty,
) (*modules_sr.GetResponse, error) {
	dashboardId := ctx.Value("dashboardId").(string)
	entity := &model.ChannelModulesSettings{}
	if err := c.Db.
		WithContext(ctx).
		Where(`"channelId" = ? AND "type" = ?`, dashboardId, SrType).
		First(entity).Error; err != nil {
		return nil, err
	}

	settings := &modules.YouTubeSettings{}
	if err := json.Unmarshal(entity.Settings, settings); err != nil {
		return nil, err
	}

	return &modules_sr.GetResponse{
		Data: &modules_sr.YouTubeSettings{
			Enabled:               *settings.Enabled,
			AcceptOnlyWhenOnline:  *settings.AcceptOnlyWhenOnline,
			MaxRequests:           int32(settings.MaxRequests),
			ChannelPointsRewardId: settings.ChannelPointsRewardId,
			AnnouncePlay:          *settings.AnnouncePlay,
			NeededVotesVorSkip:    float32(settings.NeededVotesVorSkip),
			User: &modules_sr.YouTubeUserSettings{
				MaxRequests:   int32(settings.User.MaxRequests),
				MinWatchTime:  int32(settings.User.MinWatchTime),
				MinMessages:   int32(settings.User.MinMessages),
				MinFollowTime: int32(settings.User.MinFollowTime),
			},
			Song: &modules_sr.YouTubeSongSettings{
				MinLength:          int32(settings.Song.MinLength),
				MaxLength:          int32(settings.Song.MaxLength),
				MinViews:           int32(settings.Song.MinViews),
				AcceptedCategories: settings.Song.AcceptedCategories,
			},
			DenyList: &modules_sr.YouTubeDenyList{
				Users:        settings.DenyList.Users,
				Songs:        settings.DenyList.Songs,
				Channels:     settings.DenyList.Channels,
				ArtistsNames: settings.DenyList.ArtistsNames,
			},
			Translations: &modules_sr.YouTubeTranslations{
				NowPlaying:           settings.Translations.NowPlaying,
				NotEnabled:           settings.Translations.NotEnabled,
				NoText:               settings.Translations.NoText,
				AcceptOnlyWhenOnline: settings.Translations.AcceptOnlineWhenOnline,
				User: &modules_sr.YouTubeUserTranslations{
					Denied:      settings.Translations.User.Denied,
					MaxRequests: settings.Translations.User.MaxRequests,
					MinMessages: settings.Translations.User.MinMessages,
					MinWatched:  settings.Translations.User.MinWatched,
					MinFollow:   settings.Translations.User.MinFollow,
				},
				Song: &modules_sr.YouTubeSongTranslations{
					Denied:               settings.Translations.Song.Denied,
					NotFound:             settings.Translations.Song.NotFound,
					AlreadyInQueue:       settings.Translations.Song.AlreadyInQueue,
					AgeRestrictions:      settings.Translations.Song.AgeRestrictions,
					CannotGetInformation: settings.Translations.Song.CannotGetInformation,
					Live:                 settings.Translations.Song.Live,
					MaxLength:            settings.Translations.Song.MaxLength,
					MinLength:            settings.Translations.Song.MinLength,
					RequestedMessage:     settings.Translations.Song.RequestedMessage,
					MaximumOrdered:       settings.Translations.Song.MaximumOrdered,
					MinViews:             settings.Translations.Song.MinViews,
				},
				Channel: &modules_sr.YouTubeChannelTranslations{
					Denied: settings.Translations.Channel.Denied,
				},
			},
		},
	}, nil
}

func (c *Modules) ModulesSRSearchVideosOrChannels(
	_ context.Context,
	request *modules_sr.GetSearchRequest,
) (*modules_sr.GetSearchResponse, error) {
	response := &modules_sr.GetSearchResponse{
		Items: make([]*modules_sr.GetSearchResponse_Result, 0, len(request.Query)),
	}

	if len(request.Query) == 0 {
		return response, nil
	}

	var mu sync.Mutex

	loParallel.ForEach(request.Query, func(query string, _ int) {
		if query == "" {
			return
		}

		var search *ytsearch.SearchClient
		if request.Type == modules_sr.GetSearchRequest_CHANNEL {
			search = ytsearch.ChannelSearch(query)
		} else {
			search = ytsearch.VideoSearch(query)
		}

		res, err := search.Next()
		if err != nil {
			zap.S().Error(err)
			return
		}

		mu.Lock()
		defer mu.Unlock()
		if request.Type == modules_sr.GetSearchRequest_CHANNEL {
			channels := lo.Map(res.Channels, func(item *ytsearch.ChannelItem, _ int) *modules_sr.GetSearchResponse_Result {
				thumb := getThumbNailUrl(item.Thumbnails[len(item.Thumbnails)-1].URL)
				return &modules_sr.GetSearchResponse_Result{
					Id:        item.ID,
					Title:     item.Title,
					Thumbnail: thumb,
				}
			})
			response.Items = append(
				response.Items,
				channels...,
			)
		} else {
			videos := lo.Map(res.Videos, func(item *ytsearch.VideoItem, _ int) *modules_sr.GetSearchResponse_Result {
				thumb := getThumbNailUrl(item.Thumbnails[len(item.Thumbnails)-1].URL)

				return &modules_sr.GetSearchResponse_Result{
					Id:        item.ID,
					Title:     item.Title,
					Thumbnail: thumb,
				}
			})
			response.Items = append(
				response.Items,
				videos...,
			)
		}
	})

	return response, nil
}

func (c *Modules) ModulesSRUpdate(
	ctx context.Context,
	request *modules_sr.PostRequest,
) (*emptypb.Empty, error) {
	dashboardId := ctx.Value("dashboardId").(string)
	entity := &model.ChannelModulesSettings{}
	if err := c.Db.
		WithContext(ctx).
		Where(`"channelId" = ? AND "type" = ?`, dashboardId, SrType).
		Find(entity).Error; err != nil {
		return nil, err
	}

	if entity.ID == "" {
		entity.ID = uuid.New().String()
	}

	entity.ChannelId = dashboardId
	entity.Type = SrType

	settings := &modules.YouTubeSettings{
		Enabled:               &request.Data.Enabled,
		AcceptOnlyWhenOnline:  &request.Data.AcceptOnlyWhenOnline,
		MaxRequests:           int(request.Data.MaxRequests),
		ChannelPointsRewardId: request.Data.ChannelPointsRewardId,
		AnnouncePlay:          &request.Data.AnnouncePlay,
		NeededVotesVorSkip:    float64(request.Data.NeededVotesVorSkip),
		User: modules.YouTubeUserSettings{
			MaxRequests:   int(request.Data.User.MaxRequests),
			MinWatchTime:  int64(request.Data.User.MinWatchTime),
			MinMessages:   int(request.Data.User.MinMessages),
			MinFollowTime: int(request.Data.User.MinFollowTime),
		},
		Song: modules.YouTubeSongSettings{
			MinLength:          int(request.Data.Song.MinLength),
			MaxLength:          int(request.Data.Song.MaxLength),
			MinViews:           int(request.Data.Song.MinViews),
			AcceptedCategories: request.Data.Song.AcceptedCategories,
		},
		DenyList: modules.YouTubeDenyList{
			Users:        request.Data.DenyList.Users,
			Songs:        request.Data.DenyList.Songs,
			Channels:     request.Data.DenyList.Channels,
			ArtistsNames: request.Data.DenyList.ArtistsNames,
		},
		Translations: modules.YouTubeTranslations{
			NowPlaying:             request.Data.Translations.NowPlaying,
			NotEnabled:             request.Data.Translations.NotEnabled,
			NoText:                 request.Data.Translations.NoText,
			AcceptOnlineWhenOnline: request.Data.Translations.AcceptOnlyWhenOnline,
			User: modules.YouTubeUserTranslations{
				Denied:      request.Data.Translations.User.Denied,
				MaxRequests: request.Data.Translations.User.MaxRequests,
				MinMessages: request.Data.Translations.User.MinMessages,
				MinWatched:  request.Data.Translations.User.MinWatched,
				MinFollow:   request.Data.Translations.User.MinFollow,
			},
			Song: modules.YouTubeSongTranslations{
				Denied:               request.Data.Translations.Song.Denied,
				NotFound:             request.Data.Translations.Song.NotFound,
				AlreadyInQueue:       request.Data.Translations.Song.AlreadyInQueue,
				AgeRestrictions:      request.Data.Translations.Song.AgeRestrictions,
				CannotGetInformation: request.Data.Translations.Song.CannotGetInformation,
				Live:                 request.Data.Translations.Song.Live,
				MaxLength:            request.Data.Translations.Song.MaxLength,
				MinLength:            request.Data.Translations.Song.MinLength,
				RequestedMessage:     request.Data.Translations.Song.RequestedMessage,
				MaximumOrdered:       request.Data.Translations.Song.MaximumOrdered,
				MinViews:             request.Data.Translations.Song.MinViews,
			},
			Channel: modules.YouTubeChannelTranslations{
				Denied: request.Data.Translations.Channel.Denied,
			},
		},
	}

	bytes, err := json.Marshal(settings)
	if err != nil {
		return nil, err
	}

	entity.Settings = bytes
	if err := c.Db.WithContext(ctx).Save(entity).Error; err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
