package song_requests

import (
	"github.com/gofiber/fiber/v2"
	"github.com/satont/tsuwari/apps/api/internal/types"
	model "github.com/satont/tsuwari/libs/gomodels"
	"net/http"
)

type Song struct {
	Title                string  `json:"title"`
	VideoID              string  `json:"videoId"`
	Duration             int32   `json:"duration"`
	OrderedByName        string  `json:"orderedByName"`
	OrderedByDisplayName *string `json:"orderedByDisplayName"`
}

func handleGet(channelId string, services types.Services) ([]Song, error) {
	songs := []model.RequestedSong{}

	err := services.DB.
		Where(`"channelId" = ? AND "deletedAt" IS NULL`, channelId).
		Order(`"queuePosition" asc`).
		Find(&songs).Error

	if err != nil {
		return nil, fiber.NewError(http.StatusNotFound, "cannot find commands")
	}

	songsResponse := []Song{}

	for _, song := range songs {
		songsResponse = append(songsResponse, Song{
			Title:                song.Title,
			Duration:             song.Duration,
			VideoID:              song.VideoID,
			OrderedByName:        song.OrderedByName,
			OrderedByDisplayName: song.OrderedByDisplayName.Ptr(),
		})
	}

	return songsResponse, nil
}
