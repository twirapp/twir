package variables_cache

import (
	"fmt"
	model "tsuwari/models"

	"github.com/lib/pq"
	"github.com/satont/go-helix/v2"
)

func (c *VariablesCacheService) GetChannelStream() *model.ChannelsStreams {
	c.locks.stream.Lock()
	defer c.locks.stream.Unlock()

	if c.cache.Stream != nil {
		return c.cache.Stream
	}

	stream := model.ChannelsStreams{}

	err := c.Services.Db.Where(`"userId" = ?`, c.ChannelId).First(&stream).Error

	fmt.Printf("%+v\n", stream)

	if err != nil {
		fmt.Println(err)
		streams, err := c.Services.Twitch.Client.GetStreams(&helix.StreamsParams{
			UserIDs: []string{c.ChannelId},
		})

		if err != nil || len(streams.Data.Streams) == 0 {
			return nil
		}

		helixStream := streams.Data.Streams[0]

		tags := pq.StringArray{}
		for _, t := range helixStream.TagIDs {
			tags = append(tags, t)
		}
		stream = model.ChannelsStreams{
			ID:             helixStream.ID,
			UserId:         helixStream.UserID,
			UserLogin:      helixStream.UserLogin,
			UserName:       helixStream.UserName,
			GameId:         helixStream.GameID,
			GameName:       helixStream.GameName,
			CommunityIds:   []string{},
			Type:           helixStream.Type,
			Title:          helixStream.Title,
			ViewerCount:    helixStream.ViewerCount,
			StartedAt:      helixStream.StartedAt,
			Language:       helixStream.Language,
			ThumbnailUrl:   helixStream.ThumbnailURL,
			TagIds:         &tags,
			IsMature:       helixStream.IsMature,
			ParsedMessages: 0,
		}

		c.Services.Db.Save(&stream)
		c.cache.Stream = &stream
	} else {
		c.cache.Stream = &stream
	}

	return c.cache.Stream
}
