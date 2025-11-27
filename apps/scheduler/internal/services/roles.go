package services

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/samber/lo"
	model "github.com/twirapp/twir/libs/gomodels"
	"gorm.io/gorm"
)

type Roles struct {
	db     *gorm.DB
	logger *slog.Logger
}

func NewRoles(db *gorm.DB, l *slog.Logger) *Roles {
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

type channelWithRolesToCreate struct {
	ChannelID     string         `gorm:"column:channelId" db:"channelId"`
	RolesToCreate pq.StringArray `gorm:"column:rolesToCreate" db:"rolesToCreate"`
}

func (c *Roles) CreateDefaultRoles(ctx context.Context) error {
	var placeholders []string
	var args []interface{}

	for _, roleType := range rolesForCreate {
		placeholders = append(placeholders, "(?)")
		args = append(args, roleType)
	}

	valuesClause := strings.Join(placeholders, ", ")

	query := fmt.Sprintf(
		`
		SELECT
			c.id AS "channelId",
			array_agg(required_roles.type) AS "rolesToCreate"
		FROM
			public.channels c
			CROSS JOIN (
				VALUES %s
			) AS required_roles(type)
		LEFT JOIN
			public.channels_roles r ON c.id = r."channelId" AND required_roles.type::channels_roles_type_enum = r.type
		WHERE
			r.id IS NULL
		GROUP BY
			c.id;
	`, valuesClause,
	)

	var channels []channelWithRolesToCreate
	if err := c.db.
		WithContext(ctx).
		Raw(query, args...).
		Scan(&channels).
		Error; err != nil {
		return fmt.Errorf("cannot get channels with roles to create: %w", err)
	}

	forCreate := make([]model.ChannelRole, 0, len(channels))
	for _, channel := range channels {
		for _, role := range channel.RolesToCreate {
			foundRole, ok := lo.Find(
				rolesForCreate,
				func(item string) bool {
					return item == role
				},
			)
			if !ok {
				continue
			}

			forCreate = append(
				forCreate,
				model.ChannelRole{
					ID:                        uuid.New().String(),
					ChannelID:                 channel.ChannelID,
					Name:                      foundRole,
					Type:                      model.ChannelRoleEnum(foundRole),
					Permissions:               pq.StringArray{},
					RequiredMessages:          0,
					RequiredWatchTime:         0,
					RequiredUsedChannelPoints: 0,
				},
			)
		}
	}

	if len(forCreate) == 0 {
		return nil
	}

	if err := c.db.WithContext(ctx).CreateInBatches(&forCreate, 1000).Error; err != nil {
		return fmt.Errorf("cannot create default roles: %w", err)
	}

	c.logger.Info(
		"Created default roles for channels",
		slog.Int("channels", len(channels)),
		slog.Int("roles", len(forCreate)),
	)

	return nil
}
