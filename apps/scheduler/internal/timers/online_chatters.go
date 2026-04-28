package timers

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/nicklaw5/helix/v2"
	"github.com/samber/lo"
	buscore "github.com/twirapp/twir/libs/bus-core"
	config "github.com/twirapp/twir/libs/config"
	model "github.com/twirapp/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/logger"
	"github.com/twirapp/twir/libs/twitch"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type OnlineUsersOpts struct {
	fx.In
	Lc fx.Lifecycle

	Logger *slog.Logger
	Config config.Config

	Gorm    *gorm.DB
	TwirBus *buscore.Bus
}

type onlineUsers struct {
	config  config.Config
	logger  *slog.Logger
	db      *gorm.DB
	twirBus *buscore.Bus
}

func NewOnlineUsers(opts OnlineUsersOpts) {
	timeTick := 15 * time.Second
	if opts.Config.AppEnv == "production" {
		timeTick = 5 * time.Minute
	}
	ticker := time.NewTicker(timeTick)

	ctx, cancel := context.WithCancel(context.Background())

	s := &onlineUsers{
		config:  opts.Config,
		logger:  opts.Logger,
		db:      opts.Gorm,
		twirBus: opts.TwirBus,
	}

	opts.Lc.Append(
		fx.Hook{
			OnStart: func(_ context.Context) error {
				go func() {
					for {
						select {
						case <-ctx.Done():
							ticker.Stop()
							return
						case <-ticker.C:
							s.updateOnlineUsers(ctx)
						}
					}
				}()

				return nil
			},
			OnStop: func(_ context.Context) error {
				cancel()
				return nil
			},
		},
	)
}

func (c *onlineUsers) updateOnlineUsers(ctx context.Context) {
	streams, err := c.getStreams(ctx)
	if err != nil {
		c.logger.Error("cannot get streams", logger.Error(err))
		return
	}

	var wg sync.WaitGroup
	for _, stream := range streams {
		if c.shouldSkipStream(stream) {
			continue
		}

		s := stream
		wg.Add(1)

		go func() {
			defer wg.Done()
			if updateErr := c.updateStreamUsers(ctx, s); updateErr != nil {
				c.logger.Error("cannot update stream users", slog.Any("err", updateErr))
			}
		}()
	}

	wg.Wait()
}

func (c *onlineUsers) getStreams(
	ctx context.Context,
) ([]*model.ChannelsStreams, error) {
	var streams []*model.ChannelsStreams
	err := c.db.WithContext(ctx).Preload("Channel").Preload("Channel.User").Find(&streams).Error
	return streams, err
}

func (c *onlineUsers) shouldSkipStream(stream *model.ChannelsStreams) bool {
	return stream.Channel == nil || (!stream.Channel.IsEnabled || stream.Channel.User.IsBanned)
}

func (c *onlineUsers) updateStreamUsers(
	ctx context.Context,
	stream *model.ChannelsStreams,
) error {
	broadcasterID := stream.UserId
	channelUUID := stream.Channel.ID

	twitchClient, err := twitch.NewUserClientWithContext(
		ctx,
		stream.Channel.User.ID,
		c.config,
		c.twirBus,
	)
	if err != nil {
		return err
	}

	var cursor string
	for {
		params := &helix.GetChatChattersParams{
			BroadcasterID: broadcasterID,
			ModeratorID:   broadcasterID,
			After:         cursor,
			First:         "1000",
		}
		req, err := twitchClient.GetChannelChatChatters(params)
		if err != nil {
			return fmt.Errorf("cannot get channel chat chatters: %w", err)
		}
		if req.ErrorMessage != "" {
			return fmt.Errorf("cannot get channel chat chatters: %s", req.ErrorMessage)
		}

		chatters := req.Data.Chatters
		if len(chatters) == 0 {
			return nil
		}

		chatterPlatformIDs := lo.Map(
			chatters,
			func(chatter helix.ChatChatter, _ int) string {
				return chatter.UserID
			},
		)

		err = c.db.WithContext(ctx).Transaction(
			func(tx *gorm.DB) error {
				insertParts := make([]string, 0, len(chatters))
				insertArgs := make([]interface{}, 0, len(chatters)*2)
				for i, chatter := range chatters {
					base := i * 2
					insertParts = append(
						insertParts,
						fmt.Sprintf("(uuidv7(), 'twitch', $%d, $%d)", base+1, base+2),
					)
					insertArgs = append(insertArgs, chatter.UserID, uuid.New().String())
				}
				upsertUsersSQL := `INSERT INTO users (id, platform, platform_id, "apiKey") VALUES ` +
					strings.Join(insertParts, ", ") +
					` ON CONFLICT (platform, platform_id) DO NOTHING`
				if err := tx.Exec(upsertUsersSQL, insertArgs...).Error; err != nil {
					return fmt.Errorf("cannot upsert users: %w", err)
				}

				type userRow struct {
					ID         string `gorm:"column:id"`
					PlatformID string `gorm:"column:platform_id"`
				}
				var userRows []userRow
				if err := tx.Raw(
					`SELECT id, platform_id FROM users WHERE platform_id IN ? AND platform = 'twitch'`,
					chatterPlatformIDs,
				).Scan(&userRows).Error; err != nil {
					return fmt.Errorf("cannot fetch user UUIDs: %w", err)
				}

				platformIDToUUID := make(map[string]string, len(userRows))
				for _, row := range userRows {
					platformIDToUUID[row.PlatformID] = row.ID
				}

				if len(userRows) > 0 {
					statsParts := make([]string, 0, len(userRows))
					statsArgs := make([]interface{}, 0, len(userRows)*3)
					for i, row := range userRows {
						base := i * 3
						statsParts = append(
							statsParts,
							fmt.Sprintf("($%d, $%d::uuid, $%d::uuid)", base+1, base+2, base+3),
						)
						statsArgs = append(statsArgs, uuid.New().String(), row.ID, channelUUID)
					}
					statsSQL := `INSERT INTO users_stats (id, "userId", "channelId") VALUES ` +
						strings.Join(statsParts, ", ") +
						` ON CONFLICT ("channelId", "userId") DO NOTHING`
					if err := tx.Exec(statsSQL, statsArgs...).Error; err != nil {
						return fmt.Errorf("cannot upsert users stats: %w", err)
					}
				}

				if err := tx.Exec(
					`DELETE FROM users_online WHERE "channelId" = $1::uuid`,
					channelUUID,
				).Error; err != nil {
					return fmt.Errorf("cannot delete online users: %w", err)
				}

				onlineParts := make([]string, 0, len(chatters))
				onlineArgs := make([]interface{}, 0, len(chatters)*4)
				argIdx := 0
				for _, chatter := range chatters {
					userUUID, ok := platformIDToUUID[chatter.UserID]
					if !ok {
						continue
					}
					onlineParts = append(
						onlineParts,
						fmt.Sprintf(
							"($%d, $%d::uuid, $%d::uuid, $%d)",
							argIdx+1, argIdx+2, argIdx+3, argIdx+4,
						),
					)
					onlineArgs = append(
						onlineArgs,
						uuid.New().String(), channelUUID, userUUID, chatter.UserLogin,
					)
					argIdx += 4
				}
				if len(onlineParts) > 0 {
					onlineSQL := `INSERT INTO users_online (id, "channelId", "userId", "userName") VALUES ` +
						strings.Join(onlineParts, ", ")
					if err := tx.Exec(onlineSQL, onlineArgs...).Error; err != nil {
						return fmt.Errorf("cannot insert online users: %w", err)
					}
				}

				return nil
			},
		)

		if err != nil {
			return fmt.Errorf("cannot update stream users: %w", err)
		}

		if req.Data.Pagination.Cursor == "" {
			return nil
		}

		cursor = req.Data.Pagination.Cursor
	}
}
