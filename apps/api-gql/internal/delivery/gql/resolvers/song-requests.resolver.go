package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.76

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"sync"

	"github.com/guregu/null"
	"github.com/raitonoberu/ytsearch"
	"github.com/samber/lo"
	loParallel "github.com/samber/lo/parallel"
	data_loader "github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/dataloader"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/graph"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/mappers"
	model "github.com/twirapp/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/logger/audit"
	"gorm.io/gorm"
)

// SongRequestsUpdate is the resolver for the songRequestsUpdate field.
func (r *mutationResolver) SongRequestsUpdate(ctx context.Context, opts gqlmodel.SongRequestsSettingsOpts) (bool, error) {
	dashboardId, err := r.deps.Sessions.GetSelectedDashboard(ctx)
	if err != nil {
		return false, err
	}

	user, err := r.deps.Sessions.GetAuthenticatedUser(ctx)
	if err != nil {
		return false, err
	}

	entity := model.ChannelSongRequestsSettings{
		ChannelID: dashboardId,
	}
	err = r.deps.Gorm.WithContext(ctx).Where(
		`"channel_id" = ?`,
		dashboardId,
	).FirstOrCreate(&entity).Error
	if err != nil {
		return false, fmt.Errorf("failed to update song requests settings: %w", err)
	}

	entity = model.ChannelSongRequestsSettings{
		ID:                                   entity.ID,
		ChannelID:                            dashboardId,
		Enabled:                              opts.Enabled,
		AcceptOnlyWhenOnline:                 opts.AcceptOnlyWhenOnline,
		PlayerNoCookieMode:                   opts.PlayerNoCookieMode,
		TakeSongFromDonationMessage:          opts.TakeSongFromDonationMessages,
		MaxRequests:                          opts.MaxRequests,
		ChannelPointsRewardID:                null.StringFromPtr(opts.ChannelPointsRewardID.Value()),
		AnnouncePlay:                         opts.AnnouncePlay,
		NeededVotesForSkip:                   float64(opts.NeededVotesForSkip),
		UserMaxRequests:                      opts.User.MaxRequests,
		UserMinWatchTime:                     int64(opts.User.MinWatchTime),
		UserMinMessages:                      opts.User.MinMessages,
		UserMinFollowTime:                    opts.User.MinFollowTime,
		SongMinLength:                        opts.Song.MinLength,
		SongMaxLength:                        opts.Song.MaxLength,
		SongMinViews:                         opts.Song.MinViews,
		SongAcceptedCategories:               opts.Song.AcceptedCategories,
		SongWordsDenyList:                    opts.DenyList.Words,
		DenyListUsers:                        opts.DenyList.Users,
		DenyListSongs:                        opts.DenyList.Songs,
		DenyListChannels:                     opts.DenyList.Channels,
		DenyListArtistsNames:                 opts.DenyList.ArtistsNames,
		DenyListWords:                        opts.DenyList.Words,
		TranslationsNowPlaying:               opts.Translations.NowPlaying,
		TranslationsNotEnabled:               opts.Translations.NotEnabled,
		TranslationsNoText:                   opts.Translations.NoText,
		TranslationsAcceptOnlineWhenOnline:   opts.Translations.AcceptOnlyWhenOnline,
		TranslationsUserDenied:               opts.Translations.User.Denied,
		TranslationsUserMaxRequests:          opts.Translations.User.MaxRequests,
		TranslationsUserMinMessages:          opts.Translations.User.MinMessages,
		TranslationsUserMinWatched:           opts.Translations.User.MinWatched,
		TranslationsUserMinFollow:            opts.Translations.User.MinFollow,
		TranslationsSongDenied:               opts.Translations.Song.Denied,
		TranslationsSongNotFound:             opts.Translations.Song.NotFound,
		TranslationsSongAlreadyInQueue:       opts.Translations.Song.AlreadyInQueue,
		TranslationsSongAgeRestrictions:      opts.Translations.Song.AgeRestrictions,
		TranslationsSongCannotGetInformation: opts.Translations.Song.CannotGetInformation,
		TranslationsSongLive:                 opts.Translations.Song.Live,
		TranslationsSongMaxLength:            opts.Translations.Song.MaxLength,
		TranslationsSongMinLength:            opts.Translations.Song.MinLength,
		TranslationsSongRequestedMessage:     opts.Translations.Song.RequestedMessage,
		TranslationsSongMaximumOrdered:       opts.Translations.Song.MaximumOrdered,
		TranslationsSongMinViews:             opts.Translations.Song.MinViews,
		TranslationsChannelDenied:            opts.Translations.Channel.Denied,
	}

	err = r.deps.Gorm.WithContext(ctx).Save(&entity).Error
	if err != nil {
		return false, fmt.Errorf("failed to update song requests settings: %w", err)
	}

	r.deps.Logger.Audit(
		"Song requests updated",
		audit.Fields{
			NewValue:      entity,
			ActorID:       lo.ToPtr(user.ID),
			ChannelID:     lo.ToPtr(dashboardId),
			System:        mappers.AuditSystemToTableName(gqlmodel.AuditLogSystemChannelSongRequests),
			OperationType: audit.OperationUpdate,
			ObjectID:      &entity.ID,
		},
	)

	if err := r.deps.ChannelSongRequestsSettingsCache.Invalidate(ctx, dashboardId); err != nil {
		r.deps.Logger.Error("failed to invalidate song requests settings cache", err)
	}

	return true, nil
}

// SongRequests is the resolver for the songRequests field.
func (r *queryResolver) SongRequests(ctx context.Context) (*gqlmodel.SongRequestsSettings, error) {
	dashboardId, err := r.deps.Sessions.GetSelectedDashboard(ctx)
	if err != nil {
		return nil, err
	}

	entity := model.ChannelSongRequestsSettings{}
	err = r.deps.Gorm.WithContext(ctx).Where(`"channel_id" = ?`, dashboardId).First(&entity).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get song requests settings: %w", err)
	}

	return &gqlmodel.SongRequestsSettings{
		Enabled:               entity.Enabled,
		AcceptOnlyWhenOnline:  entity.AcceptOnlyWhenOnline,
		MaxRequests:           entity.MaxRequests,
		ChannelPointsRewardID: entity.ChannelPointsRewardID.Ptr(),
		AnnouncePlay:          entity.AnnouncePlay,
		NeededVotesForSkip:    int(entity.NeededVotesForSkip),
		User: &gqlmodel.SongRequestsUserSettings{
			MaxRequests:   entity.UserMaxRequests,
			MinWatchTime:  int(entity.UserMinWatchTime),
			MinMessages:   entity.UserMinMessages,
			MinFollowTime: entity.UserMinFollowTime,
		},
		Song: &gqlmodel.SongRequestsSongSettings{
			MinLength:          entity.SongMinLength,
			MaxLength:          entity.SongMaxLength,
			MinViews:           entity.SongMinViews,
			AcceptedCategories: entity.SongAcceptedCategories,
		},
		DenyList: &gqlmodel.SongRequestsDenyList{
			Users:        entity.DenyListUsers,
			Songs:        entity.DenyListSongs,
			Channels:     entity.DenyListChannels,
			ArtistsNames: entity.DenyListArtistsNames,
			Words:        entity.DenyListWords,
		},
		Translations: &gqlmodel.SongRequestsTranslations{
			NowPlaying:           entity.TranslationsNowPlaying,
			NotEnabled:           entity.TranslationsNotEnabled,
			NoText:               entity.TranslationsNoText,
			AcceptOnlyWhenOnline: entity.TranslationsAcceptOnlineWhenOnline,
			User: &gqlmodel.SongRequestsUserTranslations{
				Denied:      entity.TranslationsUserDenied,
				MaxRequests: entity.TranslationsUserMaxRequests,
				MinMessages: entity.TranslationsUserMinMessages,
				MinWatched:  entity.TranslationsUserMinWatched,
				MinFollow:   entity.TranslationsUserMinFollow,
			},
			Song: &gqlmodel.SongRequestsSongTranslations{
				Denied:               entity.TranslationsSongDenied,
				NotFound:             entity.TranslationsSongNotFound,
				AlreadyInQueue:       entity.TranslationsSongAlreadyInQueue,
				AgeRestrictions:      entity.TranslationsSongAgeRestrictions,
				CannotGetInformation: entity.TranslationsSongCannotGetInformation,
				Live:                 entity.TranslationsSongLive,
				MaxLength:            entity.TranslationsSongMaxLength,
				MinLength:            entity.TranslationsSongMinLength,
				RequestedMessage:     entity.TranslationsSongRequestedMessage,
				MaximumOrdered:       entity.TranslationsSongMaximumOrdered,
				MinViews:             entity.TranslationsSongMinViews,
			},
			Channel: &gqlmodel.SongRequestsChannelTranslations{
				Denied: entity.TranslationsChannelDenied,
			},
		},
		TakeSongFromDonationMessages: entity.TakeSongFromDonationMessage,
		PlayerNoCookieMode:           entity.PlayerNoCookieMode,
	}, nil
}

// SongRequestsSearchChannelOrVideo is the resolver for the songRequestsSearchChannelOrVideo field.
func (r *queryResolver) SongRequestsSearchChannelOrVideo(ctx context.Context, opts gqlmodel.SongRequestsSearchChannelOrVideoOpts) (*gqlmodel.SongRequestsSearchChannelOrVideoResponse, error) {
	response := &gqlmodel.SongRequestsSearchChannelOrVideoResponse{
		Items: make([]gqlmodel.SongRequestsSearchChannelOrVideoItem, 0, len(opts.Query)),
	}

	var getThumbNailUrl = func(url string) string {
		return strings.Replace(url, "http://", "https://", 1)
	}

	if len(opts.Query) == 0 {
		return response, nil
	}

	var mu sync.Mutex
	loParallel.ForEach(
		opts.Query,
		func(query string, _ int) {
			if query == "" {
				return
			}

			var search *ytsearch.SearchClient
			if opts.Type == gqlmodel.SongRequestsSearchChannelOrVideoOptsTypeChannel {
				search = ytsearch.ChannelSearch(query)
			} else {
				search = ytsearch.VideoSearch(query)
			}

			res, err := search.Next()
			if err != nil {
				r.deps.Logger.Error(
					"cannot find",
					slog.String("query", query),
				)
				return
			}

			mu.Lock()
			defer mu.Unlock()
			if opts.Type == gqlmodel.SongRequestsSearchChannelOrVideoOptsTypeChannel {
				channels := lo.Map(
					res.Channels,
					func(item *ytsearch.ChannelItem, _ int) gqlmodel.SongRequestsSearchChannelOrVideoItem {
						thumb := getThumbNailUrl(item.Thumbnails[len(item.Thumbnails)-1].URL)
						return gqlmodel.SongRequestsSearchChannelOrVideoItem{
							ID:        item.ID,
							Title:     item.Title,
							Thumbnail: thumb,
						}
					},
				)
				response.Items = append(
					response.Items,
					channels...,
				)
			} else {
				videos := lo.Map(
					res.Videos,
					func(item *ytsearch.VideoItem, _ int) gqlmodel.SongRequestsSearchChannelOrVideoItem {
						thumb := getThumbNailUrl(item.Thumbnails[len(item.Thumbnails)-1].URL)

						return gqlmodel.SongRequestsSearchChannelOrVideoItem{
							ID:        item.ID,
							Title:     item.Title,
							Thumbnail: thumb,
						}
					},
				)
				response.Items = append(
					response.Items,
					videos...,
				)
			}
		},
	)

	return response, nil
}

// SongRequestsPublicQueue is the resolver for the songRequestsPublicQueue field.
func (r *queryResolver) SongRequestsPublicQueue(ctx context.Context, channelID string) ([]gqlmodel.SongRequestPublic, error) {
	queue, err := r.deps.SongRequestsService.GetPublicQueue(ctx, channelID)
	if err != nil {
		return nil, err
	}

	mapped := make([]gqlmodel.SongRequestPublic, 0, len(queue))
	for _, song := range queue {
		mapped = append(mapped, mappers.SongRequestPublicTo(song))
	}

	return mapped, nil
}

// TwitchProfile is the resolver for the twitchProfile field.
func (r *songRequestPublicResolver) TwitchProfile(ctx context.Context, obj *gqlmodel.SongRequestPublic) (*gqlmodel.TwirUserTwitchInfo, error) {
	return data_loader.GetHelixUserById(ctx, obj.UserID)
}

// SongRequestPublic returns graph.SongRequestPublicResolver implementation.
func (r *Resolver) SongRequestPublic() graph.SongRequestPublicResolver {
	return &songRequestPublicResolver{r}
}

type songRequestPublicResolver struct{ *Resolver }
