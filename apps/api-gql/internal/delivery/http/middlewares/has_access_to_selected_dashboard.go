package middlewares

import (
	"errors"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	model "github.com/twirapp/twir/libs/gomodels"
	"gorm.io/gorm"
)

func (c *Middlewares) HasAccessToSelectedDashboard(hc huma.Context, next func(huma.Context)) {
	ctx := hc.Context()

	user, err := c.auth.GetAuthenticatedUser(ctx)
	if err != nil {
		huma.WriteErr(
			c.huma,
			hc,
			http.StatusInternalServerError,
			"Cannot get authenticated user",
			err,
		)

		return
	}

	dashboardId, err := c.auth.GetSelectedDashboard(ctx)
	if err != nil {
		huma.WriteErr(
			c.huma,
			hc,
			http.StatusInternalServerError,
			"Cannot get selected dashboard",
			err,
		)

		return
	}

	if user.ID == dashboardId || user.IsBotAdmin {
		next(hc)
		return
	}

	var channelRoles []model.ChannelRole
	if err := c.gorm.
		WithContext(ctx).
		Where(`"channelId" = ?`, dashboardId).
		Preload("Users", `"userId" = ?`, user.ID).
		Find(&channelRoles).
		Error; err != nil {
		huma.WriteErr(
			c.huma,
			hc,
			http.StatusInternalServerError,
			"Cannot get channel roles",
			err,
		)

		return
	}

	var userStat model.UsersStats
	if err := c.gorm.
		WithContext(ctx).
		Where(`"userId" = ? AND "channelId" = ?`, user.ID, dashboardId).
		First(&userStat).
		Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {

		huma.WriteErr(
			c.huma,
			hc,
			http.StatusInternalServerError,
			"Cannot get user stats",
			err,
		)

		return
	}

	roleToStats := map[model.ChannelRoleEnum]bool{
		model.ChannelRoleTypeModerator:  userStat.IsMod,
		model.ChannelRoleTypeVip:        userStat.IsVip,
		model.ChannelRoleTypeSubscriber: userStat.IsSubscriber,
	}

	for i, role := range channelRoles {
		if roleToStats[role.Type] {
			channelRoles[i].Users = append(
				role.Users,
				&model.ChannelRoleUser{
					ID:     "", // not needed
					UserID: user.ID,
					RoleID: role.ID,
				},
			)
		}
	}

	for _, role := range channelRoles {
		// we do not check does role.Users contains request author user
		// because we are doing preload by user id
		if len(role.Users) == 0 || len(role.Permissions) == 0 {
			continue
		}

		for _, roleUser := range role.Users {
			if roleUser.UserID == user.ID {
				next(hc)

				return
			}
		}
	}

	huma.WriteErr(
		c.huma,
		hc,
		http.StatusForbidden,
		"user does not have access to selected dashboard",
		nil,
	)

	return
}
