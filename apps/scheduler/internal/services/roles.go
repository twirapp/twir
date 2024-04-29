package services

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/samber/lo"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/logger"
	"gorm.io/gorm"
)

type Roles struct {
	db     *gorm.DB
	logger logger.Logger
}

func NewRoles(db *gorm.DB, l logger.Logger) *Roles {
	return &Roles{
		db:     db,
		logger: l,
	}
}

var rolesForCreate = []string{
	model.ChannelRoleTypeBroadcaster.String(),
	model.ChannelRoleTypeModerator.String(),
	model.ChannelRoleTypeSubscriber.String(),
	model.ChannelRoleTypeVip.String(),
}

func (c *Roles) CreateDefaultRoles(ctx context.Context, channelsIds []string) error {
	var channels []model.Channels
	if err := c.db.Where(
		`"id" IN ?`,
		channelsIds,
	).Preload("Roles").Find(&channels).Error; err != nil {
		return fmt.Errorf("cannot get channels: %w", err)
	}

	var wg sync.WaitGroup

	for _, channel := range channels {
		channel := channel

		wg.Add(1)
		go func() {
			defer wg.Done()
			for _, roleType := range rolesForCreate {
				isExists := lo.SomeBy(
					channel.Roles,
					func(item *model.ChannelRole) bool {
						return item.Type.String() == roleType
					},
				)
				if isExists {
					continue
				}

				if err := c.db.WithContext(ctx).Create(
					&model.ChannelRole{
						ID:                        uuid.New().String(),
						ChannelID:                 channel.ID,
						Name:                      roleType,
						Type:                      model.ChannelRoleEnum(roleType),
						Permissions:               pq.StringArray{},
						RequiredMessages:          0,
						RequiredWatchTime:         0,
						RequiredUsedChannelPoints: 0,
					},
				).Error; err != nil {
					c.logger.Error("cannot create role", slog.Any("err", err))
					return
				}
			}

			return
		}()
	}

	wg.Wait()

	return nil
}
