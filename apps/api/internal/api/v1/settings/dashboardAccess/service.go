package dashboardaccess

import (
	"fmt"
	"github.com/samber/do"
	"github.com/satont/tsuwari/apps/api/internal/di"
	"github.com/satont/tsuwari/apps/api/internal/interfaces"
	cfg "github.com/satont/tsuwari/libs/config"
	"github.com/satont/tsuwari/libs/grpc/generated/tokens"
	"github.com/satont/tsuwari/libs/twitch"
	"net/http"

	model "github.com/satont/tsuwari/libs/gomodels"

	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
	"github.com/satont/go-helix/v2"
	"github.com/satont/tsuwari/apps/api/internal/types"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type Entity struct {
	model.ChannelsDashboardAccess
	TwitchUser *helix.User `json:"twitchUser"`
}

func handleGet(channelId string, services types.Services) ([]Entity, error) {
	logger := do.MustInvoke[interfaces.Logger](di.Provider)
	tokensGrpc := do.MustInvoke[tokens.TokensClient](di.Provider)
	config := do.MustInvoke[cfg.Config](di.Provider)

	twitchClient, err := twitch.NewAppClient(config, tokensGrpc)
	if err != nil {
		logger.Error(err)
		return nil, fiber.NewError(http.StatusInternalServerError, "internal error")
	}

	dbEntities := []model.ChannelsDashboardAccess{}
	err = services.DB.Where(`"channelId" = ?`, channelId).Find(&dbEntities).Error
	if err != nil {
		logger.Error(err)
		return nil, fiber.NewError(
			http.StatusInternalServerError,
			"cannot get dashboard users from db",
		)
	}

	usersIds := make([]string, 0, len(dbEntities))
	for _, u := range dbEntities {
		// usersIds[i] = u.UserID
		usersIds = append(usersIds, u.UserID)
	}

	twitchUsers, err := twitchClient.GetUsers(&helix.UsersParams{
		IDs: usersIds,
	})
	if err != nil {
		logger.Error(err)
		return nil, fiber.NewError(
			http.StatusInternalServerError,
			"error when getting users from twitch",
		)
	}

	if len(twitchUsers.Data.Users) == 0 {
		return []Entity{}, nil
	}

	entities := make([]Entity, 0, len(dbEntities))
	for _, dbEntity := range dbEntities {
		helixUser, ok := lo.Find(twitchUsers.Data.Users, func(u helix.User) bool {
			return u.ID == dbEntity.UserID
		})
		fmt.Println(ok, helixUser)
		entity := Entity{
			ChannelsDashboardAccess: dbEntity,
			TwitchUser:              lo.If(ok, &helixUser).Else(nil),
		}
		entities = append(entities, entity)
	}

	return entities, nil
}

func handlePost(channelId string, dto *addUserDto, services types.Services) (*Entity, error) {
	logger := do.MustInvoke[interfaces.Logger](di.Provider)
	tokensGrpc := do.MustInvoke[tokens.TokensClient](di.Provider)
	config := do.MustInvoke[cfg.Config](di.Provider)

	twitchClient, err := twitch.NewAppClient(config, tokensGrpc)
	if err != nil {
		logger.Error(err)
		return nil, fiber.NewError(http.StatusInternalServerError, "internal error")
	}

	twitchUsers, err := twitchClient.GetUsers(&helix.UsersParams{
		Logins: []string{dto.UserName},
	})
	if err != nil {
		logger.Error(err)
		return nil, fiber.NewError(http.StatusInternalServerError, "cannot get user from twitch")
	}

	if len(twitchUsers.Data.Users) == 0 {
		return nil, fiber.NewError(http.StatusNotFound, "cannot find user on twitch")
	}

	err = services.DB.
		Where(`"channelId" = ? AND "userId" = ?`, channelId, twitchUsers.Data.Users[0].ID).
		First(&model.ChannelsDashboardAccess{}).Error

	if err == nil {
		return nil, fiber.NewError(400, "that user already exists in db")
	}

	err = services.DB.Where(`"id" = ?`, twitchUsers.Data.Users[0].ID).First(&model.Users{}).Error
	if err != nil && err == gorm.ErrRecordNotFound {
		user := model.Users{
			ID:     twitchUsers.Data.Users[0].ID,
			ApiKey: uuid.NewV4().String(),
		}
		err = services.DB.Create(&user).Error
		fmt.Println(services.DB.ToSQL(func(tx *gorm.DB) *gorm.DB {
			return tx.Create(&user)
		}))
		if err != nil {
			logger.Error(err)
			return nil, fiber.NewError(http.StatusInternalServerError, "cannot save user in db")
		}
	}

	newAccess := model.ChannelsDashboardAccess{
		ID:        uuid.NewV4().String(),
		ChannelID: channelId,
		UserID:    twitchUsers.Data.Users[0].ID,
	}
	err = services.DB.Save(&newAccess).Error
	if err != nil {
		logger.Error(err)
		return nil, fiber.NewError(http.StatusInternalServerError, "cannot save user in db")
	}

	return &Entity{
		ChannelsDashboardAccess: newAccess,
		TwitchUser:              &twitchUsers.Data.Users[0],
	}, nil
}

func handleDelete(entityId string, services types.Services) error {
	access := model.ChannelsDashboardAccess{}
	err := services.DB.Where(`"id" = ?`, entityId).
		First(&access).
		Error
	if err != nil && err == gorm.ErrRecordNotFound {
		return fiber.NewError(http.StatusNotFound, "that entity not found in database")
	}
	if err != nil {
		return fiber.NewError(
			http.StatusInternalServerError,
			"something unexpected happend on our side",
		)
	}

	err = services.DB.Delete(&access).Error
	if err != nil {
		return fiber.NewError(
			http.StatusInternalServerError,
			"something unexpected happend on our side",
		)
	}

	return nil
}
