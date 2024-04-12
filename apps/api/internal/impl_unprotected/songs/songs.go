package songs

import (
	"context"
	"fmt"

	"github.com/samber/lo"
	"github.com/satont/twir/apps/api/internal/impl_deps"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/api/messages/songs_unprotected"
)

type Songs struct {
	*impl_deps.Deps
}

func (c *Songs) GetSongsQueue(
	ctx context.Context,
	req *songs_unprotected.GetSongsQueueRequest,
) (*songs_unprotected.GetSongsQueueResponse, error) {
	channel := &model.Channels{}
	if err := c.Db.
		WithContext(ctx).
		Where(`channels.id = ?`, req.ChannelId).
		Joins("User").
		First(channel).Error; err != nil {
		return nil, err
	}

	if channel.User.IsBanned {
		return &songs_unprotected.GetSongsQueueResponse{}, nil
	}

	var queue []*model.RequestedSong
	if err := c.Db.
		WithContext(ctx).
		Where(`"channelId" = ? AND "deletedAt" IS NULL`, req.ChannelId).
		Order(`"queuePosition" asc`).
		Find(&queue).Error; err != nil {
		return nil, fmt.Errorf("failed to get queue: %w", err)
	}

	songs := make([]*songs_unprotected.GetSongsQueueResponse_Song, len(queue))

	for i, s := range queue {
		requestedBy := lo.If(
			s.OrderedByDisplayName.Valid,
			s.OrderedByDisplayName.String,
		).Else(s.OrderedByName)
		songLink := lo.If(s.SongLink.Valid, s.SongLink.String).Else(
			fmt.Sprintf(
				"https://youtu.be/%s",
				s.ID,
			),
		)

		songs[i] = &songs_unprotected.GetSongsQueueResponse_Song{
			Title:       s.Title,
			RequestedBy: requestedBy,
			CreatedAt:   fmt.Sprint(s.CreatedAt.UnixMilli()),
			SongLink:    songLink,
			Duration:    s.Duration,
		}
	}

	return &songs_unprotected.GetSongsQueueResponse{
		Songs: songs,
	}, nil
}
