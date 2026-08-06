package mappers

import (
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
	entity "github.com/twirapp/twir/libs/entities/spotify_song_request"
)

func SpotifySongRequestToGQL(request entity.SpotifySongRequest) gqlmodel.SpotifySongRequest {
	return gqlmodel.SpotifySongRequest{
		ID:                   request.ID,
		Title:                request.Title,
		Artist:               request.Artist,
		Album:                request.Album,
		DurationMs:           request.DurationMs,
		RequesterName:        request.RequesterName,
		RequesterDisplayName: request.RequesterDisplayName,
		Source:               request.Source,
		QueuePosition:        request.QueuePosition,
		Status:               SpotifySongRequestStatusToGQL(request.Status),
		CreatedAt:            request.CreatedAt,
	}
}

func SpotifySongRequestStatusToGQL(status entity.Status) gqlmodel.SpotifySongRequestStatus {
	switch status {
	case entity.StatusQueued:
		return gqlmodel.SpotifySongRequestStatusQueued
	case entity.StatusPlaying:
		return gqlmodel.SpotifySongRequestStatusPlaying
	case entity.StatusPlayed:
		return gqlmodel.SpotifySongRequestStatusPlayed
	case entity.StatusCancelledPendingSkip:
		return gqlmodel.SpotifySongRequestStatusCancelledPendingSkip
	case entity.StatusSkippedByTwir:
		return gqlmodel.SpotifySongRequestStatusSkippedByTwir
	case entity.StatusRemovedOrReconciled:
		return gqlmodel.SpotifySongRequestStatusRemovedOrReconciled
	default:
		return gqlmodel.SpotifySongRequestStatusUnknown
	}
}
