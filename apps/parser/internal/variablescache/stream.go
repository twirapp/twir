package variables_cache

import (
	"fmt"
	"github.com/samber/do"
	"github.com/satont/tsuwari/apps/parser/internal/di"
	config "github.com/satont/tsuwari/libs/config"
	"github.com/satont/tsuwari/libs/grpc/generated/tokens"
	"github.com/satont/tsuwari/libs/twitch"
	"gorm.io/gorm"

	model "github.com/satont/tsuwari/libs/gomodels"

	"github.com/lib/pq"
	"github.com/satont/go-helix/v2"
)

func (c *VariablesCacheService) GetChannelStream() *model.ChannelsStreams {
	cfg := do.MustInvoke[config.Config](di.Provider)
	tokensGrpc := do.MustInvoke[tokens.TokensClient](di.Provider)

	twitchClient, err := twitch.NewAppClient(cfg, tokensGrpc)
	if err != nil {
		return nil
	}

	c.locks.stream.Lock()
	defer c.locks.stream.Unlock()

	if c.cache.Stream != nil {
		return c.cache.Stream
	}

	db := do.MustInvoke[gorm.DB](di.Provider)
	stream := model.ChannelsStreams{}

	err = db.Where(`"userId" = ?`, c.ChannelId).First(&stream).Error

	if err != nil {
		fmt.Println(err)
		streams, err := twitchClient.GetStreams(&helix.StreamsParams{
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

		db.Save(&stream)
		c.cache.Stream = &stream
	} else {
		c.cache.Stream = &stream
	}

	return c.cache.Stream
}
