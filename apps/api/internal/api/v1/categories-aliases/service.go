package categories_aliases

import (
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/samber/do"
	"github.com/satont/tsuwari/apps/api/internal/di"
	"github.com/satont/tsuwari/apps/api/internal/interfaces"
	"github.com/satont/tsuwari/apps/api/internal/types"
	model "github.com/satont/tsuwari/libs/gomodels"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

func handleGetCategoryAliases(channelId string, services types.Services) ([]model.ChannelCategoryAlias, error) {
	db := do.MustInvoke[*gorm.DB](di.Provider)
	logger := do.MustInvoke[interfaces.Logger](di.Provider)

	var aliases []model.ChannelCategoryAlias
	err := db.Where(`"channelId" = ?`, channelId).Find(&aliases).Error
	if err != nil {
		logger.Error(err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}

	return aliases, nil
}

func handlePost(channelId string, dto *categoryAliasDto, services types.Services) (*model.ChannelCategoryAlias, error) {
	logger := do.MustInvoke[interfaces.Logger](di.Provider)
	db := do.MustInvoke[*gorm.DB](di.Provider)

	existedAlias := &model.ChannelCategoryAlias{}
	err := db.Where(`"channelId" = ? AND "alias" = ?`, channelId, dto.Alias).First(existedAlias).Error
	if err == nil {
		return nil, fiber.NewError(http.StatusBadRequest, "alias with this name already exists")
	}
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Error(err)
			return nil, fiber.NewError(http.StatusInternalServerError, "cannot get info about aliases")
		}
	}

	alias := &model.ChannelCategoryAlias{
		ID:         uuid.NewV4().String(),
		ChannelID:  channelId,
		Category:   dto.Category,
		CategoryId: dto.CategoryId,
		Alias:      dto.Alias,
	}

	err = db.Save(alias).Error
	if err != nil {
		logger.Error(err)
		return nil, fiber.NewError(http.StatusInternalServerError, "cannot create category alias")
	}

	return alias, nil
}

func handleDelete(categoryAliasId string, services types.Services) error {
	logger := do.MustInvoke[interfaces.Logger](di.Provider)
	db := do.MustInvoke[*gorm.DB](di.Provider)

	existedAlias := &model.ChannelCategoryAlias{}
	err := db.Where(`"id" = ?`, categoryAliasId).First(existedAlias).Error
	if err != nil {
		return fiber.NewError(http.StatusNotFound, "alias not found")
	}

	err = db.Where(`"id" = ?`, categoryAliasId).Delete(&model.ChannelCategoryAlias{}).Error
	if err != nil {
		logger.Error(err)
		return fiber.NewError(http.StatusInternalServerError, "cannot delete category alias")
	}

	return nil
}

func handlePatch(
	channelId string,
	categoryAliasId string,
	dto *categoryAliasPatchDto,
	services types.Services,
) (*model.ChannelCategoryAlias, error) {
	logger := do.MustInvoke[interfaces.Logger](di.Provider)
	db := do.MustInvoke[*gorm.DB](di.Provider)

	existedAlias := &model.ChannelCategoryAlias{}
	err := db.Where(`"id" = ?`, categoryAliasId).First(existedAlias).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fiber.NewError(http.StatusNotFound, "alias not found")
		}
		logger.Error(err)
		return nil, fiber.NewError(http.StatusInternalServerError)
	}

	err = db.Where(`"channelId" = ? AND "alias" = ?`, channelId, dto.Alias).First(&model.ChannelCategoryAlias{}).Error
	if err == nil {
		return nil, fiber.NewError(http.StatusBadRequest, "alias with this name already exists")
	}
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Error(err)
			return nil, fiber.NewError(http.StatusInternalServerError)
		}
	}

	if dto.CategoryId != nil {
		existedAlias.CategoryId = *dto.CategoryId
	}
	if dto.Alias != nil {
		existedAlias.Alias = *dto.Alias
	}
	if dto.Category != nil {
		existedAlias.Category = *dto.Category
	}
	err = db.Model(existedAlias).Select("*").Updates(existedAlias).Error
	if err != nil {
		return nil, fiber.NewError(http.StatusInternalServerError, "cannot update category alias")
	}

	return existedAlias, nil
}
