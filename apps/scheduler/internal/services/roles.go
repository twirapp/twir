package services

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/samber/lo"
	model "github.com/satont/twir/libs/gomodels"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
)

type Roles struct {
	db *gorm.DB
}

func NewRoles(db *gorm.DB) *Roles {
	return &Roles{db}
}

var rolesForCreate = []string{
	model.ChannelRoleTypeBroadcaster.String(),
	model.ChannelRoleTypeModerator.String(),
	model.ChannelRoleTypeSubscriber.String(),
	model.ChannelRoleTypeVip.String(),
}

func (c *Roles) CreateDefaultRoles(ctx context.Context, channelsIds []string) error {
	var channels []model.Channels
	if err := c.db.Where(`"id" IN ?`, channelsIds).Preload("Roles").Find(&channels).Error; err != nil {
		return err
	}

	g, ctx := errgroup.WithContext(ctx)
	for _, channel := range channels {
		channel := channel

		g.Go(
			func() error {
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

					settings, err := json.Marshal(&model.ChannelRoleSettings{})
					if err != nil {
						return err
					}

					if err := c.db.Create(
						&model.ChannelRole{
							ID:          uuid.New().String(),
							ChannelID:   channel.ID,
							Name:        roleType,
							Type:        model.ChannelRoleEnum(roleType),
							Permissions: pq.StringArray{},
							Settings:    settings,
						},
					).Error; err != nil {
						return err
					}
				}

				return nil
			},
		)
	}

	if err := g.Wait(); err != nil {
		return err
	}

	return nil
}
